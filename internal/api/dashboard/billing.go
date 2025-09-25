// internal/api/dashboard/billing.go
package dashboard

import (
	"llmapisrv/internal/service"
	"llmapisrv/pkg/logger"
	"llmapisrv/pkg/util"

	"github.com/gin-gonic/gin"
)

type BillingHandler struct {
	newAPIService *service.NewAPIService
	userService   *service.UserService
}

func NewBillingHandler(newAPIService *service.NewAPIService, userService *service.UserService) *BillingHandler {
	return &BillingHandler{
		newAPIService: newAPIService,
		userService:   userService,
	}
}

/**
{
    "code": 200,
    "message": "success",
    "data": {
        "access_until": 1746278367,  // 到期时间，0为永不过期
        "hard_limit_usd": 11,//总额度
        "has_payment_method": true,
        "object": "billing_subscription",
        "soft_limit_usd": 11,//总额度
        "system_hard_limit_usd": 11//总额度
    }
}
*/
// GetSubscription 获取订阅信息
func (h *BillingHandler) GetSubscription(c *gin.Context) {
	apiKey := c.GetString("api_key")

	// 获取账单信息
	billing, err := h.newAPIService.GetBillingInfo(apiKey, true)
	if err != nil {
		util.ServerError(c, err)
		return
	}

	util.Success(c, billing)
}

/**
{
    "code": 200,
    "message": "success",
    "data": {
        "object": "list",
        "total_usage": 4928475 // 已使用的额度， 4928475/500000 后为美元额度
    }
}
*/
// GetUsage 获取使用情况
func (h *BillingHandler) GetUsage(c *gin.Context) {
	userID := c.GetUint("user_id")
	logger.Infof("GetUsage userID: %v", userID)

	// 获取用户信息
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		util.ServerError(c, err)
		return
	}

	// 构建响应
	response := map[string]interface{}{
		"object":      "list",
		"total_usage": user.UsedQuota * 100 / 500000, // 转换为美元（为了避免浮点数问题，先将其乘以100倍）
	}

	util.Success(c, response)
}
