package awscatalog

import (
	"testing"
)

func TestLoadServicesFromJSON(t *testing.T) {
	services, err := LoadServicesFromJSON("data/services.json")
	if err != nil {
		t.Fatalf("Failed to load services from JSON: %v", err)
	}

	if len(services) == 0 {
		t.Error("No services loaded from JSON")
	}

	t.Logf("Loaded %d services from JSON", len(services))

	// Verify some key services exist
	serviceIDs := make(map[string]bool)
	for _, s := range services {
		serviceIDs[s.ID] = true
	}

	expectedServices := []string{
		"ec2", "lambda", "s3", "rds", "dynamodb", "vpc", "iam", "cloudwatch",
		"sqs", "sns", "eventbridge", "step_functions", "sagemaker",
	}

	for _, id := range expectedServices {
		if !serviceIDs[id] {
			t.Errorf("Expected service %s not found in JSON", id)
		}
	}

	// Verify all services have required fields
	for _, service := range services {
		if service.ID == "" {
			t.Error("Service has empty ID")
		}
		if service.Name == "" {
			t.Errorf("Service %s has empty Name", service.ID)
		}
		if service.Category == "" {
			t.Errorf("Service %s has empty Category", service.ID)
		}
		if len(service.SAAExamTopics) == 0 {
			t.Errorf("Service %s has no SAA exam topics", service.ID)
		}
	}
}

func TestLoadServicesFromJSON_Categories(t *testing.T) {
	services, err := LoadServicesFromJSON("data/services.json")
	if err != nil {
		t.Fatalf("Failed to load services: %v", err)
	}

	// Count services by category
	categoryCounts := make(map[string]int)
	for _, s := range services {
		categoryCounts[s.Category]++
	}

	// Verify we have services in all expected categories
	expectedCategories := []string{
		"compute", "storage", "database", "networking", "security",
		"analytics", "management", "integration", "machine_learning",
	}

	for _, cat := range expectedCategories {
		if categoryCounts[cat] == 0 {
			t.Errorf("No services found for category: %s", cat)
		} else {
			t.Logf("Category %s: %d services", cat, categoryCounts[cat])
		}
	}
}
