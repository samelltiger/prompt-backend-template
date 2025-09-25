// internal/service/log_service.go
package service

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"llmapisrv/config"
	"llmapisrv/internal/model"
)

type LogService struct {
	gatewayDB *gorm.DB
	newAPIDB  *gorm.DB
	config    *config.Config
}

func NewLogService(gatewayDB, newAPIDB *gorm.DB, config *config.Config) *LogService {
	return &LogService{
		gatewayDB: gatewayDB,
		newAPIDB:  newAPIDB,
		config:    config,
	}
}

// SaveLog 保存调用日志到本地数据库
func (s *LogService) SaveLog(log *model.Log) error {
	return s.gatewayDB.Create(log).Error
}

// GetLogsByUserID 获取用户日志（从本地数据库）
func (s *LogService) GetLogsByUserID(userID uint, page, pageSize int) ([]model.Log, int64, error) {
	var logs []model.Log
	var total int64

	// 计算总数
	if err := s.gatewayDB.Model(&model.Log{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := s.gatewayDB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// CleanupOldLogs 清理旧日志（仅清理本地数据库）
func (s *LogService) CleanupOldLogs() error {
	// 计算保留期限
	retentionDays := s.config.Log.RetentionDays
	if retentionDays <= 0 {
		retentionDays = 30 // 默认30天
	}

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays).Unix()

	// 删除旧日志
	return s.gatewayDB.Where("created_at < ?", cutoffTime).Delete(&model.Log{}).Error
}

// ProcessLogFromQueue 从队列处理日志
func (s *LogService) ProcessLogFromQueue(data []byte, syncSrv *SyncService) error {
	var logData map[string]interface{}
	if err := json.Unmarshal(data, &logData); err != nil {
		return err
	}

	// 获取必要字段
	apiKey, _ := logData["api_key"].(string)
	// modelName, _ := logData["model"].(string)
	// usage, _ := logData["usage"].(map[string]interface{})
	// duration, _ := logData["duration"].(float64)

	// 获取用户ID
	var user model.User
	if err := s.gatewayDB.Where("api_key = ?", apiKey).First(&user).Error; err != nil {
		return err
	}

	var syncState model.SyncState
	if err := s.gatewayDB.Model(&model.SyncState{}).
		Where("token_id = ?", user.TokenID).
		First(&syncState).Error; err != nil {
		return err
	}

	// 同步日志
	syncSrv.SyncLogsByTokenID(user.TokenID, uint(syncState.LastSyncID))

	// 同步额度
	syncSrv.SyncUserByAPIKey(apiKey)

	return nil
}

// GetLatestRemoteLogID 获取最新的远程日志ID
func (s *LogService) GetLatestRemoteLogID(tokenID uint) (uint, error) {
	var syncState model.SyncState
	err := s.gatewayDB.Where("token_id = ?", tokenID).First(&syncState).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}

	return syncState.LastSyncID, nil
}
