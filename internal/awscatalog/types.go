package awscatalog

// AWSService AWS 服務
type AWSService struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Category      string       `json:"category"` // "compute", "storage", "database", "networking", "security"
	Description   string       `json:"description"`
	UseCases      []string     `json:"use_cases"`
	PricingModel  PricingModel `json:"pricing_model"`
	SAAExamTopics []string     `json:"saa_exam_topics"` // SAA 考試相關主題
	BestPractices []string     `json:"best_practices"`
}

// PricingModel 定價模型
type PricingModel struct {
	Type      string  `json:"type"` // "on_demand", "reserved", "spot", "savings_plan"
	BasePrice float64 `json:"base_price"`
	Unit      string  `json:"unit"` // "hour", "GB", "request"
	FreeTier  string  `json:"free_tier"`
}

// ServiceCategory 服務類別
type ServiceCategory string

const (
	CategoryCompute    ServiceCategory = "compute"
	CategoryStorage    ServiceCategory = "storage"
	CategoryDatabase   ServiceCategory = "database"
	CategoryNetworking ServiceCategory = "networking"
	CategorySecurity   ServiceCategory = "security"
	CategoryAnalytics  ServiceCategory = "analytics"
	CategoryML         ServiceCategory = "machine_learning"
	CategoryManagement ServiceCategory = "management"
)

// ScenarioType 情境類型
type ScenarioType string

const (
	ScenarioCompute       ScenarioType = "compute"
	ScenarioStorage       ScenarioType = "storage"
	ScenarioDatabase      ScenarioType = "database"
	ScenarioNetworking    ScenarioType = "networking"
	ScenarioSecurity      ScenarioType = "security"
	ScenarioHighAvail     ScenarioType = "high_availability"
	ScenarioDisasterRecov ScenarioType = "disaster_recovery"
	ScenarioCostOptim     ScenarioType = "cost_optimization"
	ScenarioMigration     ScenarioType = "migration"
)

// ScenarioCategoryMapping 情境到類別的映射
var ScenarioCategoryMapping = map[ScenarioType][]string{
	ScenarioCompute:       {"compute"},
	ScenarioStorage:       {"storage"},
	ScenarioDatabase:      {"database"},
	ScenarioNetworking:    {"networking"},
	ScenarioSecurity:      {"security"},
	ScenarioHighAvail:     {"compute", "database", "networking"},
	ScenarioDisasterRecov: {"storage", "database", "networking"},
	ScenarioCostOptim:     {"compute", "storage"},
	ScenarioMigration:     {"compute", "storage", "database"},
}
