// internal/api/redemption.go
package api

import (
	"net/http"

	"llmapisrv/internal/middleware"
	"llmapisrv/internal/service"
	"llmapisrv/pkg/util"

	"github.com/gin-gonic/gin"
)

type RedemptionRequest struct {
	Code string `json:"code" binding:"required"`
}

type RedemptionHandler struct {
	newAPIService     *service.NewAPIService
	redemptionService *service.RedemptionService
	userService       *service.UserService
}

func NewRedemptionHandler(
	newAPIService *service.NewAPIService,
	redemptionService *service.RedemptionService, userService *service.UserService) *RedemptionHandler {
	return &RedemptionHandler{
		newAPIService:     newAPIService,
		redemptionService: redemptionService,
		userService:       userService,
	}
}

// RedeemCode 兑换码兑换
func (h *RedemptionHandler) RedeemCode(c *gin.Context) {
	var req RedemptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ParamError(c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clientInfoInterface, exists := c.Get("client_info")
	if !exists {
		util.Fail(c, util.FailCode, "client_info is not exists")
		return
	}
	clientInfo := clientInfoInterface.(*middleware.ClientInfo)

	userID := c.GetUint("user_id")

	// 兑换码验证与使用
	quota, err := h.redemptionService.RedeemCode(req.Code, userID)
	if err != nil {
		util.Fail(c, util.FailCode, err.Error())
		return
	}

	// 更新用户额度
	if err := h.userService.AddQuota(userID, quota); err != nil {
		util.Fail(c, util.FailCode, "Failed to add quota")
		return
	}
	h.newAPIService.GetBillingInfo(clientInfo.AuthNoSk, false)
	util.Success(c, gin.H{
		"quota":  quota,
		"amount": int(float64(quota) / 500000), // 转换为美元显示
	})
}

// RedeemCodeInfo 兑换码额度查询
func (h *RedemptionHandler) RedeemCodeInfo(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		util.ParamError(c, "code is required")
		return
	}
	userID := c.GetUint("user_id")

	// 兑换码验证与使用
	quota, err := h.redemptionService.RedeemCodeInfo(code, userID)
	if err != nil {
		util.Fail(c, util.FailCode, err.Error())
		return
	}

	util.Success(c, gin.H{
		"quota":  quota,
		"amount": int(float64(quota) / 500000), // 转换为美元显示
	})
}
