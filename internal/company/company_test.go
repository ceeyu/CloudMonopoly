package company

import (
	"testing"

	"pgregory.net/rapid"
)

// Feature: aws-learning-game, Property 1: Company Initialization Completeness
// For any valid company type selection, the created Company SHALL have all predefined
// attributes (Capital, Employees, SecurityLevel, CloudAdoption) set to non-zero default
// values matching the CompanyDefaults configuration.
// **Validates: Requirements 1.2**
func TestProperty1_CompanyInitializationCompleteness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random valid company type
		companyTypes := []CompanyType{Startup, Traditional, CloudReseller, CloudNative}
		typeIndex := rapid.IntRange(0, len(companyTypes)-1).Draw(t, "typeIndex")
		companyType := companyTypes[typeIndex]

		// Create company
		manager := NewCompanyManager()
		company, err := manager.CreateCompany(companyType)

		// Verify no error
		if err != nil {
			t.Fatalf("CreateCompany failed: %v", err)
		}

		// Get expected defaults
		defaults := CompanyDefaults[companyType]

		// Verify all attributes are set to non-zero default values
		if company.Capital <= 0 {
			t.Errorf("Capital should be positive, got %d", company.Capital)
		}
		if company.Capital != defaults.Capital {
			t.Errorf("Capital mismatch: expected %d, got %d", defaults.Capital, company.Capital)
		}

		if company.Employees <= 0 {
			t.Errorf("Employees should be positive, got %d", company.Employees)
		}
		if company.Employees != defaults.Employees {
			t.Errorf("Employees mismatch: expected %d, got %d", defaults.Employees, company.Employees)
		}

		if company.SecurityLevel <= 0 {
			t.Errorf("SecurityLevel should be positive, got %d", company.SecurityLevel)
		}
		if company.SecurityLevel != defaults.SecurityLevel {
			t.Errorf("SecurityLevel mismatch: expected %d, got %d", defaults.SecurityLevel, company.SecurityLevel)
		}

		if company.CloudAdoption <= 0 {
			t.Errorf("CloudAdoption should be positive, got %f", company.CloudAdoption)
		}
		if company.CloudAdoption != defaults.CloudAdoption {
			t.Errorf("CloudAdoption mismatch: expected %f, got %f", defaults.CloudAdoption, company.CloudAdoption)
		}

		// Verify company type matches
		if company.Type != companyType {
			t.Errorf("Type mismatch: expected %s, got %s", companyType, company.Type)
		}

		// Verify ID is set
		if company.ID == "" {
			t.Error("Company ID should not be empty")
		}

		// Verify ProductCycle is set
		if company.ProductCycle == "" {
			t.Error("ProductCycle should not be empty")
		}
		if company.ProductCycle != defaults.ProductCycle {
			t.Errorf("ProductCycle mismatch: expected %s, got %s", defaults.ProductCycle, company.ProductCycle)
		}
	})
}

// Test invalid company type
func TestCreateCompany_InvalidType(t *testing.T) {
	manager := NewCompanyManager()
	_, err := manager.CreateCompany(CompanyType("invalid"))
	if err != ErrInvalidCompanyType {
		t.Errorf("Expected ErrInvalidCompanyType, got %v", err)
	}
}

// Test UpdateCapital
func TestUpdateCapital(t *testing.T) {
	manager := NewCompanyManager()
	company, _ := manager.CreateCompany(Startup)

	initialCapital := company.Capital

	// Test positive delta
	err := manager.UpdateCapital(company.ID, 100)
	if err != nil {
		t.Errorf("UpdateCapital failed: %v", err)
	}

	updated, _ := manager.GetCompanyState(company.ID)
	if updated.Capital != initialCapital+100 {
		t.Errorf("Expected capital %d, got %d", initialCapital+100, updated.Capital)
	}

	// Test negative delta (should not go below 0)
	err = manager.UpdateCapital(company.ID, -10000)
	if err != nil {
		t.Errorf("UpdateCapital failed: %v", err)
	}

	updated, _ = manager.GetCompanyState(company.ID)
	if updated.Capital < 0 {
		t.Errorf("Capital should not be negative, got %d", updated.Capital)
	}
}

// Test UpdateEmployees
func TestUpdateEmployees(t *testing.T) {
	manager := NewCompanyManager()
	company, _ := manager.CreateCompany(Startup)

	initialEmployees := company.Employees

	// Test positive delta
	err := manager.UpdateEmployees(company.ID, 5)
	if err != nil {
		t.Errorf("UpdateEmployees failed: %v", err)
	}

	updated, _ := manager.GetCompanyState(company.ID)
	if updated.Employees != initialEmployees+5 {
		t.Errorf("Expected employees %d, got %d", initialEmployees+5, updated.Employees)
	}

	// Test negative delta (should not go below 0)
	err = manager.UpdateEmployees(company.ID, -10000)
	if err != nil {
		t.Errorf("UpdateEmployees failed: %v", err)
	}

	updated, _ = manager.GetCompanyState(company.ID)
	if updated.Employees < 0 {
		t.Errorf("Employees should not be negative, got %d", updated.Employees)
	}
}

// Test company not found errors
func TestCompanyNotFound(t *testing.T) {
	manager := NewCompanyManager()

	_, err := manager.GetCompanyState("nonexistent")
	if err != ErrCompanyNotFound {
		t.Errorf("Expected ErrCompanyNotFound, got %v", err)
	}

	err = manager.UpdateCapital("nonexistent", 100)
	if err != ErrCompanyNotFound {
		t.Errorf("Expected ErrCompanyNotFound, got %v", err)
	}

	err = manager.UpdateEmployees("nonexistent", 5)
	if err != ErrCompanyNotFound {
		t.Errorf("Expected ErrCompanyNotFound, got %v", err)
	}
}
