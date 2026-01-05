package awscatalog

import "github.com/aws-learning-game/internal/company"

// AWSServiceCatalog AWS 服務目錄介面
type AWSServiceCatalog interface {
	// GetService 取得服務
	GetService(serviceID string) (*AWSService, error)
	// GetServicesByCategory 依類別取得服務
	GetServicesByCategory(category string) ([]AWSService, error)
	// GetRecommendedServices 取得適合的服務建議
	GetRecommendedServices(scenario string, comp *company.Company) ([]AWSService, error)
}
