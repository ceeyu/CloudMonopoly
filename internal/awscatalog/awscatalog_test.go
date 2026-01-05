package awscatalog

import (
	"testing"

	"github.com/aws-learning-game/internal/company"
	"pgregory.net/rapid"
)

// Feature: aws-learning-game, Property 10: AWS Service Category Relevance
// For any scenario-based service query, all returned AWSService items SHALL have
// a Category matching the scenario's technical requirements.
// **Validates: Requirements 4.2**
func TestProperty10_AWSServiceCategoryRelevance(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Create catalog with default services
		catalog := NewAWSServiceCatalogWithServices(DefaultAWSServices)

		// Generate random scenario from known scenarios
		scenarios := []string{
			"compute", "storage", "database", "networking", "security",
			"high_availability", "disaster_recovery", "cost_optimization", "migration",
		}
		scenarioIdx := rapid.IntRange(0, len(scenarios)-1).Draw(t, "scenarioIdx")
		scenario := scenarios[scenarioIdx]

		// Get expected categories for this scenario
		expectedCategories := getRelevantCategories(scenario)
		if len(expectedCategories) == 0 {
			// Skip if no categories mapped
			return
		}

		// Create a test company
		comp := &company.Company{
			ID:            "test-company",
			Type:          company.Startup,
			Capital:       1000,
			Employees:     20,
			SecurityLevel: 3,
			CloudAdoption: 50,
		}

		// Get recommended services
		services, err := catalog.GetRecommendedServices(scenario, comp)
		if err != nil {
			t.Fatalf("GetRecommendedServices failed: %v", err)
		}

		// Verify all returned services have categories matching the scenario's requirements
		for _, service := range services {
			categoryMatches := false
			for _, expectedCat := range expectedCategories {
				if service.Category == expectedCat {
					categoryMatches = true
					break
				}
			}
			if !categoryMatches {
				t.Errorf("Service %s has category %s which does not match scenario %s expected categories %v",
					service.Name, service.Category, scenario, expectedCategories)
			}
		}
	})
}

// Test GetService
func TestGetService(t *testing.T) {
	catalog := NewAWSServiceCatalogWithServices(DefaultAWSServices)

	// Test getting existing service
	service, err := catalog.GetService("ec2")
	if err != nil {
		t.Errorf("GetService failed: %v", err)
	}
	if service == nil {
		t.Fatal("Service should not be nil")
	}
	if service.Name != "Amazon EC2" {
		t.Errorf("Service name mismatch: got %s, expected Amazon EC2", service.Name)
	}

	// Test getting non-existing service
	_, err = catalog.GetService("non-existing")
	if err != ErrServiceNotFound {
		t.Errorf("Expected ErrServiceNotFound, got %v", err)
	}
}

// Test GetServicesByCategory
func TestGetServicesByCategory(t *testing.T) {
	catalog := NewAWSServiceCatalogWithServices(DefaultAWSServices)

	// Test getting compute services
	services, err := catalog.GetServicesByCategory("compute")
	if err != nil {
		t.Errorf("GetServicesByCategory failed: %v", err)
	}
	if len(services) == 0 {
		t.Error("Should have compute services")
	}

	// Verify all returned services are compute category
	for _, s := range services {
		if s.Category != "compute" {
			t.Errorf("Service %s has category %s, expected compute", s.Name, s.Category)
		}
	}

	// Test getting non-existing category
	services, err = catalog.GetServicesByCategory("non-existing")
	if err != nil {
		t.Errorf("GetServicesByCategory should not error for non-existing category: %v", err)
	}
	if len(services) != 0 {
		t.Errorf("Should return empty slice for non-existing category, got %d services", len(services))
	}
}

// Test GetRecommendedServices
func TestGetRecommendedServices(t *testing.T) {
	catalog := NewAWSServiceCatalogWithServices(DefaultAWSServices)

	comp := &company.Company{
		ID:            "test-company",
		Type:          company.Startup,
		Capital:       1000,
		Employees:     20,
		SecurityLevel: 3,
		CloudAdoption: 50,
	}

	// Test compute scenario
	services, err := catalog.GetRecommendedServices("compute", comp)
	if err != nil {
		t.Errorf("GetRecommendedServices failed: %v", err)
	}
	if len(services) == 0 {
		t.Error("Should have recommended services for compute scenario")
	}

	// Verify all services are compute category
	for _, s := range services {
		if s.Category != "compute" {
			t.Errorf("Service %s has category %s, expected compute for compute scenario", s.Name, s.Category)
		}
	}

	// Test high_availability scenario (should return multiple categories)
	services, err = catalog.GetRecommendedServices("high_availability", comp)
	if err != nil {
		t.Errorf("GetRecommendedServices failed: %v", err)
	}
	if len(services) == 0 {
		t.Error("Should have recommended services for high_availability scenario")
	}

	// Verify services are from expected categories
	expectedCategories := map[string]bool{"compute": true, "database": true, "networking": true}
	for _, s := range services {
		if !expectedCategories[s.Category] {
			t.Errorf("Service %s has unexpected category %s for high_availability scenario", s.Name, s.Category)
		}
	}
}

// Test GetRecommendedServices with nil company
func TestGetRecommendedServices_NilCompany(t *testing.T) {
	catalog := NewAWSServiceCatalogWithServices(DefaultAWSServices)

	services, err := catalog.GetRecommendedServices("compute", nil)
	if err != nil {
		t.Errorf("GetRecommendedServices should not fail with nil company: %v", err)
	}
	if len(services) == 0 {
		t.Error("Should have recommended services even with nil company")
	}
}

// Test GetRecommendedServices with unknown scenario
func TestGetRecommendedServices_UnknownScenario(t *testing.T) {
	catalog := NewAWSServiceCatalogWithServices(DefaultAWSServices)

	comp := &company.Company{
		ID:   "test-company",
		Type: company.Startup,
	}

	services, err := catalog.GetRecommendedServices("unknown_scenario", comp)
	if err != nil {
		t.Errorf("GetRecommendedServices should not error for unknown scenario: %v", err)
	}
	// Unknown scenario should return empty slice
	if len(services) != 0 {
		t.Logf("Unknown scenario returned %d services (may match keywords)", len(services))
	}
}

// Test default catalog
func TestGetDefaultCatalog(t *testing.T) {
	catalog := GetDefaultCatalog()
	if catalog == nil {
		t.Fatal("GetDefaultCatalog should not return nil")
	}

	// Test getting a service
	service, err := catalog.GetService("s3")
	if err != nil {
		t.Errorf("GetService failed: %v", err)
	}
	if service == nil {
		t.Error("Service should not be nil")
	}
	if service.Name != "Amazon S3" {
		t.Errorf("Service name mismatch: got %s", service.Name)
	}
}

// Test that default services cover all categories
func TestDefaultServicesCoverAllCategories(t *testing.T) {
	expectedCategories := []string{
		"compute", "storage", "database", "networking", "security", "analytics", "management",
	}

	categoryCounts := make(map[string]int)
	for _, service := range DefaultAWSServices {
		categoryCounts[service.Category]++
	}

	for _, cat := range expectedCategories {
		if categoryCounts[cat] == 0 {
			t.Errorf("Category %s has no services", cat)
		}
		t.Logf("Category %s: %d services", cat, categoryCounts[cat])
	}
}

// Test that all default services have required fields
func TestDefaultServicesHaveRequiredFields(t *testing.T) {
	for _, service := range DefaultAWSServices {
		if service.ID == "" {
			t.Errorf("Service has empty ID")
		}
		if service.Name == "" {
			t.Errorf("Service %s has empty Name", service.ID)
		}
		if service.Category == "" {
			t.Errorf("Service %s has empty Category", service.ID)
		}
		if service.Description == "" {
			t.Errorf("Service %s has empty Description", service.ID)
		}
		if len(service.UseCases) == 0 {
			t.Errorf("Service %s has no UseCases", service.ID)
		}
		if len(service.SAAExamTopics) == 0 {
			t.Errorf("Service %s has no SAAExamTopics", service.ID)
		}
		if len(service.BestPractices) == 0 {
			t.Errorf("Service %s has no BestPractices", service.ID)
		}
	}
}

// Test that services have SAA exam topics
func TestDefaultServicesHaveSAATopics(t *testing.T) {
	servicesWithTopics := 0
	for _, service := range DefaultAWSServices {
		if len(service.SAAExamTopics) > 0 {
			servicesWithTopics++
		}
	}

	// All services should have SAA topics
	if servicesWithTopics != len(DefaultAWSServices) {
		t.Errorf("All services should have SAA exam topics, got %d/%d",
			servicesWithTopics, len(DefaultAWSServices))
	}
}

// Test scenario category mapping
func TestScenarioCategoryMapping(t *testing.T) {
	testCases := []struct {
		scenario           string
		expectedCategories []string
	}{
		{"compute", []string{"compute"}},
		{"storage", []string{"storage"}},
		{"database", []string{"database"}},
		{"networking", []string{"networking"}},
		{"security", []string{"security"}},
		{"high_availability", []string{"compute", "database", "networking"}},
		{"disaster_recovery", []string{"storage", "database", "networking"}},
		{"cost_optimization", []string{"compute", "storage"}},
		{"migration", []string{"compute", "storage", "database"}},
	}

	for _, tc := range testCases {
		categories := getRelevantCategories(tc.scenario)
		if len(categories) != len(tc.expectedCategories) {
			t.Errorf("Scenario %s: expected %d categories, got %d",
				tc.scenario, len(tc.expectedCategories), len(categories))
			continue
		}

		for i, cat := range categories {
			if cat != tc.expectedCategories[i] {
				t.Errorf("Scenario %s: expected category %s at index %d, got %s",
					tc.scenario, tc.expectedCategories[i], i, cat)
			}
		}
	}
}

// Test keyword-based scenario matching
func TestKeywordBasedScenarioMatching(t *testing.T) {
	testCases := []struct {
		scenario           string
		expectedCategories []string
	}{
		{"server_deployment", []string{"compute"}},
		{"s3_backup", []string{"storage"}},
		{"rds_setup", []string{"database"}},
		{"vpc_configuration", []string{"networking"}},
		{"iam_policy", []string{"security"}},
	}

	for _, tc := range testCases {
		categories := getRelevantCategories(tc.scenario)
		if len(categories) == 0 {
			t.Errorf("Scenario %s should match some categories", tc.scenario)
			continue
		}

		// Check if at least one expected category is present
		found := false
		for _, expected := range tc.expectedCategories {
			for _, got := range categories {
				if got == expected {
					found = true
					break
				}
			}
		}
		if !found {
			t.Errorf("Scenario %s: expected one of %v, got %v",
				tc.scenario, tc.expectedCategories, categories)
		}
	}
}
