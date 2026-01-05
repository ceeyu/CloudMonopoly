package decision

import "github.com/aws-learning-game/internal/event"

// DecisionEvaluation 決策評估結果
type DecisionEvaluation struct {
	IsEligible         bool     `json:"is_eligible"`
	EligibilityIssues  []string `json:"eligibility_issues"`
	RiskLevel          string   `json:"risk_level"` // "low", "medium", "high"
	ExpectedROI        float64  `json:"expected_roi"`
	ImplementationTime int      `json:"implementation_time"`
	Recommendation     string   `json:"recommendation"`
}

// DecisionResult 決策執行結果
type DecisionResult struct {
	Success         bool                 `json:"success"`
	Message         string               `json:"message"`
	ActualOutcome   event.ChoiceOutcomes `json:"actual_outcome"`
	Explanation     string               `json:"explanation"`
	Penalties       []Penalty            `json:"penalties"`
	Rewards         []Reward             `json:"rewards"`
	LearningPoints  []string             `json:"learning_points"`   // 學習要點
	AWSBestPractice string               `json:"aws_best_practice"` // AWS 最佳實踐說明
	SAAExamTopics   []string             `json:"saa_exam_topics"`   // 相關 AWS SAA 考點
}

// Penalty 懲罰
type Penalty struct {
	Type        string `json:"type"` // "budget_overrun", "delay", "security_breach"
	Description string `json:"description"`
	Impact      int64  `json:"impact"`
}

// Reward 獎勵
type Reward struct {
	Type        string `json:"type"` // "efficiency", "growth", "innovation"
	Description string `json:"description"`
	Bonus       int64  `json:"bonus"`
}

// Comparison 方案比較
type Comparison struct {
	Choices        []ChoiceComparison `json:"choices"`
	Recommendation int                `json:"recommendation"` // 推薦選項 index
	ReasoningSteps []string           `json:"reasoning_steps"`
}

// ChoiceComparison 選項比較
type ChoiceComparison struct {
	ChoiceID         int          `json:"choice_id"`
	CostAnalysis     CostAnalysis `json:"cost_analysis"`
	ScalabilityScore int          `json:"scalability_score"`
	ComplexityScore  int          `json:"complexity_score"`
	SecurityScore    int          `json:"security_score"`
	AWSExamRelevance []string     `json:"aws_exam_relevance"`
}

// CostAnalysis 成本分析
type CostAnalysis struct {
	InitialCost   int64            `json:"initial_cost"`
	MonthlyCost   int64            `json:"monthly_cost"`
	ThreeYearTCO  int64            `json:"three_year_tco"`
	CostBreakdown map[string]int64 `json:"cost_breakdown"`
}
