package decision

import (
	"github.com/aws-learning-game/internal/company"
	"github.com/aws-learning-game/internal/event"
)

// DecisionEngine 決策引擎介面
type DecisionEngine interface {
	// EvaluateDecision 評估決策
	EvaluateDecision(comp *company.Company, choice *event.EventChoice) (*DecisionEvaluation, error)
	// ExecuteDecision 執行決策
	ExecuteDecision(comp *company.Company, choice *event.EventChoice) (*DecisionResult, error)
	// GetComparison 取得方案比較
	GetComparison(choices []event.EventChoice, comp *company.Company) (*Comparison, error)
}
