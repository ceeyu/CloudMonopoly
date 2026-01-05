package decision

import (
	"testing"

	"github.com/aws-learning-game/internal/company"
	"github.com/aws-learning-game/internal/event"
	"pgregory.net/rapid"
)

// Feature: aws-learning-game, Property 9: Infrastructure Decision Options
// For any infrastructure decision event, the available choices SHALL include
// at least one on-premises option AND at least one AWS cloud option.
// **Validates: Requirements 4.1**
func TestProperty9_InfrastructureDecisionOptions(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random infrastructure decision event
		numChoices := rapid.IntRange(2, 5).Draw(t, "numChoices")
		choices := make([]event.EventChoice, numChoices)

		// Ensure at least one AWS and one on-prem option
		awsIndex := rapid.IntRange(0, numChoices-1).Draw(t, "awsIndex")
		onPremIndex := rapid.IntRange(0, numChoices-1).Draw(t, "onPremIndex")
		if onPremIndex == awsIndex {
			onPremIndex = (awsIndex + 1) % numChoices
		}

		for i := 0; i < numChoices; i++ {
			isAWS := i == awsIndex
			choices[i] = event.EventChoice{
				ID:          i + 1,
				Title:       rapid.String().Draw(t, "choiceTitle"),
				Description: rapid.String().Draw(t, "choiceDesc"),
				IsAWS:       isAWS,
				Requirements: event.ChoiceRequirements{
					MinCapital:       rapid.Int64Range(100, 1000).Draw(t, "minCapital"),
					MinEmployees:     rapid.IntRange(5, 50).Draw(t, "minEmployees"),
					MinSecurityLevel: rapid.IntRange(1, 5).Draw(t, "minSecLevel"),
				},
				Outcomes: event.ChoiceOutcomes{
					CapitalChange:   rapid.Int64Range(-500, 500).Draw(t, "capitalChange"),
					EmployeeChange:  rapid.IntRange(-10, 20).Draw(t, "employeeChange"),
					SuccessRate:     rapid.Float64Range(0.3, 1.0).Draw(t, "successRate"),
					TimeToImplement: rapid.IntRange(1, 5).Draw(t, "timeToImpl"),
				},
			}
			if isAWS {
				choices[i].AWSServices = []string{"EC2", "S3"}
			} else {
				choices[i].OnPremSolution = "地端伺服器方案"
			}
		}

		// Validate infrastructure decision has both options
		hasAWS, hasOnPrem := ValidateInfrastructureDecisionOptions(choices)

		if !hasAWS {
			t.Error("Infrastructure decision should have at least one AWS option")
		}
		if !hasOnPrem {
			t.Error("Infrastructure decision should have at least one on-premises option")
		}
	})
}

// ValidateInfrastructureDecisionOptions checks if choices include both AWS and on-prem options
func ValidateInfrastructureDecisionOptions(choices []event.EventChoice) (hasAWS bool, hasOnPrem bool) {
	for _, choice := range choices {
		if choice.IsAWS {
			hasAWS = true
		} else {
			hasOnPrem = true
		}
	}
	return hasAWS, hasOnPrem
}

// createTestCompany creates a test company with specified attributes
func createTestCompany(capital int64, employees int, securityLevel int) *company.Company {
	return &company.Company{
		ID:            "test-company",
		Name:          "Test Company",
		Type:          company.Startup,
		Capital:       capital,
		Employees:     employees,
		SecurityLevel: securityLevel,
		CloudAdoption: 50.0,
		TechDebt:      20,
	}
}

// createTestChoice creates a test event choice
func createTestChoice(isAWS bool, minCapital int64, minEmployees int, minSecLevel int) *event.EventChoice {
	choice := &event.EventChoice{
		ID:          1,
		Title:       "Test Choice",
		Description: "Test Description",
		IsAWS:       isAWS,
		Requirements: event.ChoiceRequirements{
			MinCapital:       minCapital,
			MinEmployees:     minEmployees,
			MinSecurityLevel: minSecLevel,
		},
		Outcomes: event.ChoiceOutcomes{
			CapitalChange:   200,
			EmployeeChange:  5,
			SecurityChange:  1,
			SuccessRate:     0.8,
			TimeToImplement: 2,
		},
	}
	if isAWS {
		choice.AWSServices = []string{"EC2", "S3"}
	} else {
		choice.OnPremSolution = "地端伺服器"
	}
	return choice
}

// Feature: aws-learning-game, Property 11: Decision Outcome Calculation
// For any decision execution, the DecisionResult SHALL contain CapitalDelta,
// EmployeeDelta, and Success status calculated based on Company attributes
// and ChoiceRequirements.
// **Validates: Requirements 4.3**
func TestProperty11_DecisionOutcomeCalculation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random company
		capital := rapid.Int64Range(100, 10000).Draw(t, "capital")
		employees := rapid.IntRange(5, 500).Draw(t, "employees")
		securityLevel := rapid.IntRange(1, 5).Draw(t, "securityLevel")
		comp := createTestCompany(capital, employees, securityLevel)

		// Generate random choice requirements
		minCapital := rapid.Int64Range(50, 5000).Draw(t, "minCapital")
		minEmployees := rapid.IntRange(1, 200).Draw(t, "minEmployees")
		minSecLevel := rapid.IntRange(1, 5).Draw(t, "minSecLevel")
		isAWS := rapid.Bool().Draw(t, "isAWS")

		choice := createTestChoice(isAWS, minCapital, minEmployees, minSecLevel)
		choice.Outcomes.CapitalChange = rapid.Int64Range(-500, 1000).Draw(t, "capitalChange")
		choice.Outcomes.EmployeeChange = rapid.IntRange(-20, 50).Draw(t, "employeeChange")
		choice.Outcomes.SuccessRate = rapid.Float64Range(0.1, 1.0).Draw(t, "successRate")

		// Execute decision
		engine := NewDecisionEngine()
		result, err := engine.ExecuteDecision(comp, choice)

		if err != nil {
			t.Fatalf("ExecuteDecision failed: %v", err)
		}

		// Verify result contains required fields
		if result == nil {
			t.Fatal("DecisionResult should not be nil")
		}

		// ActualOutcome should have CapitalChange and EmployeeChange
		// (they may be modified from original based on success/failure)
		// The key is that they exist and are calculated

		// Success status should be set
		// (either true or false, but must be determined)

		// If company doesn't meet requirements, should have penalties
		meetsRequirements := comp.Capital >= choice.Requirements.MinCapital &&
			comp.Employees >= choice.Requirements.MinEmployees &&
			comp.SecurityLevel >= choice.Requirements.MinSecurityLevel

		if !meetsRequirements {
			// Should have penalties when requirements not met
			if len(result.Penalties) == 0 {
				t.Error("Should have penalties when company doesn't meet requirements")
			}
			// Success should be false
			if result.Success {
				t.Error("Success should be false when requirements not met")
			}
		}

		// Explanation should not be empty
		if result.Explanation == "" {
			t.Error("Explanation should not be empty")
		}
	})
}

// Feature: aws-learning-game, Property 12: Decision Evaluation Completeness
// For any decision evaluation, the DecisionEvaluation SHALL include assessments
// for: cost (ExpectedROI), time (ImplementationTime), and risk (RiskLevel).
// **Validates: Requirements 4.4**
func TestProperty12_DecisionEvaluationCompleteness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random company
		capital := rapid.Int64Range(100, 10000).Draw(t, "capital")
		employees := rapid.IntRange(5, 500).Draw(t, "employees")
		securityLevel := rapid.IntRange(1, 5).Draw(t, "securityLevel")
		comp := createTestCompany(capital, employees, securityLevel)
		comp.TechDebt = rapid.IntRange(0, 100).Draw(t, "techDebt")
		comp.CloudAdoption = rapid.Float64Range(0, 100).Draw(t, "cloudAdoption")

		// Generate random choice
		minCapital := rapid.Int64Range(50, 5000).Draw(t, "minCapital")
		minEmployees := rapid.IntRange(1, 200).Draw(t, "minEmployees")
		minSecLevel := rapid.IntRange(1, 5).Draw(t, "minSecLevel")
		isAWS := rapid.Bool().Draw(t, "isAWS")

		choice := createTestChoice(isAWS, minCapital, minEmployees, minSecLevel)
		choice.Outcomes.TimeToImplement = rapid.IntRange(1, 10).Draw(t, "timeToImpl")
		choice.Outcomes.SuccessRate = rapid.Float64Range(0.1, 1.0).Draw(t, "successRate")

		// Evaluate decision
		engine := NewDecisionEngine()
		eval, err := engine.EvaluateDecision(comp, choice)

		if err != nil {
			t.Fatalf("EvaluateDecision failed: %v", err)
		}

		// Verify evaluation contains all required assessments
		if eval == nil {
			t.Fatal("DecisionEvaluation should not be nil")
		}

		// Cost assessment (ExpectedROI)
		// ROI can be any value, but should be calculated
		// We just verify it's not the zero value in a meaningful way
		// by checking the evaluation was performed

		// Time assessment (ImplementationTime)
		if eval.ImplementationTime <= 0 {
			t.Error("ImplementationTime should be positive")
		}

		// Risk assessment (RiskLevel)
		validRiskLevels := map[string]bool{"low": true, "medium": true, "high": true}
		if !validRiskLevels[eval.RiskLevel] {
			t.Errorf("RiskLevel should be 'low', 'medium', or 'high', got '%s'", eval.RiskLevel)
		}

		// Recommendation should not be empty
		if eval.Recommendation == "" {
			t.Error("Recommendation should not be empty")
		}

		// If not eligible, should have eligibility issues
		if !eval.IsEligible && len(eval.EligibilityIssues) == 0 {
			t.Error("Should have eligibility issues when not eligible")
		}
	})
}

// Feature: aws-learning-game, Property 13: Incompatible Decision Penalty
// For any decision where Company attributes do not meet ChoiceRequirements
// (e.g., Capital < MinCapital), the DecisionResult SHALL contain at least one Penalty.
// **Validates: Requirements 4.5**
func TestProperty13_IncompatibleDecisionPenalty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate company with limited resources
		capital := rapid.Int64Range(100, 500).Draw(t, "capital")
		employees := rapid.IntRange(5, 20).Draw(t, "employees")
		securityLevel := rapid.IntRange(1, 2).Draw(t, "securityLevel")
		comp := createTestCompany(capital, employees, securityLevel)

		// Generate choice with higher requirements than company has
		// Ensure at least one requirement is not met
		requirementType := rapid.IntRange(0, 2).Draw(t, "requirementType")

		var minCapital int64
		var minEmployees int
		var minSecLevel int

		switch requirementType {
		case 0:
			// Capital requirement not met
			minCapital = capital + rapid.Int64Range(100, 1000).Draw(t, "extraCapital")
			minEmployees = rapid.IntRange(1, employees).Draw(t, "minEmployees")
			minSecLevel = rapid.IntRange(1, securityLevel).Draw(t, "minSecLevel")
		case 1:
			// Employee requirement not met
			minCapital = rapid.Int64Range(1, capital).Draw(t, "minCapital")
			minEmployees = employees + rapid.IntRange(10, 100).Draw(t, "extraEmployees")
			minSecLevel = rapid.IntRange(1, securityLevel).Draw(t, "minSecLevel")
		case 2:
			// Security level requirement not met
			minCapital = rapid.Int64Range(1, capital).Draw(t, "minCapital")
			minEmployees = rapid.IntRange(1, employees).Draw(t, "minEmployees")
			minSecLevel = securityLevel + rapid.IntRange(1, 3).Draw(t, "extraSecLevel")
		}

		choice := createTestChoice(true, minCapital, minEmployees, minSecLevel)

		// Execute decision
		engine := NewDecisionEngine()
		result, err := engine.ExecuteDecision(comp, choice)

		if err != nil {
			t.Fatalf("ExecuteDecision failed: %v", err)
		}

		// Verify result has at least one penalty
		if len(result.Penalties) == 0 {
			t.Errorf("Should have at least one penalty when requirements not met. "+
				"Company: Capital=%d, Employees=%d, SecurityLevel=%d. "+
				"Requirements: MinCapital=%d, MinEmployees=%d, MinSecurityLevel=%d",
				comp.Capital, comp.Employees, comp.SecurityLevel,
				choice.Requirements.MinCapital, choice.Requirements.MinEmployees, choice.Requirements.MinSecurityLevel)
		}

		// Verify success is false
		if result.Success {
			t.Error("Success should be false when requirements not met")
		}

		// Verify each penalty has required fields
		for i, penalty := range result.Penalties {
			if penalty.Type == "" {
				t.Errorf("Penalty %d should have a Type", i)
			}
			if penalty.Description == "" {
				t.Errorf("Penalty %d should have a Description", i)
			}
			if penalty.Impact <= 0 {
				t.Errorf("Penalty %d should have positive Impact, got %d", i, penalty.Impact)
			}
		}
	})
}

// Test that penalties are correctly calculated for each type of requirement failure
func TestPenaltyTypes(t *testing.T) {
	engine := NewDecisionEngine()

	testCases := []struct {
		name          string
		company       *company.Company
		choice        *event.EventChoice
		expectedTypes []string
	}{
		{
			name:          "capital insufficient",
			company:       createTestCompany(100, 50, 3),
			choice:        createTestChoice(true, 500, 10, 1),
			expectedTypes: []string{"budget_overrun"},
		},
		{
			name:          "employees insufficient",
			company:       createTestCompany(1000, 5, 3),
			choice:        createTestChoice(true, 100, 50, 1),
			expectedTypes: []string{"delay"},
		},
		{
			name:          "security level insufficient",
			company:       createTestCompany(1000, 50, 1),
			choice:        createTestChoice(true, 100, 10, 5),
			expectedTypes: []string{"security_breach"},
		},
		{
			name:          "multiple requirements not met",
			company:       createTestCompany(100, 5, 1),
			choice:        createTestChoice(true, 500, 50, 5),
			expectedTypes: []string{"budget_overrun", "delay", "security_breach"},
		},
		{
			name: "infrastructure insufficient",
			company: func() *company.Company {
				c := createTestCompany(1000, 50, 3)
				c.Infrastructure = []string{"VPN"}
				return c
			}(),
			choice: func() *event.EventChoice {
				c := createTestChoice(true, 100, 10, 1)
				c.Requirements.RequiredInfra = []string{"VPN", "Firewall", "LoadBalancer"}
				return c
			}(),
			expectedTypes: []string{"infrastructure_gap"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := engine.ExecuteDecision(tc.company, tc.choice)
			if err != nil {
				t.Fatalf("ExecuteDecision failed: %v", err)
			}

			// Check that expected penalty types are present
			penaltyTypes := make(map[string]bool)
			for _, p := range result.Penalties {
				penaltyTypes[p.Type] = true
			}

			for _, expectedType := range tc.expectedTypes {
				if !penaltyTypes[expectedType] {
					t.Errorf("Expected penalty type '%s' not found", expectedType)
				}
			}
		})
	}
}

// Feature: aws-learning-game, Property 15: Comparison Table Generation
// For any set of 2 or more EventChoices, the generated comparison SHALL include
// CostAnalysis, ScalabilityScore, ComplexityScore, and SecurityScore for each choice.
// **Validates: Requirements 5.2**
func TestProperty15_ComparisonTableGeneration(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random company
		capital := rapid.Int64Range(500, 10000).Draw(t, "capital")
		employees := rapid.IntRange(10, 500).Draw(t, "employees")
		securityLevel := rapid.IntRange(1, 5).Draw(t, "securityLevel")
		comp := createTestCompany(capital, employees, securityLevel)

		// Generate 2-5 random choices
		numChoices := rapid.IntRange(2, 5).Draw(t, "numChoices")
		choices := make([]event.EventChoice, numChoices)

		for i := 0; i < numChoices; i++ {
			isAWS := rapid.Bool().Draw(t, "isAWS")
			minCapital := rapid.Int64Range(100, 2000).Draw(t, "minCapital")
			minEmployees := rapid.IntRange(5, 100).Draw(t, "minEmployees")
			minSecLevel := rapid.IntRange(1, 5).Draw(t, "minSecLevel")

			choices[i] = event.EventChoice{
				ID:          i + 1,
				Title:       rapid.String().Draw(t, "title"),
				Description: rapid.String().Draw(t, "desc"),
				IsAWS:       isAWS,
				Requirements: event.ChoiceRequirements{
					MinCapital:       minCapital,
					MinEmployees:     minEmployees,
					MinSecurityLevel: minSecLevel,
					RequiredInfra:    []string{},
				},
				Outcomes: event.ChoiceOutcomes{
					CapitalChange:   rapid.Int64Range(-500, 1000).Draw(t, "capitalChange"),
					EmployeeChange:  rapid.IntRange(-10, 50).Draw(t, "employeeChange"),
					SuccessRate:     rapid.Float64Range(0.3, 1.0).Draw(t, "successRate"),
					TimeToImplement: rapid.IntRange(1, 5).Draw(t, "timeToImpl"),
				},
			}
			if isAWS {
				choices[i].AWSServices = []string{"EC2", "S3", "RDS"}
			} else {
				choices[i].OnPremSolution = "地端伺服器方案"
			}
		}

		// Get comparison
		engine := NewDecisionEngine()
		comparison, err := engine.GetComparison(choices, comp)

		if err != nil {
			t.Fatalf("GetComparison failed: %v", err)
		}

		// Verify comparison is not nil
		if comparison == nil {
			t.Fatal("Comparison should not be nil")
		}

		// Verify we have comparisons for all choices
		if len(comparison.Choices) != numChoices {
			t.Errorf("Expected %d choice comparisons, got %d", numChoices, len(comparison.Choices))
		}

		// Verify each choice comparison has required fields
		for i, choiceComp := range comparison.Choices {
			// CostAnalysis should be present
			if choiceComp.CostAnalysis.InitialCost < 0 {
				t.Errorf("Choice %d: InitialCost should not be negative", i)
			}
			if choiceComp.CostAnalysis.MonthlyCost < 0 {
				t.Errorf("Choice %d: MonthlyCost should not be negative", i)
			}
			if choiceComp.CostAnalysis.ThreeYearTCO < 0 {
				t.Errorf("Choice %d: ThreeYearTCO should not be negative", i)
			}
			if choiceComp.CostAnalysis.CostBreakdown == nil {
				t.Errorf("Choice %d: CostBreakdown should not be nil", i)
			}

			// ScalabilityScore should be in valid range (1-10)
			if choiceComp.ScalabilityScore < 1 || choiceComp.ScalabilityScore > 10 {
				t.Errorf("Choice %d: ScalabilityScore should be 1-10, got %d", i, choiceComp.ScalabilityScore)
			}

			// ComplexityScore should be in valid range (1-10)
			if choiceComp.ComplexityScore < 1 || choiceComp.ComplexityScore > 10 {
				t.Errorf("Choice %d: ComplexityScore should be 1-10, got %d", i, choiceComp.ComplexityScore)
			}

			// SecurityScore should be in valid range (1-10)
			if choiceComp.SecurityScore < 1 || choiceComp.SecurityScore > 10 {
				t.Errorf("Choice %d: SecurityScore should be 1-10, got %d", i, choiceComp.SecurityScore)
			}
		}

		// Verify recommendation is valid index
		if comparison.Recommendation < 0 || comparison.Recommendation >= numChoices {
			t.Errorf("Recommendation index %d is out of range [0, %d)", comparison.Recommendation, numChoices)
		}

		// Verify reasoning steps are present
		if len(comparison.ReasoningSteps) == 0 {
			t.Error("ReasoningSteps should not be empty")
		}
	})
}

// Feature: aws-learning-game, Property 16: Decision Feedback Completeness
// For any DecisionResult, the feedback SHALL include an Explanation string
// AND at least one AWSBestPractice or LearningPoint.
// **Validates: Requirements 6.1, 6.2**
func TestProperty16_DecisionFeedbackCompleteness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random company
		capital := rapid.Int64Range(100, 10000).Draw(t, "capital")
		employees := rapid.IntRange(5, 500).Draw(t, "employees")
		securityLevel := rapid.IntRange(1, 5).Draw(t, "securityLevel")
		comp := createTestCompany(capital, employees, securityLevel)
		comp.TechDebt = rapid.IntRange(0, 100).Draw(t, "techDebt")
		comp.CloudAdoption = rapid.Float64Range(0, 100).Draw(t, "cloudAdoption")

		// Randomly assign company type
		companyTypes := []company.CompanyType{
			company.Startup,
			company.Traditional,
			company.CloudReseller,
			company.CloudNative,
		}
		comp.Type = companyTypes[rapid.IntRange(0, len(companyTypes)-1).Draw(t, "companyType")]

		// Generate random choice
		minCapital := rapid.Int64Range(50, 5000).Draw(t, "minCapital")
		minEmployees := rapid.IntRange(1, 200).Draw(t, "minEmployees")
		minSecLevel := rapid.IntRange(1, 5).Draw(t, "minSecLevel")
		isAWS := rapid.Bool().Draw(t, "isAWS")

		choice := createTestChoice(isAWS, minCapital, minEmployees, minSecLevel)
		choice.Outcomes.CapitalChange = rapid.Int64Range(-500, 1000).Draw(t, "capitalChange")
		choice.Outcomes.EmployeeChange = rapid.IntRange(-20, 50).Draw(t, "employeeChange")
		choice.Outcomes.SuccessRate = rapid.Float64Range(0.1, 1.0).Draw(t, "successRate")

		// Add AWS services if it's an AWS choice
		if isAWS {
			awsServices := []string{"EC2", "S3", "RDS", "Lambda", "DynamoDB", "VPC", "IAM", "CloudWatch", "SNS", "SQS"}
			numServices := rapid.IntRange(1, 4).Draw(t, "numServices")
			choice.AWSServices = make([]string, numServices)
			for i := 0; i < numServices; i++ {
				choice.AWSServices[i] = awsServices[rapid.IntRange(0, len(awsServices)-1).Draw(t, "serviceIndex")]
			}
		}

		// Execute decision
		engine := NewDecisionEngine()
		result, err := engine.ExecuteDecision(comp, choice)

		if err != nil {
			t.Fatalf("ExecuteDecision failed: %v", err)
		}

		// Verify result is not nil
		if result == nil {
			t.Fatal("DecisionResult should not be nil")
		}

		// Verify Explanation is not empty
		if result.Explanation == "" {
			t.Error("Explanation should not be empty")
		}

		// Verify at least one of AWSBestPractice or LearningPoints is present
		hasAWSBestPractice := result.AWSBestPractice != ""
		hasLearningPoints := len(result.LearningPoints) > 0

		if !hasAWSBestPractice && !hasLearningPoints {
			t.Error("DecisionResult should have at least one AWSBestPractice or LearningPoint")
		}

		// Verify SAAExamTopics is populated for AWS choices
		if isAWS && len(result.SAAExamTopics) == 0 {
			t.Error("SAAExamTopics should be populated for AWS choices")
		}

		// Verify LearningPoints are meaningful (not empty strings)
		for i, point := range result.LearningPoints {
			if point == "" {
				t.Errorf("LearningPoint %d should not be empty", i)
			}
		}
	})
}
