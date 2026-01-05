package awscatalog

import (
	"errors"
	"sync"

	"github.com/aws-learning-game/internal/company"
)

var (
	ErrServiceNotFound = errors.New("service not found")
	ErrInvalidCategory = errors.New("invalid category")
	ErrInvalidScenario = errors.New("invalid scenario")
)

// catalogImpl AWS 服務目錄實作
type catalogImpl struct {
	services           map[string]*AWSService
	servicesByCategory map[string][]*AWSService
	mu                 sync.RWMutex
}

// NewAWSServiceCatalog 建立 AWS 服務目錄
func NewAWSServiceCatalog() AWSServiceCatalog {
	c := &catalogImpl{
		services:           make(map[string]*AWSService),
		servicesByCategory: make(map[string][]*AWSService),
	}
	return c
}

// NewAWSServiceCatalogWithServices 建立 AWS 服務目錄並載入服務
func NewAWSServiceCatalogWithServices(services []*AWSService) AWSServiceCatalog {
	c := &catalogImpl{
		services:           make(map[string]*AWSService),
		servicesByCategory: make(map[string][]*AWSService),
	}
	for _, s := range services {
		c.services[s.ID] = s
		c.servicesByCategory[s.Category] = append(c.servicesByCategory[s.Category], s)
	}
	return c
}

// GetService 取得服務
func (c *catalogImpl) GetService(serviceID string) (*AWSService, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	service, ok := c.services[serviceID]
	if !ok {
		return nil, ErrServiceNotFound
	}
	return service, nil
}

// GetServicesByCategory 依類別取得服務
func (c *catalogImpl) GetServicesByCategory(category string) ([]AWSService, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	services, ok := c.servicesByCategory[category]
	if !ok || len(services) == 0 {
		return []AWSService{}, nil
	}

	result := make([]AWSService, len(services))
	for i, s := range services {
		result[i] = *s
	}
	return result, nil
}

// GetRecommendedServices 取得適合的服務建議
func (c *catalogImpl) GetRecommendedServices(scenario string, comp *company.Company) ([]AWSService, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 根據情境取得相關類別
	categories := getRelevantCategories(scenario)
	if len(categories) == 0 {
		return []AWSService{}, nil
	}

	// 收集所有相關類別的服務
	var result []AWSService
	seen := make(map[string]bool)

	for _, cat := range categories {
		services, ok := c.servicesByCategory[cat]
		if !ok {
			continue
		}
		for _, s := range services {
			if !seen[s.ID] {
				seen[s.ID] = true
				// 根據公司屬性過濾服務
				if isServiceSuitableForCompany(s, comp) {
					result = append(result, *s)
				}
			}
		}
	}

	return result, nil
}

// getRelevantCategories 根據情境取得相關類別
func getRelevantCategories(scenario string) []string {
	// 先嘗試直接映射
	if categories, ok := ScenarioCategoryMapping[ScenarioType(scenario)]; ok {
		return categories
	}

	// 根據關鍵字判斷
	scenarioKeywords := map[string][]string{
		"compute":           {"compute"},
		"server":            {"compute"},
		"ec2":               {"compute"},
		"lambda":            {"compute"},
		"container":         {"compute"},
		"storage":           {"storage"},
		"s3":                {"storage"},
		"backup":            {"storage"},
		"database":          {"database"},
		"rds":               {"database"},
		"dynamodb":          {"database"},
		"network":           {"networking"},
		"vpc":               {"networking"},
		"cdn":               {"networking"},
		"security":          {"security"},
		"iam":               {"security"},
		"waf":               {"security"},
		"high_availability": {"compute", "database", "networking"},
		"disaster_recovery": {"storage", "database", "networking"},
		"cost":              {"compute", "storage"},
		"migration":         {"compute", "storage", "database"},
		"analytics":         {"analytics"},
		"machine_learning":  {"machine_learning"},
		"ml":                {"machine_learning"},
	}

	for keyword, categories := range scenarioKeywords {
		if containsKeyword(scenario, keyword) {
			return categories
		}
	}

	// 預設返回空
	return []string{}
}

// containsKeyword 檢查字串是否包含關鍵字
func containsKeyword(s, keyword string) bool {
	if len(s) < len(keyword) {
		return false
	}
	for i := 0; i <= len(s)-len(keyword); i++ {
		if s[i:i+len(keyword)] == keyword {
			return true
		}
	}
	return false
}

// isServiceSuitableForCompany 判斷服務是否適合公司
func isServiceSuitableForCompany(service *AWSService, comp *company.Company) bool {
	if comp == nil {
		return true
	}

	// 根據公司類型和屬性過濾
	switch comp.Type {
	case company.Startup:
		// 新創公司偏好低成本、易上手的服務
		return true
	case company.Traditional:
		// 傳產公司可能需要更穩定的服務
		return true
	case company.CloudReseller:
		// 雲端代理商熟悉各種服務
		return true
	case company.CloudNative:
		// 雲端公司可以使用進階服務
		return true
	}

	return true
}

// AddService 新增服務 (用於測試或動態載入)
func (c *catalogImpl) AddService(service *AWSService) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.services[service.ID] = service
	c.servicesByCategory[service.Category] = append(c.servicesByCategory[service.Category], service)
}

// GetAllServices 取得所有服務
func (c *catalogImpl) GetAllServices() []AWSService {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]AWSService, 0, len(c.services))
	for _, s := range c.services {
		result = append(result, *s)
	}
	return result
}

// GetAllCategories 取得所有類別
func (c *catalogImpl) GetAllCategories() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	categories := make([]string, 0, len(c.servicesByCategory))
	for cat := range c.servicesByCategory {
		categories = append(categories, cat)
	}
	return categories
}
