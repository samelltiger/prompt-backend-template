// internal/model/new_api_models.go
package model

// New API数据库中的tokens表模型
type NewAPIToken struct {
	ID          uint   `gorm:"primaryKey;column:id" json:"id"`
	Key         string `gorm:"column:key" json:"key"`
	Status      int    `gorm:"column:status" json:"status"`
	Name        string `gorm:"column:name" json:"name"`
	ExpiredTime int64  `gorm:"column:expired_time" json:"expired_time"`
	RemainQuota int64  `gorm:"column:remain_quota" json:"remain_quota"`
	UsedQuota   int64  `gorm:"column:used_quota" json:"used_quota"`
}

// 设置表名
func (NewAPIToken) TableName() string {
	return "tokens"
}

// New API数据库中的logs表模型
type NewAPILog struct {
	ID               uint   `gorm:"primaryKey;column:id" json:"id"`
	UserID           uint   `gorm:"column:user_id" json:"user_id"`
	CreatedAt        int64  `gorm:"column:created_at" json:"created_at"`
	Type             int    `gorm:"column:type" json:"type"`
	Content          string `gorm:"column:content" json:"content"`
	Username         string `gorm:"column:username" json:"-"`
	TokenName        string `gorm:"column:token_name" json:"-"`
	ModelName        string `gorm:"column:model_name" json:"model_name"`
	Quota            int64  `gorm:"column:quota" json:"quota"`
	PromptTokens     int    `gorm:"column:prompt_tokens" json:"prompt_tokens"`
	CompletionTokens int    `gorm:"column:completion_tokens" json:"completion_tokens"`
	UseTime          int    `gorm:"column:use_time" json:"use_time"`
	IsStream         bool   `gorm:"column:is_stream" json:"is_stream"`
	Channel          int    `gorm:"column:channel" json:"channel"`
	ChannelName      string `gorm:"column:channel_name" json:"channel_name"`
	TokenID          uint   `gorm:"column:token_id" json:"token_id"`
	Group            string `gorm:"column:group" json:"group"`
	Other            string `gorm:"column:other" json:"other"`
}

// 设置表名
func (NewAPILog) TableName() string {
	return "logs"
}
