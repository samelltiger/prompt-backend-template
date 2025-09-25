// pkg/cron/cron.go
package cron

import (
	"log"

	"github.com/robfig/cron/v3"

	"llmapisrv/internal/service"
)

// CronManager 定时任务管理器
type CronManager struct {
	cron         *cron.Cron
	logService   *service.LogService
	modelService *service.ModelService
	syncService  *service.SyncService
}

// NewCronManager 创建定时任务管理器
func NewCronManager(
	logService *service.LogService,
	modelService *service.ModelService,
	syncService *service.SyncService,
) *CronManager {
	c := cron.New(cron.WithSeconds())
	return &CronManager{
		cron:         c,
		logService:   logService,
		modelService: modelService,
		syncService:  syncService,
	}
}

// Start 启动定时任务
func (m *CronManager) Start() {
	// 每天凌晨3点清理旧日志
	_, err := m.cron.AddFunc("0 0 3 * * *", m.cleanupOldLogs)
	if err != nil {
		log.Printf("Failed to add cleanup logs task: %v", err)
	}

	// 每5分钟检查一次模型状态
	_, err = m.cron.AddFunc("0 */5 * * * *", m.checkModelStatus)
	if err != nil {
		log.Printf("Failed to add check model status task: %v", err)
	}

	// 每10分钟同步一次用户信息
	_, err = m.cron.AddFunc("0 */10 * * * *", m.syncUsers)
	if err != nil {
		log.Printf("Failed to add sync users task: %v", err)
	}

	// 每5分钟同步一次日志
	_, err = m.cron.AddFunc("0 */5 * * * *", m.syncLogs)
	if err != nil {
		log.Printf("Failed to add sync logs task: %v", err)
	}

	m.cron.Start()
}

// Stop 停止定时任务
func (m *CronManager) Stop() {
	m.cron.Stop()
}

// 清理旧日志
func (m *CronManager) cleanupOldLogs() {
	log.Println("Starting cleanup of old logs")
	if err := m.logService.CleanupOldLogs(); err != nil {
		log.Printf("Error cleaning up old logs: %v", err)
	}
	log.Println("Finished cleanup of old logs")
}

// 检查模型状态
func (m *CronManager) checkModelStatus() {
	log.Println("Starting model status check")
	// 这里需要实现检查模型状态的逻辑
	// 可以通过调用 modelService 中的方法来检查每个模型的可用性
	log.Println("Finished model status check")
}

// 同步用户信息
func (m *CronManager) syncUsers() {
	log.Println("Starting user sync")
	if err := m.syncService.SyncAllUsers(); err != nil {
		log.Printf("Error syncing users: %v", err)
	}
	log.Println("Finished user sync")
}

// 同步日志
func (m *CronManager) syncLogs() {
	log.Println("Starting logs sync")
	if err := m.syncService.SyncAllLogs(); err != nil {
		log.Printf("Error syncing logs: %v", err)
	}
	log.Println("Finished logs sync")
}
