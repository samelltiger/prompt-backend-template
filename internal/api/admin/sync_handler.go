// internal/api/admin/sync_handler.go
package admin

import (
	"github.com/gin-gonic/gin"

	"llmapisrv/internal/service"
	"llmapisrv/pkg/util"
)

type SyncRequest struct {
	APIKey string `json:"api_key"` // 可选，如果提供则仅同步指定用户
}

type SyncHandler struct {
	syncService *service.SyncService
}

func NewSyncHandler(syncService *service.SyncService) *SyncHandler {
	return &SyncHandler{
		syncService: syncService,
	}
}

// SyncUser 同步用户信息
func (h *SyncHandler) SyncUser(c *gin.Context) {
	var req SyncRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.APIKey != "" {
		// 同步指定用户
		user, err := h.syncService.SyncUserByAPIKey(util.RemoveStartSk(req.APIKey))
		if err != nil {
			util.Fail(c, util.FailCode, err.Error())
			return
		}

		util.Success(c, user)
		return
	}

	// 同步所有用户
	if err := h.syncService.SyncAllUsers(); err != nil {
		util.Fail(c, util.FailCode, err.Error())
		return
	}

	util.Success(c, "All users synced successfully")
}

// SyncLogs 同步日志
func (h *SyncHandler) SyncLogs(c *gin.Context) {
	var req SyncRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.APIKey != "" {
		// 同步指定用户的日志
		user, err := h.syncService.SyncUserByAPIKey(util.RemoveStartSk(req.APIKey))
		if err != nil {
			util.Fail(c, util.FailCode, err.Error())
			return
		}

		if err := h.syncService.SyncLogsByTokenID(user.TokenID, 0); err != nil {
			util.Fail(c, util.FailCode, err.Error())
			return
		}

		util.Success(c, "Logs for user synced successfully")
		return
	}

	// 同步所有日志
	if err := h.syncService.SyncAllLogs(); err != nil {
		util.Fail(c, util.FailCode, err.Error())
		return
	}

	util.Success(c, "All logs synced successfully")
}

// SyncAll 同步所有数据
func (h *SyncHandler) SyncAll(c *gin.Context) {
	// 先同步所有用户
	if err := h.syncService.SyncAllUsers(); err != nil {
		util.Fail(c, util.FailCode, "Failed to sync logs: "+err.Error())
		return
	}

	// 再同步所有日志
	if err := h.syncService.SyncAllLogs(); err != nil {
		util.Fail(c, util.FailCode, "Failed to sync logs: "+err.Error())
		return
	}

	util.Success(c, "All data synced successfully")
}
