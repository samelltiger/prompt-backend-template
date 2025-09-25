// internal/service/sync_service.go
package service

import (
	"encoding/json"
	"log"
	"time"

	"gorm.io/gorm"

	"llmapisrv/config"
	"llmapisrv/internal/model"
	"llmapisrv/pkg/cache"
	"llmapisrv/pkg/logger"
)

type SyncService struct {
	gatewayDB *gorm.DB
	newAPIDB  *gorm.DB
	config    *config.Config
	cache     *cache.RedisCache
}

func NewSyncService(gatewayDB, newAPIDB *gorm.DB, config *config.Config, cache *cache.RedisCache) *SyncService {
	return &SyncService{
		gatewayDB: gatewayDB,
		newAPIDB:  newAPIDB,
		config:    config,
		cache:     cache,
	}
}

// SyncUserByAPIKey 通过API Key同步用户信息
func (s *SyncService) SyncUserByAPIKey(apiKey string) (*model.User, error) {
	logger.Infof("in SyncUserByAPIKey: %v", apiKey)
	time.Sleep(3 * time.Second)
	// 从New API数据库查询token信息
	var newAPIToken model.NewAPIToken
	if err := s.newAPIDB.Where("`key` = ?", apiKey).First(&newAPIToken).Error; err != nil {
		return nil, err
	}

	// 在gateway数据库中查找或创建用户
	var user model.User
	result := s.gatewayDB.Where("api_key = ?", apiKey).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新用户
		user = model.User{
			APIKey:      apiKey,
			TokenID:     newAPIToken.ID,
			RemainQuota: newAPIToken.RemainQuota,
			UsedQuota:   newAPIToken.UsedQuota,
			ExpiredTime: newAPIToken.ExpiredTime,
			Status:      newAPIToken.Status,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.gatewayDB.Create(&user).Error; err != nil {
			return nil, err
		}
	} else if result.Error != nil {
		return nil, result.Error
	} else if user.RemainQuota != newAPIToken.RemainQuota {
		// 更新现有用户
		user.TokenID = newAPIToken.ID
		user.RemainQuota = newAPIToken.RemainQuota
		user.UsedQuota = newAPIToken.UsedQuota
		user.ExpiredTime = newAPIToken.ExpiredTime
		user.Status = newAPIToken.Status
		user.UpdatedAt = time.Now()

		if err := s.gatewayDB.Save(&user).Error; err != nil {
			return nil, err
		}
		// NewNewAPIService(s.config, s.cache).GetBillingInfo("sk-"+apiKey, false)
	}

	return &user, nil
}

// SyncLogsByTokenID 同步指定TokenID的日志
func (s *SyncService) SyncLogsByTokenID(tokenID uint, lastSyncID uint) error {
	// 查询New API数据库中新的日志
	var newAPILogs []model.NewAPILog
	query := s.newAPIDB.Where("token_id = ?", tokenID)

	if lastSyncID > 0 {
		query = query.Where("id > ?", lastSyncID)
	}

	if err := query.Order("id asc").Find(&newAPILogs).Error; err != nil {
		return err
	}

	if len(newAPILogs) == 0 {
		return nil // 没有新日志
	}

	// 查询对应的用户
	var user model.User
	if err := s.gatewayDB.Where("token_id = ?", tokenID).First(&user).Error; err != nil {
		return err
	}

	// 同步日志
	for _, newAPILog := range newAPILogs {
		// 解析Other字段，提取upstream_model_name
		var otherData map[string]interface{}
		upstreamModelName := newAPILog.ModelName

		if newAPILog.Other != "" {
			if err := json.Unmarshal([]byte(newAPILog.Other), &otherData); err == nil {
				if name, ok := otherData["upstream_model_name"].(string); ok && name != "" {
					upstreamModelName = name
				}
			}
		}

		// 创建本地日志
		log := model.Log{
			UserID:            user.ID,
			RemoteLogID:       newAPILog.ID,
			CreatedAt:         newAPILog.CreatedAt,
			Type:              newAPILog.Type,
			Content:           newAPILog.Content,
			Username:          newAPILog.Username,
			TokenName:         newAPILog.TokenName,
			ModelName:         newAPILog.ModelName,
			Quota:             newAPILog.Quota,
			PromptTokens:      newAPILog.PromptTokens,
			CompletionTokens:  newAPILog.CompletionTokens,
			UseTime:           newAPILog.UseTime,
			IsStream:          newAPILog.IsStream,
			Channel:           newAPILog.Channel,
			ChannelName:       newAPILog.ChannelName,
			TokenID:           newAPILog.TokenID,
			Group:             newAPILog.Group,
			Other:             newAPILog.Other,
			UpstreamModelName: upstreamModelName,
		}

		if err := s.gatewayDB.Create(&log).Error; err != nil {
			return err
		}
	}

	// 更新同步状态
	lastID := newAPILogs[len(newAPILogs)-1].ID
	if err := s.gatewayDB.Model(&model.SyncState{}).
		Where("token_id = ?", tokenID).
		Save(map[string]interface{}{
			"token_id":     tokenID,
			"last_sync_id": lastID,
			"updated_at":   time.Now(),
		}).Error; err != nil {
		return err
	}

	return nil
}

// SyncAllUsers 同步所有用户信息
func (s *SyncService) SyncAllUsers() error {
	// 从调用层数据库获取所有用户
	var users []model.User
	if err := s.gatewayDB.Find(&users).Error; err != nil {
		return err
	}

	// 同步每个用户的信息
	for _, user := range users {
		if _, err := s.SyncUserByAPIKey(user.APIKey); err != nil {
			log.Printf("Failed to sync user %d: %v", user.ID, err)
			// 继续同步其他用户
		}
	}

	return nil
}

// SyncAllLogs 同步所有日志
func (s *SyncService) SyncAllLogs() error {
	// 获取所有用户的TokenID
	var users []model.User
	if err := s.gatewayDB.Find(&users).Error; err != nil {
		return err
	}

	// 获取每个用户的最后同步ID
	var syncStates []model.SyncState
	if err := s.gatewayDB.Find(&syncStates).Error; err != nil {
		return err
	}

	// 创建TokenID到最后同步ID的映射
	lastSyncMap := make(map[uint]uint)
	for _, state := range syncStates {
		lastSyncMap[state.TokenID] = state.LastSyncID
	}

	// 同步每个用户的日志
	for _, user := range users {
		lastSyncID := lastSyncMap[user.TokenID]
		if err := s.SyncLogsByTokenID(user.TokenID, lastSyncID); err != nil {
			log.Printf("Failed to sync logs for token %d: %v", user.TokenID, err)
			// 继续同步其他用户的日志
		}
	}

	return nil
}
