package architecture

import (
	"testing"

	"github.com/aws-learning-game/internal/company"
	"github.com/aws-learning-game/internal/event"
	"pgregory.net/rapid"
)

// createTestChoice creates a test event choice for testing
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

// createTestCompany creates a test company for testing
func createTestCompany(capital int64, employees int, securityLevel int) *company.Company {
	return &company.Company{
		ID:             "test-company",
		Name:           "Test Company",
		Type:           company.Startup,
		Capital:        capital,
		Employees:      employees,
		SecurityLevel:  securityLevel,
		CloudAdoption:  50.0,
		TechDebt:       20,
		ProductCycle:   "growth",
		Infrastructure: []string{},
	}
}

// Feature: aws-learning-game, Property 14: Architecture Diagram Generation
// For any EventChoice, the ArchitectureVisualizer SHALL produce a non-empty diagram string.
// **Validates: Requirements 5.1**
func TestProperty14_ArchitectureDiagramGeneration(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random choice
		isAWS := rapid.Bool().Draw(t, "isAWS")
		minCapital := rapid.Int64Range(100, 5000).Draw(t, "minCapital")
		minEmployees := rapid.IntRange(5, 200).Draw(t, "minEmployees")
		minSecLevel := rapid.IntRange(1, 5).Draw(t, "minSecLevel")

		choice := &event.EventChoice{
			ID:          rapid.IntRange(1, 100).Draw(t, "choiceID"),
			Title:       rapid.StringN(5, 50, -1).Draw(t, "title"),
			Description: rapid.StringN(10, 200, -1).Draw(t, "description"),
			IsAWS:       isAWS,
			Requirements: event.ChoiceRequirements{
				MinCapital:       minCapital,
				MinEmployees:     minEmployees,
				MinSecurityLevel: minSecLevel,
			},
			Outcomes: event.ChoiceOutcomes{
				CapitalChange:   rapid.Int64Range(-500, 1000).Draw(t, "capitalChange"),
				EmployeeChange:  rapid.IntRange(-10, 50).Draw(t, "employeeChange"),
				SuccessRate:     rapid.Float64Range(0.1, 1.0).Draw(t, "successRate"),
				TimeToImplement: rapid.IntRange(1, 10).Draw(t, "timeToImpl"),
			},
		}

		if isAWS {
			// Generate random AWS services
			numServices := rapid.IntRange(1, 5).Draw(t, "numServices")
			services := []string{"EC2", "S3", "RDS", "Lambda", "DynamoDB", "ECS", "EKS", "CloudFront"}
			choice.AWSServices = make([]string, numServices)
			for i := 0; i < numServices; i++ {
				idx := rapid.IntRange(0, len(services)-1).Draw(t, "serviceIdx")
				choice.AWSServices[i] = services[idx]
			}
		} else {
			choice.OnPremSolution = rapid.StringN(10, 100, -1).Draw(t, "onPremSolution")
		}

		// Generate diagram
		visualizer := NewArchitectureVisualizer()
		diagram, err := visualizer.GenerateDiagram(choice)

		if err != nil {
			t.Fatalf("GenerateDiagram failed: %v", err)
		}

		// Verify diagram is not empty
		if diagram == "" {
			t.Error("Diagram should not be empty")
		}

		// Verify diagram contains expected content
		if isAWS {
			if !containsString(diagram, "AWS") {
				t.Error("AWS diagram should contain 'AWS'")
			}
		} else {
			if !containsString(diagram, "地端") {
				t.Error("On-prem diagram should contain '地端'")
			}
		}
	})
}

// Test diagram generation with different formats
func TestDiagramFormats(t *testing.T) {
	visualizer := NewArchitectureVisualizer()

	testCases := []struct {
		name   string
		format DiagramFormat
		isAWS  bool
	}{
		{"ASCII AWS", FormatASCII, true},
		{"ASCII OnPrem", FormatASCII, false},
		{"Mermaid AWS", FormatMermaid, true},
		{"Mermaid OnPrem", FormatMermaid, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			choice := createTestChoice(tc.isAWS, 500, 20, 3)

			opts := DiagramOptions{
				Format:      tc.format,
				ShowDetails: true,
			}

			diagram, err := visualizer.GenerateDiagramWithOptions(choice, opts)
			if err != nil {
				t.Fatalf("GenerateDiagramWithOptions failed: %v", err)
			}

			if diagram == "" {
				t.Error("Diagram should not be empty")
			}

			// Check format-specific content
			if tc.format == FormatMermaid {
				if !containsString(diagram, "```mermaid") {
					t.Error("Mermaid diagram should contain mermaid code block")
				}
				if !containsString(diagram, "graph TB") {
					t.Error("Mermaid diagram should contain graph definition")
				}
			} else if tc.format == FormatASCII {
				if !containsString(diagram, "===") {
					t.Error("ASCII diagram should contain title markers")
				}
			}
		})
	}
}

// Test comparison table generation
func TestComparisonTableGeneration(t *testing.T) {
	visualizer := NewArchitectureVisualizer()

	choices := []event.EventChoice{
		*createTestChoice(true, 500, 20, 3),
		*createTestChoice(false, 1000, 30, 2),
	}
	choices[0].Title = "AWS 方案"
	choices[1].Title = "地端方案"

	table, err := visualizer.GenerateComparisonTable(choices)
	if err != nil {
		t.Fatalf("GenerateComparisonTable failed: %v", err)
	}

	if table == "" {
		t.Error("Comparison table should not be empty")
	}

	// Verify table contains expected headers
	expectedHeaders := []string{"項目", "類型", "初始成本", "員工需求", "資安等級", "實施時間", "成功率", "資本變化"}
	for _, header := range expectedHeaders {
		if !containsString(table, header) {
			t.Errorf("Comparison table should contain header '%s'", header)
		}
	}

	// Verify table contains choice titles
	if !containsString(table, "AWS 方案") {
		t.Error("Comparison table should contain 'AWS 方案'")
	}
	if !containsString(table, "地端方案") {
		t.Error("Comparison table should contain '地端方案'")
	}
}

// Test company architecture generation
func TestCompanyArchitectureGeneration(t *testing.T) {
	visualizer := NewArchitectureVisualizer()

	comp := createTestCompany(1000, 50, 3)
	comp.Name = "測試公司"
	comp.Infrastructure = []string{"EC2", "S3", "RDS"}

	diagram, err := visualizer.GenerateCompanyArchitecture(comp)
	if err != nil {
		t.Fatalf("GenerateCompanyArchitecture failed: %v", err)
	}

	if diagram == "" {
		t.Error("Company architecture diagram should not be empty")
	}

	// Verify diagram contains company info
	if !containsString(diagram, "測試公司") {
		t.Error("Diagram should contain company name")
	}

	// Verify diagram contains infrastructure
	for _, infra := range comp.Infrastructure {
		if !containsString(diagram, infra) {
			t.Errorf("Diagram should contain infrastructure '%s'", infra)
		}
	}
}

// Test error handling
func TestErrorHandling(t *testing.T) {
	visualizer := NewArchitectureVisualizer()

	t.Run("nil choice", func(t *testing.T) {
		_, err := visualizer.GenerateDiagram(nil)
		if err == nil {
			t.Error("Should return error for nil choice")
		}
	})

	t.Run("nil company", func(t *testing.T) {
		_, err := visualizer.GenerateCompanyArchitecture(nil)
		if err == nil {
			t.Error("Should return error for nil company")
		}
	})

	t.Run("empty choices", func(t *testing.T) {
		_, err := visualizer.GenerateComparisonTable([]event.EventChoice{})
		if err == nil {
			t.Error("Should return error for empty choices")
		}
	})
}

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
