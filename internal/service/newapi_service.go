// internal/service/newapi_service.go
package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"llmapisrv/config"
	"llmapisrv/pkg/cache"
	"llmapisrv/pkg/logger"
)

type NewAPIService struct {
	client *http.Client
	config *config.Config
	cache  *cache.RedisCache
}

func NewNewAPIService(config *config.Config, cache *cache.RedisCache) *NewAPIService {
	return &NewAPIService{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		config: config,
		cache:  cache,
	}
}

// GetBillingInfo 获取账单信息
func (s *NewAPIService) GetBillingInfo(apiKey string, useCache bool) (map[string]interface{}, error) {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("billing:%s", apiKey)
	if data, err := s.cache.Get(cacheKey); err == nil && useCache {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(data), &result); err == nil {
			return result, nil
		}
	}

	// 从API获取
	url := fmt.Sprintf("%s/v1/dashboard/billing/subscription", s.config.NewAPI.Domain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// 保存到缓存，有效期5分钟
	s.cache.Set(cacheKey, string(body), 5*60)

	return result, nil
}

// GetModelPricing 获取模型价格
func (s *NewAPIService) GetModelPricing() (map[string]interface{}, error) {
	// 先从缓存获取
	cacheKey := "model:pricing"
	if data, err := s.cache.Get(cacheKey); err == nil {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(data), &result); err == nil {
			return result, nil
		}
	}

	// 从API获取
	url := fmt.Sprintf("%s/api/pricing", s.config.NewAPI.Domain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// 保存到缓存，有效期1小时
	s.cache.Set(cacheKey, string(body), 3600)

	return result, nil
}

// 转发聊天完成请求
func (s *NewAPIService) ChatCompletion(apiKey string, requestBody map[string]interface{}) (*http.Response, error) {
	// 获取请求的模型
	modelName, ok := requestBody["model"].(string)
	if !ok {
		return nil, fmt.Errorf("missing model parameter")
	}

	// 查找可用的实际模型
	actualModel, err := s.getAvailableModel(modelName)
	if err != nil {
		return nil, err
	}

	// 替换模型名称
	requestBody["model"] = actualModel

	// 构建请求
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/v1/chat/completions", s.config.NewAPI.Domain)
	logger.Infof("ChatCompletion url: %v, requestBody: %v", url, string(jsonData))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return s.client.Do(req)
}

// 获取可用模型
func (s *NewAPIService) getAvailableModel(modelName string) (string, error) {
	// 从映射配置中查找
	models, ok := s.config.ModelMapping[modelName]
	if !ok {
		return "", fmt.Errorf("model not supported: %s", modelName)
	}

	// 检查缓存中的可用模型状态
	for _, model := range models {
		cacheKey := fmt.Sprintf("model:status:%s", model)
		status, err := s.cache.Get(cacheKey)
		if err == nil && status == "available" {
			return model, nil
		}
	}

	// 如果没有可用模型，随机选择一个
	// 实际应用中应该有更复杂的选择策略
	return models[0], nil
}

// 更新模型状态
func (s *NewAPIService) UpdateModelStatus(modelName string, available bool) {
	cacheKey := fmt.Sprintf("model:status:%s", modelName)
	status := "unavailable"
	if available {
		status = "available"
	}
	s.cache.Set(cacheKey, status, 3600) // 1小时有效期
}
