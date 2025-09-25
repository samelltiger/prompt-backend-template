// internal/service/redemption_service.go
package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"llmapisrv/internal/model"
)

type RedemptionService struct {
	db *gorm.DB
}

func NewRedemptionService(db *gorm.DB) *RedemptionService {
	return &RedemptionService{
		db: db,
	}
}

// GenerateCodes 生成兑换码
func (s *RedemptionService) GenerateCodes(count int, quota int64, batchNum string) ([]string, error) {
	codes := make([]string, 0, count)

	// 生成批次号
	if batchNum == "" {
		batchNum = fmt.Sprintf("B%d", time.Now().Unix())
	}

	// 批量生成兑换码
	for i := 0; i < count; i++ {
		// 生成随机字符串
		b := make([]byte, 12) // 16字节 -> 24字符的base64
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}

		code := fmt.Sprintf("%s-%s", batchNum, base64.URLEncoding.EncodeToString(b)[:16])
		code = strings.ReplaceAll(code, "-", "") // 移除可能的连字符
		code = strings.ReplaceAll(code, "_", "") // 移除下划线

		// 添加连字符使其更易读
		formattedCode := fmt.Sprintf("RC-%s-%s-%s",
			code[:4], code[4:8], code[8:])

		// 保存到数据库
		redemptionCode := model.RedemptionCode{
			Code:      formattedCode,
			Quota:     quota,
			Used:      false,
			CreatedAt: time.Now(),
		}

		if err := s.db.Create(&redemptionCode).Error; err != nil {
			return nil, err
		}

		codes = append(codes, formattedCode)
	}

	return codes, nil
}

// RedeemCode 兑换码兑换
func (s *RedemptionService) RedeemCode(code string, userID uint) (int64, error) {
	var redemptionCode model.RedemptionCode

	// 查找兑换码
	if err := s.db.Where("code = ? AND used = ?", code, false).First(&redemptionCode).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("兑换码无效或已被使用")
		}
		return 0, err
	}

	// 开启事务
	tx := s.db.Begin()

	// 标记兑换码为已使用
	if err := tx.Model(&redemptionCode).Updates(map[string]interface{}{
		"used":    true,
		"used_at": time.Now(),
		"used_by": userID,
	}).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 记录兑换日志
	redemptionLog := model.RedemptionLog{
		UserID:    userID,
		Code:      code,
		Quota:     redemptionCode.Quota,
		CreatedAt: time.Now(),
	}

	if err := tx.Create(&redemptionLog).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return redemptionCode.Quota, nil
}

// RedeemCode 兑换码兑换
func (s *RedemptionService) RedeemCodeInfo(code string, userID uint) (int64, error) {
	var redemptionCode model.RedemptionCode

	// 查找兑换码
	if err := s.db.Where("code = ? AND used = ?", code, false).First(&redemptionCode).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("兑换码无效或已被使用")
		}
		return 0, err
	}

	return redemptionCode.Quota, nil
}
