// pkg/database/mysql.go
package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"llmapisrv/config"
)

// 数据库连接
type DBConnections struct {
	GatewayDB *gorm.DB // 调用层数据库
	NewAPIDB  *gorm.DB // New API 原始数据库
}

// InitDatabases 初始化数据库连接
func InitDatabases(cfg *config.Config) (*DBConnections, error) {
	// 设置GORM配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接调用层数据库
	gatewayDB, err := gorm.Open(mysql.Open(cfg.Database.GatewayDSN), gormConfig)
	if err != nil {
		return nil, err
	}

	// 连接New API原始数据库（只读访问）
	newAPIDB, err := gorm.Open(mysql.Open(cfg.Database.NewAPIDSN), gormConfig)
	if err != nil {
		return nil, err
	}

	return &DBConnections{
		GatewayDB: gatewayDB,
		NewAPIDB:  newAPIDB,
	}, nil
}
