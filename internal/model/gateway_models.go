// internal/model/gateway_models.go
package model

import (
	"time"
)

// 调用层数据库中的用户表
type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	APIKey      string    `gorm:"column:api_key;uniqueIndex" json:"api_key"`
	TokenID     uint      `gorm:"column:token_id" json:"token_id"`
	RemainQuota int64     `gorm:"column:remain_quota" json:"remain_quota"` // 剩余额度（单位：0.001美元）
	UsedQuota   int64     `gorm:"column:used_quota" json:"used_quota"`     // 已用额度
	ExpiredTime int64     `gorm:"column:expired_time" json:"expired_time"` // 过期时间戳
	Status      int       `gorm:"column:status" json:"status"`             // 状态：1正常，0禁用
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// 调用层数据库中的日志表
type Log struct {
	ID                uint   `gorm:"primaryKey" json:"id"`
	UserID            uint   `gorm:"column:user_id;index" json:"user_id"`
	RemoteLogID       uint   `gorm:"column:remote_log_id;uniqueIndex" json:"remote_log_id"` // New API 中的日志ID
	CreatedAt         int64  `gorm:"column:created_at" json:"created_at"`
	Type              int    `gorm:"column:type" json:"type"`
	Content           string `gorm:"column:content;type:text" json:"content"`
	Username          string `gorm:"column:username" json:"-"`
	TokenName         string `gorm:"column:token_name" json:"-"`
	ModelName         string `gorm:"column:model_name" json:"model_name"`
	Quota             int64  `gorm:"column:quota" json:"quota"` //  1美元为 500000 token，充值都是以token计算
	PromptTokens      int    `gorm:"column:prompt_tokens" json:"prompt_tokens"`
	CompletionTokens  int    `gorm:"column:completion_tokens" json:"completion_tokens"`
	UseTime           int    `gorm:"column:use_time" json:"use_time"`
	IsStream          bool   `gorm:"column:is_stream" json:"is_stream"`
	Channel           int    `gorm:"column:channel" json:"channel"`
	ChannelName       string `gorm:"column:channel_name" json:"channel_name"`
	TokenID           uint   `gorm:"column:token_id" json:"token_id"`
	Group             string `gorm:"column:group" json:"group"`
	Other             string `gorm:"column:other;type:text" json:"other"`
	UpstreamModelName string `gorm:"column:upstream_model_name" json:"upstream_model_name"` // 真实模型名
}

// 兑换码表
type RedemptionCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"column:code;uniqueIndex" json:"code"`
	Quota     int64     `gorm:"column:quota" json:"quota"` // 额度（单位：0.001美元）
	Used      bool      `gorm:"column:used" json:"used"`   // 是否已使用
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UsedAt    time.Time `gorm:"column:used_at" json:"used_at"`
	UsedBy    uint      `gorm:"column:used_by" json:"used_by"` // 使用者ID
}

// 兑换记录表
type RedemptionLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"column:user_id;index" json:"user_id"`
	Code      string    `gorm:"column:code" json:"code"`
	Quota     int64     `gorm:"column:quota" json:"quota"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// 同步状态表
type SyncState struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TokenID    uint      `gorm:"column:token_id;uniqueIndex" json:"token_id"`
	LastSyncID uint      `gorm:"column:last_sync_id" json:"last_sync_id"` // 最后同步的日志ID
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}
