// internal/service/user_service.go
package service

import (
	"time"

	"gorm.io/gorm"

	"llmapisrv/internal/model"
	"llmapisrv/pkg/cache"
)

type UserService struct {
	gatewayDB   *gorm.DB
	newAPIDB    *gorm.DB
	cache       *cache.RedisCache
	syncService *SyncService
}

func NewUserService(gatewayDB, newAPIDB *gorm.DB, cache *cache.RedisCache, syncService *SyncService) *UserService {
	return &UserService{
		gatewayDB:   gatewayDB,
		newAPIDB:    newAPIDB,
		cache:       cache,
		syncService: syncService,
	}
}

// GetUserByAPIKey 通过API Key获取用户
func (s *UserService) GetUserByAPIKey(apiKey string) (*model.User, error) {
	// 先从缓存获取
	// cacheKey := "user:api_key:" + apiKey
	var user model.User

	// 尝试从本地数据库获取
	err := s.gatewayDB.Where("api_key = ?", apiKey).First(&user).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		// 如果本地没有，尝试从New API同步
		return s.syncService.SyncUserByAPIKey(apiKey)
	}

	// 如果本地有，但可能已过时，尝试更新
	if user.UpdatedAt.Before(time.Now().Add(-5 * time.Minute)) {
		// 从New API同步最新数据
		syncedUser, err := s.syncService.SyncUserByAPIKey(apiKey)
		if err != nil {
			// 如果同步失败，仍使用本地数据
			return &user, nil
		}
		return syncedUser, nil
	}

	return &user, nil
}

// GetUserByID 通过ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.gatewayDB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// AddQuota 添加用户额度
func (s *UserService) AddQuota(userID uint, quota int64) error {
	// 在本地数据库中添加额度
	txg := s.gatewayDB.Begin()
	txn := s.newAPIDB.Begin()

	// 更新本地用户额度
	if err := txg.Model(&model.User{}).
		Where("id = ?", userID).
		Update("remain_quota", gorm.Expr("remain_quota + ?", quota)).
		Error; err != nil {
		txg.Rollback()
		return err
	}

	// 获取用户信息以获取TokenID
	var user model.User
	if err := txg.First(&user, userID).Error; err != nil {
		txg.Rollback()
		return err
	}

	// 更新New API数据库中的额度
	if err := txn.Exec(`
        UPDATE tokens SET remain_quota = remain_quota + ?
        WHERE id = ?
    `, quota, user.TokenID).Error; err != nil {
		txn.Rollback()
		txg.Rollback()
		return err
	}

	err := txn.Commit().Error
	if err != nil {
		txn.Rollback()
		txg.Rollback()
		return err
	}
	txg.Commit()

	return nil
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(userID uint, status int) error {
	// 在本地数据库中更新状态
	tx := s.gatewayDB.Begin()

	// 更新本地用户状态
	if err := tx.Model(&model.User{}).
		Where("id = ?", userID).
		Update("status", status).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	// 获取用户信息以获取TokenID
	var user model.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新New API数据库中的状态
	if err := tx.Exec(`
        UPDATE tokens SET status = ?
        WHERE id = ?
    `, status, user.TokenID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
