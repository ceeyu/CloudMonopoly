package event

// EventType 事件類型
type EventType string

const (
	EventOpportunity EventType = "opportunity" // 機會事件
	EventFate        EventType = "fate"        // 命運事件
	EventChallenge   EventType = "challenge"   // 關卡事件
	EventSecurity    EventType = "security"    // 資安事件
)

// Event 事件
type Event struct {
	ID           string        `json:"id"`
	Type         EventType     `json:"type"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`    // 事件描述 (真實情境)
	RealWorldRef string        `json:"real_world_ref"` // 真實案例參考
	Context      EventContext  `json:"context"`        // 事件背景
	Choices      []EventChoice `json:"choices"`        // 可選方案
	AWSTopics    []string      `json:"aws_topics"`     // 相關 AWS SAA 考點
}

// EventContext 事件背景
type EventContext struct {
	Scenario       string   `json:"scenario"`        // 情境說明
	BusinessImpact string   `json:"business_impact"` // 商業影響
	TechnicalNeeds []string `json:"technical_needs"` // 技術需求
	Constraints    []string `json:"constraints"`     // 限制條件
}

// EventChoice 事件選項
type EventChoice struct {
	ID                  int                `json:"id"`
	Title               string             `json:"title"`
	Description         string             `json:"description"`
	IsAWS               bool               `json:"is_aws"`           // 是否為 AWS 方案
	AWSServices         []string           `json:"aws_services"`     // 使用的 AWS 服務
	OnPremSolution      string             `json:"on_prem_solution"` // 地端方案描述
	Requirements        ChoiceRequirements `json:"requirements"`
	Outcomes            ChoiceOutcomes     `json:"outcomes"`
	ArchitectureDiagram string             `json:"architecture_diagram"` // 架構圖 (ASCII/Mermaid)
}

// ChoiceRequirements 選項需求條件
type ChoiceRequirements struct {
	MinCapital       int64    `json:"min_capital"`
	MinEmployees     int      `json:"min_employees"`
	MinSecurityLevel int      `json:"min_security_level"`
	RequiredInfra    []string `json:"required_infra"`
}

// ChoiceOutcomes 選項結果
type ChoiceOutcomes struct {
	CapitalChange       int64   `json:"capital_change"`
	EmployeeChange      int     `json:"employee_change"`
	SecurityChange      int     `json:"security_change"`
	CloudAdoptionChange float64 `json:"cloud_adoption_change"`
	SuccessRate         float64 `json:"success_rate"`      // 成功機率 (基於公司屬性調整)
	TimeToImplement     int     `json:"time_to_implement"` // 實施時間 (回合數)
}

// EventOutcome 事件結果
type EventOutcome struct {
	Success         bool     `json:"success"`
	Message         string   `json:"message"`
	CapitalDelta    int64    `json:"capital_delta"`
	EmployeeDelta   int      `json:"employee_delta"`
	LearningPoints  []string `json:"learning_points"`   // 學習要點
	AWSBestPractice string   `json:"aws_best_practice"` // AWS 最佳實踐說明
}
