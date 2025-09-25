// internal/api/status.go
package api

import (
	"llmapisrv/internal/service"
	"llmapisrv/pkg/util"

	"github.com/gin-gonic/gin"
)

type StatusHandler struct {
	newAPIService *service.NewAPIService
}

func NewStatusHandler(newAPIService *service.NewAPIService) *StatusHandler {
	return &StatusHandler{
		newAPIService: newAPIService,
	}
}

// HealthCheck 健康检查
func (h *StatusHandler) HealthCheck(c *gin.Context) {
	// 调用上游API的健康检查
	// err := h.newAPIService.CheckHealth()
	err := error(nil)
	if err != nil {
		util.Fail(c, util.FailCode, err.Error())
		return
	}

	util.Success(c, "")
}
