// internal/api/admin/redemption.go
package admin

import (
	"net/http"

	"llmapisrv/internal/service"
	"llmapisrv/pkg/util"

	"github.com/gin-gonic/gin"
)

type GenerateCodesRequest struct {
	Count    int    `json:"count" binding:"required,min=1,max=1000"`
	Quota    int64  `json:"quota" binding:"required,min=1"`
	BatchNum string `json:"batch_num"` // 批次号，可选
}

type AddQuotaRequest struct {
	APIKey string `json:"api_key" binding:"required"`
	Quota  int64  `json:"quota" binding:"required,min=1"`
}

type RedemptionAdminHandler struct {
	redemptionService *service.RedemptionService
	userService       *service.UserService
}

func NewRedemptionAdminHandler(redemptionService *service.RedemptionService, userService *service.UserService) *RedemptionAdminHandler {
	return &RedemptionAdminHandler{
		redemptionService: redemptionService,
		userService:       userService,
	}
}

// GenerateCodes 生成兑换码
func (h *RedemptionAdminHandler) GenerateCodes(c *gin.Context) {
	var req GenerateCodesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成兑换码
	codes, err := h.redemptionService.GenerateCodes(req.Count, req.Quota, req.BatchNum)
	if err != nil {
		util.ParamError(c, err.Error())
		return
	}

	util.Success(c, gin.H{
		"codes":  codes,
		"count":  len(codes),
		"quota":  req.Quota,
		"amount": float64(req.Quota) / 100.0, // 转换为美元显示
	})
}

// AddQuota 直接添加额度
func (h *RedemptionAdminHandler) AddQuota(c *gin.Context) {
	var req AddQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ParamError(c, err.Error())
		return
	}

	// 查找用户
	user, err := h.userService.GetUserByAPIKey(util.RemoveStartSk(req.APIKey))
	if err != nil {
		util.ParamError(c, "Invalid API key")
		return
	}

	// 添加额度
	if err := h.userService.AddQuota(user.ID, req.Quota); err != nil {
		util.Fail(c, util.FailCode, "Failed to add quota")
		return
	}

	util.Success(c, gin.H{
		"api_key":        req.APIKey,
		"quota_added":    req.Quota,
		"amount":         req.Quota, // 转换为美元显示
		"current_quota":  user.RemainQuota + req.Quota,
		"current_amount": user.RemainQuota + req.Quota,
	})
}
