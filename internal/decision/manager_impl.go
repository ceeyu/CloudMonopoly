package decision

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aws-learning-game/internal/company"
	"github.com/aws-learning-game/internal/event"
)

// random source for decision engine
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// decisionEngineImpl 決策引擎實作
type decisionEngineImpl struct{}

// NewDecisionEngine 建立新的決策引擎
func NewDecisionEngine() DecisionEngine {
	return &decisionEngineImpl{}
}

// EvaluateDecision 評估決策
func (d *decisionEngineImpl) EvaluateDecision(comp *company.Company, choice *event.EventChoice) (*DecisionEvaluation, error) {
	if comp == nil {
		return nil, fmt.Errorf("company cannot be nil")
	}
	if choice == nil {
		return nil, fmt.Errorf("choice cannot be nil")
	}

	evaluation := &DecisionEvaluation{
		IsEligible:         true,
		EligibilityIssues:  []string{},
		ImplementationTime: choice.Outcomes.TimeToImplement,
	}

	// 檢查資本需求
	if comp.Capital < choice.Requirements.MinCapital {
		evaluation.IsEligible = false
		evaluation.EligibilityIssues = append(evaluation.EligibilityIssues,
			fmt.Sprintf("資本不足: 需要 %d 萬元，目前只有 %d 萬元", choice.Requirements.MinCapital, comp.Capital))
	}

	// 檢查員工數需求
	if comp.Employees < choice.Requirements.MinEmployees {
		evaluation.IsEligible = false
		evaluation.EligibilityIssues = append(evaluation.EligibilityIssues,
			fmt.Sprintf("員工不足: 需要 %d 人，目前只有 %d 人", choice.Requirements.MinEmployees, comp.Employees))
	}

	// 檢查資安等級需求
	if comp.SecurityLevel < choice.Requirements.MinSecurityLevel {
		evaluation.IsEligible = false
		evaluation.EligibilityIssues = append(evaluation.EligibilityIssues,
			fmt.Sprintf("資安等級不足: 需要等級 %d，目前只有等級 %d", choice.Requirements.MinSecurityLevel, comp.SecurityLevel))
	}

	// 檢查基礎設施需求
	for _, required := range choice.Requirements.RequiredInfra {
		found := false
		for _, infra := range comp.Infrastructure {
			if infra == required {
				found = true
				break
			}
		}
		if !found {
			evaluation.IsEligible = false
			evaluation.EligibilityIssues = append(evaluation.EligibilityIssues,
				fmt.Sprintf("缺少必要基礎設施: %s", required))
		}
	}

	// 計算風險等級
	evaluation.RiskLevel = d.calculateRiskLevel(comp, choice)

	// 計算預期 ROI
	evaluation.ExpectedROI = d.calculateExpectedROI(comp, choice)

	// 產生建議
	evaluation.Recommendation = d.generateRecommendation(evaluation, choice)

	return evaluation, nil
}

// ExecuteDecision 執行決策
func (d *decisionEngineImpl) ExecuteDecision(comp *company.Company, choice *event.EventChoice) (*DecisionResult, error) {
	if comp == nil {
		return nil, fmt.Errorf("company cannot be nil")
	}
	if choice == nil {
		return nil, fmt.Errorf("choice cannot be nil")
	}

	// 先評估決策
	evaluation, err := d.EvaluateDecision(comp, choice)
	if err != nil {
		return nil, err
	}

	result := &DecisionResult{
		ActualOutcome:   choice.Outcomes,
		Penalties:       []Penalty{},
		Rewards:         []Reward{},
		LearningPoints:  []string{},
		AWSBestPractice: "",
		SAAExamTopics:   d.generateSAAExamTopics(choice),
	}

	// 如果不符合條件，套用懲罰
	if !evaluation.IsEligible {
		result.Success = false
		result.Message = "決策執行失敗"
		result.Explanation = "決策執行失敗：公司屬性不符合要求"
		result.Penalties = d.calculatePenalties(comp, choice, evaluation)
		// 調整實際結果（減少收益或增加損失）
		result.ActualOutcome = d.applyPenaltyToOutcome(choice.Outcomes, result.Penalties)
		// 加入失敗時的學習要點
		result.LearningPoints = d.generateFailureLearningPoints(comp, choice, evaluation)
		result.AWSBestPractice = d.generateAWSBestPractice(choice, false)
		return result, nil
	}

	// 計算成功機率
	successRate := d.calculateAdjustedSuccessRate(comp, choice)

	// 決定是否成功
	result.Success = rng.Float64() < successRate

	if result.Success {
		result.Message = "決策執行成功"
		result.Explanation = d.generateSuccessExplanation(choice)
		result.Rewards = d.calculateRewards(comp, choice)
		// 可能增加額外獎勵
		result.ActualOutcome = d.applyRewardToOutcome(choice.Outcomes, result.Rewards)
		// 加入成功時的學習要點
		result.LearningPoints = d.generateSuccessLearningPoints(comp, choice)
		result.AWSBestPractice = d.generateAWSBestPractice(choice, true)
	} else {
		result.Message = "決策未達預期效果"
		result.Explanation = d.generateFailureExplanation(choice)
		// 失敗時減少收益
		result.ActualOutcome = d.reduceOutcome(choice.Outcomes)
		// 加入失敗時的學習要點
		result.LearningPoints = d.generatePartialFailureLearningPoints(comp, choice)
		result.AWSBestPractice = d.generateAWSBestPractice(choice, false)
	}

	return result, nil
}

// GetComparison 取得方案比較
func (d *decisionEngineImpl) GetComparison(choices []event.EventChoice, comp *company.Company) (*Comparison, error) {
	if len(choices) == 0 {
		return nil, fmt.Errorf("choices cannot be empty")
	}
	if comp == nil {
		return nil, fmt.Errorf("company cannot be nil")
	}

	comparison := &Comparison{
		Choices:        make([]ChoiceComparison, len(choices)),
		ReasoningSteps: []string{},
	}

	bestScore := -1.0
	bestIndex := 0

	for i, choice := range choices {
		choiceComp := d.analyzeChoice(&choice, comp)
		comparison.Choices[i] = choiceComp

		// 計算綜合分數
		score := d.calculateOverallScore(choiceComp, comp)
		if score > bestScore {
			bestScore = score
			bestIndex = i
		}
	}

	comparison.Recommendation = bestIndex
	comparison.ReasoningSteps = d.generateReasoningSteps(choices, comparison.Choices, bestIndex, comp)

	return comparison, nil
}

// calculateRiskLevel 計算風險等級
func (d *decisionEngineImpl) calculateRiskLevel(comp *company.Company, choice *event.EventChoice) string {
	riskScore := 0

	// 資本風險
	if choice.Requirements.MinCapital > comp.Capital/2 {
		riskScore += 2
	} else if choice.Requirements.MinCapital > comp.Capital/4 {
		riskScore += 1
	}

	// 實施時間風險
	if choice.Outcomes.TimeToImplement > 3 {
		riskScore += 2
	} else if choice.Outcomes.TimeToImplement > 1 {
		riskScore += 1
	}

	// 成功率風險
	if choice.Outcomes.SuccessRate < 0.5 {
		riskScore += 2
	} else if choice.Outcomes.SuccessRate < 0.7 {
		riskScore += 1
	}

	// 技術債影響
	if comp.TechDebt > 50 {
		riskScore += 1
	}

	if riskScore >= 4 {
		return "high"
	} else if riskScore >= 2 {
		return "medium"
	}
	return "low"
}

// calculateExpectedROI 計算預期投資報酬率
func (d *decisionEngineImpl) calculateExpectedROI(comp *company.Company, choice *event.EventChoice) float64 {
	// 預期收益 = 資本變化 * 成功率
	expectedGain := float64(choice.Outcomes.CapitalChange) * choice.Outcomes.SuccessRate

	// 投入成本 = 最低資本需求的一部分
	cost := float64(choice.Requirements.MinCapital)
	if cost == 0 {
		cost = 100 // 預設最低成本
	}

	// ROI = (收益 - 成本) / 成本
	roi := (expectedGain - cost) / cost * 100
	return roi
}

// generateRecommendation 產生建議
func (d *decisionEngineImpl) generateRecommendation(eval *DecisionEvaluation, choice *event.EventChoice) string {
	if !eval.IsEligible {
		return "不建議執行此決策，公司目前不符合執行條件"
	}

	if eval.RiskLevel == "high" {
		return "此決策風險較高，建議謹慎評估或考慮其他方案"
	}

	if eval.ExpectedROI > 50 {
		return "此決策預期報酬率高，建議執行"
	} else if eval.ExpectedROI > 0 {
		return "此決策預期有正向報酬，可考慮執行"
	}

	return "此決策預期報酬率較低，建議評估其他選項"
}

// calculatePenalties 計算懲罰
func (d *decisionEngineImpl) calculatePenalties(comp *company.Company, choice *event.EventChoice, eval *DecisionEvaluation) []Penalty {
	penalties := []Penalty{}

	// 資本不足懲罰
	if comp.Capital < choice.Requirements.MinCapital {
		deficit := choice.Requirements.MinCapital - comp.Capital
		penalties = append(penalties, Penalty{
			Type:        "budget_overrun",
			Description: "資本不足導致預算超支",
			Impact:      deficit / 2, // 損失差額的一半
		})
	}

	// 員工不足懲罰
	if comp.Employees < choice.Requirements.MinEmployees {
		penalties = append(penalties, Penalty{
			Type:        "delay",
			Description: "人力不足導致專案延遲",
			Impact:      int64(choice.Outcomes.TimeToImplement * 50), // 每回合延遲損失 50 萬
		})
	}

	// 資安等級不足懲罰
	if comp.SecurityLevel < choice.Requirements.MinSecurityLevel {
		penalties = append(penalties, Penalty{
			Type:        "security_breach",
			Description: "資安等級不足導致安全風險",
			Impact:      200, // 固定損失 200 萬
		})
	}

	// 基礎設施不足懲罰
	for _, required := range choice.Requirements.RequiredInfra {
		found := false
		for _, infra := range comp.Infrastructure {
			if infra == required {
				found = true
				break
			}
		}
		if !found {
			penalties = append(penalties, Penalty{
				Type:        "infrastructure_gap",
				Description: fmt.Sprintf("缺少必要基礎設施 %s 導致額外建置成本", required),
				Impact:      150, // 每項缺少的基礎設施損失 150 萬
			})
		}
	}

	return penalties
}

// applyPenaltyToOutcome 套用懲罰到結果
func (d *decisionEngineImpl) applyPenaltyToOutcome(outcomes event.ChoiceOutcomes, penalties []Penalty) event.ChoiceOutcomes {
	result := outcomes

	totalImpact := int64(0)
	for _, p := range penalties {
		totalImpact += p.Impact
	}

	// 減少資本收益或增加損失
	result.CapitalChange -= totalImpact
	// 減少成功率
	result.SuccessRate *= 0.5
	// 增加實施時間
	result.TimeToImplement += 1

	return result
}

// calculateAdjustedSuccessRate 計算調整後的成功率
func (d *decisionEngineImpl) calculateAdjustedSuccessRate(comp *company.Company, choice *event.EventChoice) float64 {
	baseRate := choice.Outcomes.SuccessRate

	// 根據公司屬性調整
	adjustment := 0.0

	// 資本充裕加成
	if comp.Capital > choice.Requirements.MinCapital*2 {
		adjustment += 0.1
	}

	// 員工充裕加成
	if comp.Employees > choice.Requirements.MinEmployees*2 {
		adjustment += 0.05
	}

	// 資安等級加成
	if comp.SecurityLevel > choice.Requirements.MinSecurityLevel {
		adjustment += 0.05
	}

	// 雲端採用率對 AWS 方案的加成
	if choice.IsAWS && comp.CloudAdoption > 50 {
		adjustment += 0.1
	}

	// 技術債減成
	if comp.TechDebt > 50 {
		adjustment -= 0.1
	}

	finalRate := baseRate + adjustment
	if finalRate > 1.0 {
		finalRate = 1.0
	}
	if finalRate < 0.1 {
		finalRate = 0.1
	}

	return finalRate
}

// generateSuccessExplanation 產生成功說明
func (d *decisionEngineImpl) generateSuccessExplanation(choice *event.EventChoice) string {
	if choice.IsAWS {
		return fmt.Sprintf("成功採用 AWS 方案: %s。雲端服務的彈性和可擴展性為公司帶來了預期的效益。", choice.Title)
	}
	return fmt.Sprintf("成功實施地端方案: %s。穩定的基礎設施為公司提供了可靠的服務。", choice.Title)
}

// generateFailureExplanation 產生失敗說明
func (d *decisionEngineImpl) generateFailureExplanation(choice *event.EventChoice) string {
	if choice.IsAWS {
		return fmt.Sprintf("AWS 方案 %s 實施過程中遇到困難，未能達到預期效果。建議檢視實施策略和團隊能力。", choice.Title)
	}
	return fmt.Sprintf("地端方案 %s 實施過程中遇到技術挑戰，效果不如預期。建議評估是否需要額外資源。", choice.Title)
}

// calculateRewards 計算獎勵
func (d *decisionEngineImpl) calculateRewards(comp *company.Company, choice *event.EventChoice) []Reward {
	rewards := []Reward{}

	// 效率獎勵 - 資本充裕時
	if comp.Capital > choice.Requirements.MinCapital*2 {
		rewards = append(rewards, Reward{
			Type:        "efficiency",
			Description: "資源充裕，專案執行效率提升",
			Bonus:       50,
		})
	}

	// 成長獎勵 - AWS 方案且雲端採用率高
	if choice.IsAWS && comp.CloudAdoption > 70 {
		rewards = append(rewards, Reward{
			Type:        "growth",
			Description: "雲端轉型加速，業務成長潛力提升",
			Bonus:       100,
		})
	}

	// 創新獎勵 - 新創公司採用新技術
	if comp.Type == company.Startup && choice.IsAWS {
		rewards = append(rewards, Reward{
			Type:        "innovation",
			Description: "新創公司採用雲端技術，創新能力提升",
			Bonus:       75,
		})
	}

	return rewards
}

// applyRewardToOutcome 套用獎勵到結果
func (d *decisionEngineImpl) applyRewardToOutcome(outcomes event.ChoiceOutcomes, rewards []Reward) event.ChoiceOutcomes {
	result := outcomes

	totalBonus := int64(0)
	for _, r := range rewards {
		totalBonus += r.Bonus
	}

	result.CapitalChange += totalBonus

	return result
}

// reduceOutcome 減少結果（失敗時）
func (d *decisionEngineImpl) reduceOutcome(outcomes event.ChoiceOutcomes) event.ChoiceOutcomes {
	result := outcomes

	// 失敗時只獲得一半的收益
	if result.CapitalChange > 0 {
		result.CapitalChange /= 2
	}
	if result.EmployeeChange > 0 {
		result.EmployeeChange /= 2
	}

	return result
}

// analyzeChoice 分析選項
func (d *decisionEngineImpl) analyzeChoice(choice *event.EventChoice, comp *company.Company) ChoiceComparison {
	comparison := ChoiceComparison{
		ChoiceID: choice.ID,
		CostAnalysis: CostAnalysis{
			InitialCost:   choice.Requirements.MinCapital,
			MonthlyCost:   d.estimateMonthlyCost(choice),
			CostBreakdown: make(map[string]int64),
		},
		AWSExamRelevance: choice.AWSServices,
	}

	// 計算三年總成本
	comparison.CostAnalysis.ThreeYearTCO = comparison.CostAnalysis.InitialCost +
		comparison.CostAnalysis.MonthlyCost*36

	// 成本細項
	if choice.IsAWS {
		comparison.CostAnalysis.CostBreakdown["compute"] = comparison.CostAnalysis.MonthlyCost * 40 / 100
		comparison.CostAnalysis.CostBreakdown["storage"] = comparison.CostAnalysis.MonthlyCost * 30 / 100
		comparison.CostAnalysis.CostBreakdown["network"] = comparison.CostAnalysis.MonthlyCost * 20 / 100
		comparison.CostAnalysis.CostBreakdown["other"] = comparison.CostAnalysis.MonthlyCost * 10 / 100
	} else {
		comparison.CostAnalysis.CostBreakdown["hardware"] = comparison.CostAnalysis.InitialCost * 60 / 100
		comparison.CostAnalysis.CostBreakdown["software"] = comparison.CostAnalysis.InitialCost * 25 / 100
		comparison.CostAnalysis.CostBreakdown["maintenance"] = comparison.CostAnalysis.InitialCost * 15 / 100
	}

	// 可擴展性分數 (1-10)
	comparison.ScalabilityScore = d.calculateScalabilityScore(choice)

	// 複雜度分數 (1-10, 越低越好)
	comparison.ComplexityScore = d.calculateComplexityScore(choice)

	// 安全性分數 (1-10)
	comparison.SecurityScore = d.calculateSecurityScore(choice)

	return comparison
}

// estimateMonthlyCost 估算月成本
func (d *decisionEngineImpl) estimateMonthlyCost(choice *event.EventChoice) int64 {
	if choice.IsAWS {
		// AWS 方案通常有較低的初始成本但持續的月費
		return choice.Requirements.MinCapital / 10
	}
	// 地端方案通常有較高的初始成本但較低的月費
	return choice.Requirements.MinCapital / 50
}

// calculateScalabilityScore 計算可擴展性分數
func (d *decisionEngineImpl) calculateScalabilityScore(choice *event.EventChoice) int {
	if choice.IsAWS {
		// AWS 方案通常有較高的可擴展性
		return 8 + len(choice.AWSServices)%3
	}
	// 地端方案可擴展性較低
	return 4
}

// calculateComplexityScore 計算複雜度分數
func (d *decisionEngineImpl) calculateComplexityScore(choice *event.EventChoice) int {
	baseScore := 5

	// 需要更多基礎設施增加複雜度
	baseScore += len(choice.Requirements.RequiredInfra)

	// AWS 服務數量增加複雜度
	baseScore += len(choice.AWSServices) / 2

	// 實施時間長增加複雜度
	baseScore += choice.Outcomes.TimeToImplement / 2

	if baseScore > 10 {
		baseScore = 10
	}

	return baseScore
}

// calculateSecurityScore 計算安全性分數
func (d *decisionEngineImpl) calculateSecurityScore(choice *event.EventChoice) int {
	baseScore := 5

	// 資安等級需求反映安全性
	baseScore += choice.Requirements.MinSecurityLevel

	// AWS 方案通常有較好的安全性
	if choice.IsAWS {
		baseScore += 2
	}

	if baseScore > 10 {
		baseScore = 10
	}

	return baseScore
}

// calculateOverallScore 計算綜合分數
func (d *decisionEngineImpl) calculateOverallScore(comparison ChoiceComparison, comp *company.Company) float64 {
	// 權重: 成本 30%, 可擴展性 25%, 複雜度 20%, 安全性 25%
	costScore := 10.0 - float64(comparison.CostAnalysis.ThreeYearTCO)/float64(comp.Capital+1)*10
	if costScore < 0 {
		costScore = 0
	}

	scalabilityScore := float64(comparison.ScalabilityScore)
	complexityScore := 10.0 - float64(comparison.ComplexityScore) // 複雜度越低越好
	securityScore := float64(comparison.SecurityScore)

	return costScore*0.3 + scalabilityScore*0.25 + complexityScore*0.2 + securityScore*0.25
}

// generateReasoningSteps 產生推理步驟
func (d *decisionEngineImpl) generateReasoningSteps(choices []event.EventChoice, comparisons []ChoiceComparison, bestIndex int, comp *company.Company) []string {
	steps := []string{}

	steps = append(steps, fmt.Sprintf("分析了 %d 個方案選項", len(choices)))

	for i, c := range comparisons {
		steps = append(steps, fmt.Sprintf("方案 %d: 三年總成本 %d 萬元, 可擴展性 %d/10, 複雜度 %d/10, 安全性 %d/10",
			c.ChoiceID, c.CostAnalysis.ThreeYearTCO, c.ScalabilityScore, c.ComplexityScore, c.SecurityScore))

		if i == bestIndex {
			steps = append(steps, fmt.Sprintf("→ 方案 %d 綜合評分最高，推薦採用", c.ChoiceID))
		}
	}

	// 根據公司類型給出額外建議
	switch comp.Type {
	case company.Startup:
		steps = append(steps, "考量新創公司特性，建議優先考慮成本效益和快速部署能力")
	case company.Traditional:
		steps = append(steps, "考量傳產公司特性，建議優先考慮穩定性和漸進式轉型")
	case company.CloudReseller:
		steps = append(steps, "考量雲端代理商特性，建議優先考慮技術先進性和客戶價值")
	case company.CloudNative:
		steps = append(steps, "考量雲端公司特性，建議優先考慮創新性和擴展能力")
	}

	return steps
}

// generateSuccessLearningPoints 產生成功時的學習要點
func (d *decisionEngineImpl) generateSuccessLearningPoints(comp *company.Company, choice *event.EventChoice) []string {
	points := []string{}

	if choice.IsAWS {
		points = append(points, "雲端服務的彈性讓您能夠根據需求快速擴展或縮減資源")

		// 根據使用的 AWS 服務加入相關學習要點
		for _, service := range choice.AWSServices {
			switch service {
			case "EC2":
				points = append(points, "EC2 提供可調整的運算容量，是 AWS 最基礎的運算服務")
			case "S3":
				points = append(points, "S3 提供高耐久性的物件儲存，適合存放靜態資源和備份")
			case "RDS":
				points = append(points, "RDS 簡化了資料庫管理，自動處理備份、修補和擴展")
			case "Lambda":
				points = append(points, "Lambda 無伺服器架構讓您只需為實際執行時間付費")
			case "CloudFront":
				points = append(points, "CloudFront CDN 可加速全球內容傳遞，降低延遲")
			case "DynamoDB":
				points = append(points, "DynamoDB 提供毫秒級延遲的 NoSQL 資料庫服務")
			case "ELB":
				points = append(points, "Elastic Load Balancing 自動分配流量，提高可用性")
			case "Auto Scaling":
				points = append(points, "Auto Scaling 根據需求自動調整 EC2 容量")
			case "VPC":
				points = append(points, "VPC 讓您在 AWS 中建立隔離的虛擬網路環境")
			case "IAM":
				points = append(points, "IAM 提供細粒度的存取控制，是 AWS 安全的基礎")
			case "CloudWatch":
				points = append(points, "CloudWatch 提供監控和可觀測性，幫助您了解系統狀態")
			case "SNS":
				points = append(points, "SNS 提供發布/訂閱訊息服務，適合事件驅動架構")
			case "SQS":
				points = append(points, "SQS 提供可靠的訊息佇列，解耦系統元件")
			}
		}

		// 根據公司類型加入建議
		switch comp.Type {
		case company.Startup:
			points = append(points, "新創公司採用雲端服務可以降低初期資本支出，專注於核心業務")
		case company.Traditional:
			points = append(points, "傳統企業透過雲端轉型可以提升競爭力，但需要注意變更管理")
		case company.CloudReseller:
			points = append(points, "作為雲端代理商，深入了解 AWS 服務有助於為客戶提供更好的建議")
		case company.CloudNative:
			points = append(points, "雲端原生公司應持續優化架構，善用最新的 AWS 服務")
		}
	} else {
		points = append(points, "地端方案提供完全的控制權，適合有特殊合規需求的場景")
		points = append(points, "地端基礎設施需要考慮長期維護成本和人力需求")
		points = append(points, "混合雲架構可以結合地端和雲端的優勢")
	}

	return points
}

// generateFailureLearningPoints 產生不符合條件時的學習要點
func (d *decisionEngineImpl) generateFailureLearningPoints(comp *company.Company, choice *event.EventChoice, eval *DecisionEvaluation) []string {
	points := []string{}

	points = append(points, "在做出重大技術決策前，務必評估公司的資源和能力是否足夠")

	for _, issue := range eval.EligibilityIssues {
		if contains(issue, "資本不足") {
			points = append(points, "AWS 提供多種成本優化方案，如 Reserved Instances 和 Savings Plans")
			points = append(points, "考慮使用 AWS Cost Explorer 分析和預測成本")
		}
		if contains(issue, "員工不足") {
			points = append(points, "AWS 託管服務可以減少維運人力需求")
			points = append(points, "考慮使用 AWS Managed Services 或尋求 AWS Partner 協助")
		}
		if contains(issue, "資安等級不足") {
			points = append(points, "AWS 提供多層次的安全服務，如 WAF、Shield、GuardDuty")
			points = append(points, "建議先提升基礎安全能力，再進行大型專案")
		}
		if contains(issue, "基礎設施") {
			points = append(points, "AWS 提供遷移服務如 AWS Migration Hub 協助基礎設施轉型")
		}
	}

	return points
}

// generatePartialFailureLearningPoints 產生執行失敗時的學習要點
func (d *decisionEngineImpl) generatePartialFailureLearningPoints(comp *company.Company, choice *event.EventChoice) []string {
	points := []string{}

	points = append(points, "技術專案的成功不僅取決於選擇正確的方案，還需要良好的執行")

	if choice.IsAWS {
		points = append(points, "AWS Well-Architected Framework 提供最佳實踐指南")
		points = append(points, "考慮使用 AWS Professional Services 或認證合作夥伴協助實施")
		points = append(points, "建議進行 Proof of Concept (PoC) 驗證方案可行性")
	} else {
		points = append(points, "地端專案需要充足的規劃和測試時間")
		points = append(points, "考慮分階段實施，降低風險")
	}

	// 根據公司技術債提供建議
	if comp.TechDebt > 50 {
		points = append(points, "高技術債會影響新專案的成功率，建議優先處理技術債")
	}

	return points
}

// generateAWSBestPractice 產生 AWS 最佳實踐說明
func (d *decisionEngineImpl) generateAWSBestPractice(choice *event.EventChoice, success bool) string {
	if !choice.IsAWS {
		if success {
			return "地端方案成功實施。建議定期評估是否有適合遷移到雲端的工作負載，以獲得更好的彈性和成本效益。"
		}
		return "地端方案實施遇到挑戰。AWS 提供多種遷移工具和服務，可以協助您逐步將工作負載遷移到雲端。"
	}

	// 根據使用的 AWS 服務提供最佳實踐
	practices := []string{}

	for _, service := range choice.AWSServices {
		switch service {
		case "EC2":
			practices = append(practices, "使用 EC2 時，建議根據工作負載選擇適當的執行個體類型，並考慮使用 Spot Instances 降低成本")
		case "S3":
			practices = append(practices, "S3 最佳實踐包括啟用版本控制、設定生命週期政策、使用適當的儲存類別")
		case "RDS":
			practices = append(practices, "RDS 建議啟用 Multi-AZ 部署以提高可用性，並設定自動備份")
		case "Lambda":
			practices = append(practices, "Lambda 最佳實踐包括最小化部署套件大小、適當設定記憶體和逾時")
		case "DynamoDB":
			practices = append(practices, "DynamoDB 建議使用 On-Demand 容量模式或仔細規劃 Provisioned 容量")
		case "VPC":
			practices = append(practices, "VPC 設計應考慮子網路規劃、安全群組和網路 ACL 的最佳實踐")
		case "IAM":
			practices = append(practices, "IAM 最佳實踐包括最小權限原則、使用角色而非長期憑證、啟用 MFA")
		}
	}

	if len(practices) > 0 {
		return practices[0] // 返回第一個最相關的最佳實踐
	}

	if success {
		return "AWS 架構設計應遵循 Well-Architected Framework 的五大支柱：卓越營運、安全性、可靠性、效能效率、成本優化。"
	}
	return "建議參考 AWS Well-Architected Framework 重新評估架構設計，確保符合最佳實踐。"
}

// contains 檢查字串是否包含子字串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// generateSAAExamTopics 產生 AWS SAA 考試相關主題
// 根據選項中使用的 AWS 服務，對應到 SAA 考試的相關主題
func (d *decisionEngineImpl) generateSAAExamTopics(choice *event.EventChoice) []string {
	topics := []string{}
	topicsSet := make(map[string]bool)

	// AWS 服務到 SAA 考試主題的對應
	serviceToTopics := map[string][]string{
		// 運算服務
		"EC2":          {"EC2 執行個體類型", "EC2 定價模式", "EC2 Auto Scaling"},
		"Lambda":       {"無伺服器架構", "Lambda 最佳實踐", "事件驅動架構"},
		"ECS":          {"容器服務", "ECS vs EKS", "微服務架構"},
		"EKS":          {"Kubernetes on AWS", "容器編排", "微服務架構"},
		"Auto Scaling": {"彈性擴展", "擴展策略", "高可用性架構"},

		// 儲存服務
		"S3":  {"S3 儲存類別", "S3 生命週期政策", "S3 安全性"},
		"EBS": {"EBS 磁碟區類型", "EBS 快照", "區塊儲存"},
		"EFS": {"彈性檔案系統", "共享儲存", "NFS 協定"},

		// 資料庫服務
		"RDS":         {"RDS 多可用區", "RDS 讀取副本", "關聯式資料庫"},
		"Aurora":      {"Aurora 架構", "Aurora Serverless", "高效能資料庫"},
		"DynamoDB":    {"NoSQL 資料庫", "DynamoDB 容量模式", "全域表"},
		"ElastiCache": {"快取策略", "Redis vs Memcached", "效能優化"},
		"Redshift":    {"資料倉儲", "資料分析", "OLAP"},

		// 網路服務
		"VPC":         {"VPC 設計", "子網路規劃", "網路 ACL vs 安全群組"},
		"CloudFront":  {"CDN", "邊緣運算", "內容分發"},
		"Route 53":    {"DNS 服務", "路由策略", "健康檢查"},
		"ELB":         {"負載平衡", "ALB vs NLB", "高可用性"},
		"API Gateway": {"API 管理", "REST API", "WebSocket API"},

		// 安全服務
		"IAM":             {"身分與存取管理", "IAM 政策", "最小權限原則"},
		"KMS":             {"金鑰管理", "加密", "資料保護"},
		"WAF":             {"Web 應用程式防火牆", "DDoS 防護", "安全規則"},
		"Shield":          {"DDoS 防護", "AWS Shield Advanced", "安全性"},
		"GuardDuty":       {"威脅偵測", "安全監控", "異常偵測"},
		"Cognito":         {"使用者認證", "身分聯合", "OAuth"},
		"Secrets Manager": {"密鑰管理", "憑證輪換", "安全性"},
		"Security Hub":    {"安全態勢管理", "合規性", "安全最佳實踐"},

		// 監控與管理
		"CloudWatch":      {"監控與日誌", "指標與警報", "可觀測性"},
		"CloudTrail":      {"稽核日誌", "API 追蹤", "合規性"},
		"AWS Config":      {"組態管理", "合規性規則", "資源追蹤"},
		"X-Ray":           {"分散式追蹤", "效能分析", "除錯"},
		"Trusted Advisor": {"最佳實踐建議", "成本優化", "安全性檢查"},

		// 應用整合
		"SNS":            {"發布/訂閱", "訊息通知", "解耦架構"},
		"SQS":            {"訊息佇列", "解耦架構", "非同步處理"},
		"Step Functions": {"工作流程編排", "狀態機", "無伺服器編排"},
		"EventBridge":    {"事件匯流排", "事件驅動架構", "整合"},

		// 資料分析
		"Kinesis":    {"串流資料處理", "即時分析", "資料管道"},
		"Glue":       {"ETL 服務", "資料目錄", "資料整合"},
		"Athena":     {"無伺服器查詢", "S3 資料分析", "SQL 查詢"},
		"QuickSight": {"商業智慧", "資料視覺化", "儀表板"},

		// 機器學習
		"SageMaker": {"機器學習平台", "模型訓練", "ML 部署"},

		// IoT
		"IoT Core":      {"物聯網", "裝置管理", "MQTT"},
		"IoT Analytics": {"IoT 資料分析", "時間序列", "裝置資料"},

		// 遷移與傳輸
		"DMS":          {"資料庫遷移", "異質遷移", "持續複寫"},
		"AWS Backup":   {"備份服務", "災難復原", "資料保護"},
		"AWS Artifact": {"合規報告", "稽核文件", "認證"},

		// 成本管理
		"Cost Explorer": {"成本分析", "成本優化", "預算管理"},
		"Savings Plans": {"成本節省", "承諾使用折扣", "定價優化"},

		// 開發工具
		"CodePipeline": {"CI/CD", "持續整合", "持續部署"},
		"CodeBuild":    {"建置服務", "自動化建置", "DevOps"},
		"CodeDeploy":   {"部署自動化", "藍綠部署", "滾動更新"},

		// 其他
		"Amplify":    {"行動開發", "前端部署", "全端開發"},
		"AppSync":    {"GraphQL", "即時資料", "離線同步"},
		"OpenSearch": {"搜尋服務", "日誌分析", "全文搜尋"},
		"ECR":        {"容器映像庫", "Docker Registry", "映像管理"},
	}

	// 根據選項中的 AWS 服務加入對應主題
	for _, service := range choice.AWSServices {
		if examTopics, ok := serviceToTopics[service]; ok {
			for _, topic := range examTopics {
				if !topicsSet[topic] {
					topicsSet[topic] = true
					topics = append(topics, topic)
				}
			}
		}
		// 也加入服務本身作為主題
		if !topicsSet[service] {
			topicsSet[service] = true
			topics = append(topics, service)
		}
	}

	// 如果是 AWS 方案，加入通用主題
	if choice.IsAWS {
		generalTopics := []string{"AWS Well-Architected Framework", "雲端架構設計"}
		for _, topic := range generalTopics {
			if !topicsSet[topic] {
				topicsSet[topic] = true
				topics = append(topics, topic)
			}
		}
	}

	// 如果沒有任何主題，返回基本主題
	if len(topics) == 0 {
		if choice.IsAWS {
			topics = []string{"AWS 雲端服務", "雲端架構設計"}
		} else {
			topics = []string{"混合雲架構", "地端與雲端整合"}
		}
	}

	return topics
}
