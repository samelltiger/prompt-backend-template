// internal/api/pricing.go
package api

import (
	"llmapisrv/internal/service"
	"llmapisrv/pkg/logger"
	"llmapisrv/pkg/util"

	"github.com/gin-gonic/gin"
)

type PricingHandler struct {
	newAPIService *service.NewAPIService
	modelService  *service.ModelService
}

func NewPricingHandler(newAPIService *service.NewAPIService, modelService *service.ModelService) *PricingHandler {
	return &PricingHandler{
		newAPIService: newAPIService,
		modelService:  modelService,
	}
}

/**
{
    "code": 200,
    "message": "success",
    "data": {
        "data": [
            {
                "completion_ratio": 1,
                "enable_groups": [
                    "default"
                ],
                "model_name": "claude37",
                "model_price": 0,
                "model_ratio": 37.5,
                "owner_by": "",
                "quota_type": 0
            },
			....
        ],
        "group_ratio": {
            "default": 1,
            "vip": 1
        },
        "success": true,
        "usable_group": {
            "": "用户分组",
            "default": "默认分组",
            "vip": "vip分组"
        }
    }
}
*/
// GetPricing 获取价格信息
func (h *PricingHandler) GetPricing(c *gin.Context) {
	// 从上游获取价格
	pricing, err := h.newAPIService.GetModelPricing()
	if err != nil {
		util.ServerError(c, err)
		return
	}

	logger.Infof("newAPIService data: %v", util.ToJSONString(pricing))
	// 处理模型映射
	if models, ok := pricing["data"].([]interface{}); ok {
		mappedModels := h.modelService.MapModels(models)
		pricing["data"] = mappedModels
	}

	util.Success(c, pricing)
}
