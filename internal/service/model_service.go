// internal/service/model_service.go
package service

import (
	"fmt"

	"gorm.io/gorm"

	"llmapisrv/config"
)

type ModelService struct {
	gatewayDB *gorm.DB
	newAPIDB  *gorm.DB
	config    *config.Config
}

func NewModelService(gatewayDB, newAPIDB *gorm.DB, config *config.Config) *ModelService {
	return &ModelService{
		gatewayDB: gatewayDB,
		newAPIDB:  newAPIDB,
		config:    config,
	}
}

// MapModels 映射模型信息
func (s *ModelService) MapModels(models []interface{}) []interface{} {
	// 创建映射表
	modelMap := make(map[string][]string)
	for k, v := range s.config.ModelMapping {
		modelMap[k] = v
	}

	// 创建反向映射表
	reverseMap := make(map[string]string)
	for displayName, actualModels := range modelMap {
		for _, actualModel := range actualModels {
			reverseMap[actualModel] = displayName
		}
	}

	// 创建结果集
	var result []interface{}
	modelPriceMap := make(map[string]map[string]interface{})

	// 处理每个模型
	for _, m := range models {
		model, ok := m.(map[string]interface{})
		if !ok {
			continue
		}

		modelName, ok := model["model_name"].(string)
		if !ok {
			continue
		}

		// 如果是映射模型，替换为显示名
		displayName, exists := reverseMap[modelName]
		if exists {
			// 存储价格信息到映射表
			modelPriceMap[displayName] = model
			continue
		}
		// else {
		// 	// 如果不在映射表中，直接添加到结果
		// 	result = append(result, model)
		// }
	}

	// 添加映射模型到结果
	for displayName, priceInfo := range modelPriceMap {
		priceInfo["model_name"] = displayName
		result = append(result, priceInfo)
	}

	return result
}

// CalculateQuota 计算使用额度
func (s *ModelService) CalculateQuota(modelName string, promptTokens, completionTokens int) (int64, error) {
	// 获取模型价格信息
	newAPIService := NewNewAPIService(s.config, nil)
	pricing, err := newAPIService.GetModelPricing()
	if err != nil {
		return 0, err
	}

	// 解析模型数据
	models, ok := pricing["data"].([]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid pricing data")
	}

	// 查找实际模型
	var actualModel string
	for displayName, actualModels := range s.config.ModelMapping {
		if displayName == modelName && len(actualModels) > 0 {
			actualModel = actualModels[0]
			break
		}
	}

	if actualModel == "" {
		actualModel = modelName // 如果没有映射，使用原始名称
	}

	// 查找模型价格信息
	var modelRatio float64 = 1.0
	var completionRatio float64 = 1.0

	for _, m := range models {
		model, ok := m.(map[string]interface{})
		if !ok {
			continue
		}

		name, ok := model["model_name"].(string)
		if !ok || name != actualModel {
			continue
		}

		// 获取价格比率
		if ratio, ok := model["model_ratio"].(float64); ok {
			modelRatio = ratio
		}

		if ratio, ok := model["completion_ratio"].(float64); ok {
			completionRatio = ratio
		}

		break
	}

	// 计算价格 (按照每百万tokens的价格)
	// 输入价格: 2 * modelRatio 每1M tokens
	// 输出价格: 2 * modelRatio * completionRatio 每1M tokens
	inputCost := float64(promptTokens) / 1000000.0 * 2.0 * modelRatio
	outputCost := float64(completionTokens) / 1000000.0 * 2.0 * modelRatio * completionRatio

	// 总价格（单位：0.001美元）
	totalCost := (inputCost + outputCost) * 1000.0

	return int64(totalCost), nil
}

// 检查模型状态并更新
func (s *ModelService) CheckModelStatus() {
	// 获取所有需要检查的模型
	allModels := make([]string, 0)
	for _, models := range s.config.ModelMapping {
		allModels = append(allModels, models...)
	}

	newAPIService := NewNewAPIService(s.config, nil)

	// 对每个模型进行状态检查
	for _, model := range allModels {
		// 这里应该实现对模型的可用性检查
		// 例如，发送一个简单请求来测试模型是否可用
		available := true // 默认假设模型可用

		// 更新模型状态
		newAPIService.UpdateModelStatus(model, available)
	}
}
