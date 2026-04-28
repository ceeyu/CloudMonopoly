package game

import (
	"testing"

	"github.com/aws-learning-game/internal/company"
)

// TestCheckVictory_Startup tests Startup victory condition
// Requirements 1.1: WHEN a Startup player's capital reaches 3000 (萬), THE Victory_Condition_System SHALL declare that player as the winner
func TestCheckVictory_Startup(t *testing.T) {
	tests := []struct {
		name     string
		capital  int64
		expected bool
	}{
		{"below threshold", 2999, false},
		{"at threshold", 3000, true},
		{"above threshold", 3500, true},
		{"zero capital", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &PlayerState{
				PlayerID: "test_player",
				Company: &company.Company{
					Type:    company.Startup,
					Capital: tt.capital,
				},
			}

			result := CheckVictory(player)
			if result != tt.expected {
				t.Errorf("CheckVictory() = %v, want %v for capital %d", result, tt.expected, tt.capital)
			}
		})
	}
}

// TestCheckVictory_Traditional tests Traditional victory condition
// Requirements 1.2: WHEN a Traditional player's cloud adoption rate reaches 80%, THE Victory_Condition_System SHALL declare that player as the winner
func TestCheckVictory_Traditional(t *testing.T) {
	tests := []struct {
		name          string
		cloudAdoption float64
		expected      bool
	}{
		{"below threshold", 79.9, false},
		{"at threshold", 80.0, true},
		{"above threshold", 90.0, true},
		{"zero adoption", 0.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &PlayerState{
				PlayerID: "test_player",
				Company: &company.Company{
					Type:          company.Traditional,
					CloudAdoption: tt.cloudAdoption,
				},
			}

			result := CheckVictory(player)
			if result != tt.expected {
				t.Errorf("CheckVictory() = %v, want %v for cloudAdoption %f", result, tt.expected, tt.cloudAdoption)
			}
		})
	}
}

// TestCheckVictory_CloudReseller tests CloudReseller victory condition
// Requirements 1.3: WHEN a CloudReseller player's employee count reaches 150, THE Victory_Condition_System SHALL declare that player as the winner
func TestCheckVictory_CloudReseller(t *testing.T) {
	tests := []struct {
		name      string
		employees int
		expected  bool
	}{
		{"below threshold", 149, false},
		{"at threshold", 150, true},
		{"above threshold", 200, true},
		{"zero employees", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &PlayerState{
				PlayerID: "test_player",
				Company: &company.Company{
					Type:      company.CloudReseller,
					Employees: tt.employees,
				},
			}

			result := CheckVictory(player)
			if result != tt.expected {
				t.Errorf("CheckVictory() = %v, want %v for employees %d", result, tt.expected, tt.employees)
			}
		})
	}
}

// TestCheckVictory_CloudNative tests CloudNative victory condition
// Requirements 1.4: WHEN a CloudNative player's capital reaches 2000 (萬) AND security level reaches 5, THE Victory_Condition_System SHALL declare that player as the winner
func TestCheckVictory_CloudNative(t *testing.T) {
	tests := []struct {
		name          string
		capital       int64
		securityLevel int
		expected      bool
	}{
		{"both below threshold", 1999, 4, false},
		{"capital below, security at", 1999, 5, false},
		{"capital at, security below", 2000, 4, false},
		{"both at threshold", 2000, 5, true},
		{"both above threshold", 2500, 6, true},
		{"capital above, security at", 2500, 5, true},
		{"capital at, security above", 2000, 6, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &PlayerState{
				PlayerID: "test_player",
				Company: &company.Company{
					Type:          company.CloudNative,
					Capital:       tt.capital,
					SecurityLevel: tt.securityLevel,
				},
			}

			result := CheckVictory(player)
			if result != tt.expected {
				t.Errorf("CheckVictory() = %v, want %v for capital %d, securityLevel %d", result, tt.expected, tt.capital, tt.securityLevel)
			}
		})
	}
}

// TestCheckVictory_NilPlayer tests nil player handling
func TestCheckVictory_NilPlayer(t *testing.T) {
	result := CheckVictory(nil)
	if result != false {
		t.Errorf("CheckVictory(nil) = %v, want false", result)
	}
}

// TestCheckVictory_NilCompany tests nil company handling
func TestCheckVictory_NilCompany(t *testing.T) {
	player := &PlayerState{
		PlayerID: "test_player",
		Company:  nil,
	}

	result := CheckVictory(player)
	if result != false {
		t.Errorf("CheckVictory() with nil company = %v, want false", result)
	}
}

// TestCalculateVictoryProgress_Startup tests Startup progress calculation
// Requirements 3.1: FOR Startup players, THE Victory_Progress SHALL calculate as (current_capital / 3000) * 100%
func TestCalculateVictoryProgress_Startup(t *testing.T) {
	tests := []struct {
		name     string
		capital  int64
		expected float64
	}{
		{"zero capital", 0, 0.0},
		{"half progress", 1500, 50.0},
		{"at threshold", 3000, 100.0},
		{"above threshold (capped)", 4500, 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &PlayerState{
				PlayerID: "test_player",
				Company: &company.Company{
					Type:    company.Startup,
					Capital: tt.capital,
				},
			}

			result := CalculateVictoryProgress(player)
			if result != tt.expected {
				t.Errorf("CalculateVictoryProgress() = %v, want %v for capital %d", result, tt.expected, tt.capital)
			}
		})
	}
}

// TestCalculateVictoryProgress_Traditional tests Traditional progress calculation
// Requirements 3.2: FOR Traditional players, THE Victory_Progress SHALL calculate as (current_cloud_adoption / 80) * 100%
func TestCalculateVictoryProgress_Traditional(t *testing.T) {
	tests := []struct {
		name          string
		cloudAdoption float64
		expected      float64
	}{
		{"zero adoption", 0.0, 0.0},
		{"half progress", 40.0, 50.0},
		{"at threshold", 80.0, 100.0},
		{"above threshold (capped)", 100.0, 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &PlayerState{
				PlayerID: "test_player",
				Company: &company.Company{
					Type:          company.Traditional,
					CloudAdoption: tt.cloudAdoption,
				},
			}

			result := CalculateVictoryProgress(player)
			if result != tt.expected {
				t.Errorf("CalculateVictoryProgress() = %v, want %v for cloudAdoption %f", result, tt.expected, tt.cloudAdoption)
			}
		})
	}
}

// TestCalculateVictoryProgress_CloudReseller tests CloudReseller progress calculation
// Requirements 3.3: FOR CloudReseller players, THE Victory_Progress SHALL calculate as (current_employees / 150) * 100%
func TestCalculateVictoryProgress_CloudReseller(t *testing.T) {
	tests := []struct {
		name      string
		employees int
		expected  float64
	}{
		{"zero employees", 0, 0.0},
		{"half progress", 75, 50.0},
		{"at threshold", 150, 100.0},
		{"above threshold (capped)", 200, 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &PlayerState{
				PlayerID: "test_player",
				Company: &company.Company{
					Type:      company.CloudReseller,
					Employees: tt.employees,
				},
			}

			result := CalculateVictoryProgress(player)
			if result != tt.expected {
				t.Errorf("CalculateVictoryProgress() = %v, want %v for employees %d", result, tt.expected, tt.employees)
			}
		})
	}
}

// TestCalculateVictoryProgress_CloudNative tests CloudNative progress calculation
// Requirements 3.4: FOR CloudNative players, THE Victory_Progress SHALL calculate as ((capital_progress + security_progress) / 2) * 100%
func TestCalculateVictoryProgress_CloudNative(t *testing.T) {
	tests := []struct {
		name          string
		capital       int64
		securityLevel int
		expected      float64
	}{
		{"zero both", 0, 0, 0.0},
		{"half capital, zero security", 1000, 0, 25.0},
		{"zero capital, half security", 0, 2, 20.0},
		{"half both", 1000, 2, 45.0},             // (0.5 + 0.4) / 2 * 100 = 45
		{"at threshold both", 2000, 5, 100.0},    // (1.0 + 1.0) / 2 * 100 = 100
		{"above threshold both", 3000, 6, 100.0}, // capped at 100
		{"capital above, security at", 3000, 5, 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &PlayerState{
				PlayerID: "test_player",
				Company: &company.Company{
					Type:          company.CloudNative,
					Capital:       tt.capital,
					SecurityLevel: tt.securityLevel,
				},
			}

			result := CalculateVictoryProgress(player)
			if result != tt.expected {
				t.Errorf("CalculateVictoryProgress() = %v, want %v for capital %d, securityLevel %d", result, tt.expected, tt.capital, tt.securityLevel)
			}
		})
	}
}

// TestCalculateVictoryProgress_NilPlayer tests nil player handling
func TestCalculateVictoryProgress_NilPlayer(t *testing.T) {
	result := CalculateVictoryProgress(nil)
	if result != 0.0 {
		t.Errorf("CalculateVictoryProgress(nil) = %v, want 0.0", result)
	}
}

// TestCalculateVictoryProgress_NilCompany tests nil company handling
func TestCalculateVictoryProgress_NilCompany(t *testing.T) {
	player := &PlayerState{
		PlayerID: "test_player",
		Company:  nil,
	}

	result := CalculateVictoryProgress(player)
	if result != 0.0 {
		t.Errorf("CalculateVictoryProgress() with nil company = %v, want 0.0", result)
	}
}

// TestCalculateVictoryProgress_Cap tests that progress is capped at 100%
// Requirements 3.5: THE Victory_Progress SHALL cap at 100% maximum
func TestCalculateVictoryProgress_Cap(t *testing.T) {
	companyTypes := []struct {
		companyType   company.CompanyType
		capital       int64
		employees     int
		cloudAdoption float64
		securityLevel int
	}{
		{company.Startup, 5000, 0, 0, 0},
		{company.Traditional, 0, 0, 100.0, 0},
		{company.CloudReseller, 0, 300, 0, 0},
		{company.CloudNative, 5000, 0, 0, 10},
	}

	for _, tt := range companyTypes {
		t.Run(string(tt.companyType), func(t *testing.T) {
			player := &PlayerState{
				PlayerID: "test_player",
				Company: &company.Company{
					Type:          tt.companyType,
					Capital:       tt.capital,
					Employees:     tt.employees,
					CloudAdoption: tt.cloudAdoption,
					SecurityLevel: tt.securityLevel,
				},
			}

			result := CalculateVictoryProgress(player)
			if result > 100.0 {
				t.Errorf("CalculateVictoryProgress() = %v, should be capped at 100.0 for %s", result, tt.companyType)
			}
		})
	}
}

// TestGetVictoryConditionDescription tests description generation
func TestGetVictoryConditionDescription(t *testing.T) {
	tests := []struct {
		companyType company.CompanyType
		expected    string
	}{
		{company.Startup, "資本達到 3000 萬"},
		{company.Traditional, "雲端採用率達到 80%"},
		{company.CloudReseller, "員工數達到 150 人"},
		{company.CloudNative, "資本達到 2000 萬 且 資安等級達到 5"},
		{company.CompanyType("unknown"), "未知勝利條件"},
	}

	for _, tt := range tests {
		t.Run(string(tt.companyType), func(t *testing.T) {
			result := GetVictoryConditionDescription(tt.companyType)
			if result != tt.expected {
				t.Errorf("GetVictoryConditionDescription(%s) = %v, want %v", tt.companyType, result, tt.expected)
			}
		})
	}
}

// TestDefaultVictoryConditions tests that default conditions are correctly defined
func TestDefaultVictoryConditions(t *testing.T) {
	// Verify all company types have conditions
	expectedTypes := []company.CompanyType{
		company.Startup,
		company.Traditional,
		company.CloudReseller,
		company.CloudNative,
	}

	for _, ct := range expectedTypes {
		condition, ok := DefaultVictoryConditions[ct]
		if !ok {
			t.Errorf("DefaultVictoryConditions missing condition for %s", ct)
			continue
		}

		if condition.CompanyType != ct {
			t.Errorf("DefaultVictoryConditions[%s].CompanyType = %s, want %s", ct, condition.CompanyType, ct)
		}
	}

	// Verify specific values
	if DefaultVictoryConditions[company.Startup].TargetCapital != 3000 {
		t.Errorf("Startup TargetCapital = %d, want 3000", DefaultVictoryConditions[company.Startup].TargetCapital)
	}

	if DefaultVictoryConditions[company.Traditional].TargetCloudAdoption != 80.0 {
		t.Errorf("Traditional TargetCloudAdoption = %f, want 80.0", DefaultVictoryConditions[company.Traditional].TargetCloudAdoption)
	}

	if DefaultVictoryConditions[company.CloudReseller].TargetEmployees != 150 {
		t.Errorf("CloudReseller TargetEmployees = %d, want 150", DefaultVictoryConditions[company.CloudReseller].TargetEmployees)
	}

	cloudNative := DefaultVictoryConditions[company.CloudNative]
	if cloudNative.TargetCapital != 2000 || cloudNative.TargetSecurityLevel != 5 {
		t.Errorf("CloudNative TargetCapital = %d, TargetSecurityLevel = %d, want 2000 and 5",
			cloudNative.TargetCapital, cloudNative.TargetSecurityLevel)
	}
}
