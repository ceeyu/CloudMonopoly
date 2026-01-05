package event

import (
	"testing"

	"github.com/aws-learning-game/internal/board"
	"pgregory.net/rapid"
)

// Feature: aws-learning-game, Property 6: Event Type Cell Matching
// For any cell landing, the returned Event type SHALL match the cell type:
// opportunity cells return opportunity events, fate cells return fate events,
// challenge cells return challenge events.
// **Validates: Requirements 3.1, 3.2, 3.3**
func TestProperty6_EventTypeCellMatching(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Create event system with events of all types
		events := createTestEventsForAllTypes()
		es := NewEventSystemWithEvents(events)

		// Generate random cell type (only event-triggering types)
		cellTypes := []board.CellType{
			board.CellOpportunity,
			board.CellFate,
			board.CellChallenge,
		}
		cellTypeIdx := rapid.IntRange(0, len(cellTypes)-1).Draw(t, "cellTypeIdx")
		cellType := cellTypes[cellTypeIdx]

		// Get event for cell
		event, err := es.GetEventForCell(cellType, nil)
		if err != nil {
			t.Fatalf("GetEventForCell failed: %v", err)
		}

		// Verify event type matches cell type
		expectedEventType := CellTypeToEventType(cellType)
		if event.Type != expectedEventType {
			t.Errorf("Event type mismatch: cell type %s should return event type %s, got %s",
				cellType, expectedEventType, event.Type)
		}
	})
}

// createTestEventsForAllTypes creates test events for all event types
func createTestEventsForAllTypes() []*Event {
	return []*Event{
		{
			ID:          "opp-1",
			Type:        EventOpportunity,
			Title:       "擴廠機會",
			Description: "公司有機會擴展業務",
			Context: EventContext{
				Scenario:       "業務成長需要擴展",
				BusinessImpact: "可增加產能",
			},
			Choices: []EventChoice{
				{ID: 1, Title: "擴展雲端", IsAWS: true, AWSServices: []string{"EC2", "S3"}},
				{ID: 2, Title: "擴展地端", IsAWS: false, OnPremSolution: "購買伺服器"},
			},
			AWSTopics: []string{"EC2", "Auto Scaling"},
		},
		{
			ID:          "opp-2",
			Type:        EventOpportunity,
			Title:       "跨國合作",
			Description: "國際合作機會",
			Context: EventContext{
				Scenario:       "跨國企業合作",
				BusinessImpact: "擴展國際市場",
			},
			Choices: []EventChoice{
				{ID: 1, Title: "使用 AWS Global", IsAWS: true},
			},
			AWSTopics: []string{"CloudFront", "Route 53"},
		},
		{
			ID:          "fate-1",
			Type:        EventFate,
			Title:       "市場變化",
			Description: "市場突然變化",
			Context: EventContext{
				Scenario:       "市場波動",
				BusinessImpact: "可能影響營收",
			},
			Choices: []EventChoice{
				{ID: 1, Title: "調整策略", IsAWS: false},
			},
			AWSTopics: []string{"Cost Explorer"},
		},
		{
			ID:          "fate-2",
			Type:        EventFate,
			Title:       "人才流動",
			Description: "關鍵人才異動",
			Context: EventContext{
				Scenario:       "人員變動",
				BusinessImpact: "影響專案進度",
			},
			Choices: []EventChoice{
				{ID: 1, Title: "招募新人", IsAWS: false},
			},
			AWSTopics: []string{},
		},
		{
			ID:          "challenge-1",
			Type:        EventChallenge,
			Title:       "系統架構挑戰",
			Description: "需要重新設計系統架構",
			Context: EventContext{
				Scenario:       "系統效能瓶頸",
				BusinessImpact: "影響用戶體驗",
			},
			Choices: []EventChoice{
				{ID: 1, Title: "使用 AWS 微服務", IsAWS: true, AWSServices: []string{"ECS", "Lambda"}},
				{ID: 2, Title: "優化現有架構", IsAWS: false},
			},
			AWSTopics: []string{"ECS", "Lambda", "API Gateway"},
		},
		{
			ID:          "challenge-2",
			Type:        EventChallenge,
			Title:       "資料庫擴展",
			Description: "資料庫需要擴展",
			Context: EventContext{
				Scenario:       "資料量成長",
				BusinessImpact: "查詢效能下降",
			},
			Choices: []EventChoice{
				{ID: 1, Title: "遷移到 RDS", IsAWS: true, AWSServices: []string{"RDS", "Aurora"}},
			},
			AWSTopics: []string{"RDS", "Aurora", "DynamoDB"},
		},
	}
}

// Test CellTypeToEventType mapping
func TestCellTypeToEventType(t *testing.T) {
	testCases := []struct {
		cellType     board.CellType
		expectedType EventType
	}{
		{board.CellOpportunity, EventOpportunity},
		{board.CellFate, EventFate},
		{board.CellChallenge, EventChallenge},
		{board.CellNormal, EventOpportunity}, // Default
		{board.CellStart, EventOpportunity},  // Default
		{board.CellBonus, EventOpportunity},  // Default
	}

	for _, tc := range testCases {
		result := CellTypeToEventType(tc.cellType)
		if result != tc.expectedType {
			t.Errorf("CellTypeToEventType(%s) = %s, expected %s",
				tc.cellType, result, tc.expectedType)
		}
	}
}

// Feature: aws-learning-game, Property 7: Security Event Mitigation Choices
// For any security incident event, the Event SHALL contain at least 2 mitigation strategy choices.
// **Validates: Requirements 3.4**
func TestProperty7_SecurityEventMitigationChoices(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random security event parameters
		id := rapid.String().Draw(t, "id")
		title := rapid.StringMatching(`[a-zA-Z\x{4e00}-\x{9fff}]{5,20}`).Draw(t, "title")
		description := rapid.StringMatching(`[a-zA-Z\x{4e00}-\x{9fff}]{10,50}`).Draw(t, "description")

		// Generate random number of initial choices (0-3)
		numChoices := rapid.IntRange(0, 3).Draw(t, "numChoices")
		choices := make([]EventChoice, numChoices)
		for i := 0; i < numChoices; i++ {
			choices[i] = EventChoice{
				ID:    i + 1,
				Title: rapid.String().Draw(t, "choiceTitle"),
			}
		}

		// Create security event using the helper function
		event := CreateSecurityEvent(id, title, description, choices)

		// Verify event type is security
		if event.Type != EventSecurity {
			t.Errorf("Event type should be security, got %s", event.Type)
		}

		// Verify at least 2 mitigation choices
		if len(event.Choices) < 2 {
			t.Errorf("Security event should have at least 2 choices, got %d", len(event.Choices))
		}

		// Verify validation function
		if !ValidateSecurityEvent(event) {
			t.Errorf("ValidateSecurityEvent should return true for valid security event")
		}
	})
}

// Test security event validation
func TestValidateSecurityEvent(t *testing.T) {
	testCases := []struct {
		name     string
		event    *Event
		expected bool
	}{
		{
			name:     "nil event",
			event:    nil,
			expected: false,
		},
		{
			name: "security event with 2 choices",
			event: &Event{
				Type:    EventSecurity,
				Choices: []EventChoice{{ID: 1}, {ID: 2}},
			},
			expected: true,
		},
		{
			name: "security event with 1 choice",
			event: &Event{
				Type:    EventSecurity,
				Choices: []EventChoice{{ID: 1}},
			},
			expected: false,
		},
		{
			name: "non-security event with 1 choice",
			event: &Event{
				Type:    EventOpportunity,
				Choices: []EventChoice{{ID: 1}},
			},
			expected: true, // Non-security events don't need validation
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateSecurityEvent(tc.event)
			if result != tc.expected {
				t.Errorf("ValidateSecurityEvent() = %v, expected %v", result, tc.expected)
			}
		})
	}
}

// Feature: aws-learning-game, Property 8: Event Display Completeness
// For any Event presentation, the output SHALL include: Title, Description,
// Context (Scenario, BusinessImpact), and at least one Choice.
// **Validates: Requirements 3.6**
func TestProperty8_EventDisplayCompleteness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a complete event with all required fields
		event := &Event{
			ID:          rapid.String().Draw(t, "id"),
			Type:        EventOpportunity,
			Title:       rapid.StringMatching(`[a-zA-Z\x{4e00}-\x{9fff}]{5,20}`).Draw(t, "title"),
			Description: rapid.StringMatching(`[a-zA-Z\x{4e00}-\x{9fff}]{10,50}`).Draw(t, "description"),
			Context: EventContext{
				Scenario:       rapid.StringMatching(`[a-zA-Z\x{4e00}-\x{9fff}]{5,30}`).Draw(t, "scenario"),
				BusinessImpact: rapid.StringMatching(`[a-zA-Z\x{4e00}-\x{9fff}]{5,30}`).Draw(t, "impact"),
			},
			Choices: []EventChoice{
				{
					ID:    1,
					Title: rapid.StringMatching(`[a-zA-Z\x{4e00}-\x{9fff}]{3,15}`).Draw(t, "choiceTitle"),
				},
			},
		}

		// Validate event display completeness
		if !ValidateEventDisplay(event) {
			t.Errorf("ValidateEventDisplay should return true for complete event")
		}

		// Get event display
		display := GetEventDisplay(event)
		if display == nil {
			t.Fatal("GetEventDisplay should not return nil for valid event")
		}

		// Verify all required fields are present
		if display.Title == "" {
			t.Error("Display should have Title")
		}
		if display.Description == "" {
			t.Error("Display should have Description")
		}
		if display.Scenario == "" {
			t.Error("Display should have Scenario")
		}
		if display.BusinessImpact == "" {
			t.Error("Display should have BusinessImpact")
		}
		if len(display.Choices) < 1 {
			t.Error("Display should have at least one Choice")
		}
	})
}

// Test ValidateEventDisplay with incomplete events
func TestValidateEventDisplay_Incomplete(t *testing.T) {
	testCases := []struct {
		name     string
		event    *Event
		expected bool
	}{
		{
			name:     "nil event",
			event:    nil,
			expected: false,
		},
		{
			name: "missing title",
			event: &Event{
				Title:       "",
				Description: "desc",
				Context:     EventContext{Scenario: "s", BusinessImpact: "b"},
				Choices:     []EventChoice{{ID: 1}},
			},
			expected: false,
		},
		{
			name: "missing description",
			event: &Event{
				Title:       "title",
				Description: "",
				Context:     EventContext{Scenario: "s", BusinessImpact: "b"},
				Choices:     []EventChoice{{ID: 1}},
			},
			expected: false,
		},
		{
			name: "missing scenario",
			event: &Event{
				Title:       "title",
				Description: "desc",
				Context:     EventContext{Scenario: "", BusinessImpact: "b"},
				Choices:     []EventChoice{{ID: 1}},
			},
			expected: false,
		},
		{
			name: "missing business impact",
			event: &Event{
				Title:       "title",
				Description: "desc",
				Context:     EventContext{Scenario: "s", BusinessImpact: ""},
				Choices:     []EventChoice{{ID: 1}},
			},
			expected: false,
		},
		{
			name: "no choices",
			event: &Event{
				Title:       "title",
				Description: "desc",
				Context:     EventContext{Scenario: "s", BusinessImpact: "b"},
				Choices:     []EventChoice{},
			},
			expected: false,
		},
		{
			name: "complete event",
			event: &Event{
				Title:       "title",
				Description: "desc",
				Context:     EventContext{Scenario: "s", BusinessImpact: "b"},
				Choices:     []EventChoice{{ID: 1}},
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateEventDisplay(tc.event)
			if result != tc.expected {
				t.Errorf("ValidateEventDisplay() = %v, expected %v", result, tc.expected)
			}
		})
	}
}

// Test GetEventDisplay
func TestGetEventDisplay(t *testing.T) {
	event := &Event{
		ID:          "test-1",
		Type:        EventOpportunity,
		Title:       "Test Event",
		Description: "Test Description",
		Context: EventContext{
			Scenario:       "Test Scenario",
			BusinessImpact: "Test Impact",
		},
		Choices: []EventChoice{
			{ID: 1, Title: "Choice 1", IsAWS: true, AWSServices: []string{"EC2"}},
			{ID: 2, Title: "Choice 2", IsAWS: false},
		},
		AWSTopics: []string{"EC2", "S3"},
	}

	display := GetEventDisplay(event)

	if display.Title != event.Title {
		t.Errorf("Title mismatch: got %s, expected %s", display.Title, event.Title)
	}
	if display.Description != event.Description {
		t.Errorf("Description mismatch")
	}
	if display.Scenario != event.Context.Scenario {
		t.Errorf("Scenario mismatch")
	}
	if display.BusinessImpact != event.Context.BusinessImpact {
		t.Errorf("BusinessImpact mismatch")
	}
	if len(display.Choices) != len(event.Choices) {
		t.Errorf("Choices count mismatch")
	}
	if len(display.AWSTopics) != len(event.AWSTopics) {
		t.Errorf("AWSTopics count mismatch")
	}
}

// Test that default events dataset has at least 50 events
func TestDefaultEventsCount(t *testing.T) {
	if len(DefaultEvents) < 50 {
		t.Errorf("DefaultEvents should have at least 50 events, got %d", len(DefaultEvents))
	}
}

// Test that all default events are valid for display
func TestDefaultEventsValidity(t *testing.T) {
	for _, event := range DefaultEvents {
		if !ValidateEventDisplay(event) {
			t.Errorf("Event %s is not valid for display", event.ID)
		}
	}
}

// Test that security events have at least 2 choices
func TestDefaultSecurityEventsHaveMitigationChoices(t *testing.T) {
	for _, event := range DefaultEvents {
		if event.Type == EventSecurity {
			if len(event.Choices) < 2 {
				t.Errorf("Security event %s should have at least 2 choices, got %d", event.ID, len(event.Choices))
			}
		}
	}
}

// Test GetDefaultEventSystem
func TestGetDefaultEventSystem(t *testing.T) {
	es := GetDefaultEventSystem()
	if es == nil {
		t.Fatal("GetDefaultEventSystem should not return nil")
	}

	// Test getting an event
	event, err := es.GetEvent("opp-001")
	if err != nil {
		t.Errorf("GetEvent failed: %v", err)
	}
	if event == nil {
		t.Error("Event should not be nil")
	}
	if event.Title != "擴廠機會" {
		t.Errorf("Event title mismatch: got %s", event.Title)
	}
}

// Test that default events cover all required scenarios
func TestDefaultEventsCoverageScenarios(t *testing.T) {
	// Required scenarios from Requirements 3.5, 3.7
	// Each scenario can be matched by multiple keywords
	requiredScenarios := map[string][]string{
		"擴廠": {"擴廠", "擴展", "擴充"},
		"跨國": {"跨國", "國際", "全球", "東南亞"},
		"資安": {"資安", "安全", "Security"},
		"外洩": {"外洩", "洩漏", "breach"},
		"故障": {"故障", "中斷", "停擺"},
		"成長": {"成長", "成長", "增加"},
		"合規": {"合規", "認證", "SOC", "GDPR"},
		"攻擊": {"攻擊", "DDoS", "駭客", "注入"},
		"災難": {"災難", "復原", "備援", "備份"},
	}

	scenarioCovered := make(map[string]bool)
	for scenario := range requiredScenarios {
		scenarioCovered[scenario] = false
	}

	for _, event := range DefaultEvents {
		// Check all text fields for scenario keywords
		text := event.Title + event.Description + event.Context.Scenario +
			event.Context.BusinessImpact + event.RealWorldRef
		for _, need := range event.Context.TechnicalNeeds {
			text += need
		}
		for _, topic := range event.AWSTopics {
			text += topic
		}
		for _, choice := range event.Choices {
			text += choice.Title + choice.Description
		}

		for scenario, keywords := range requiredScenarios {
			for _, keyword := range keywords {
				if containsKeyword(text, keyword) {
					scenarioCovered[scenario] = true
					break
				}
			}
		}
	}

	// Verify all scenarios are covered
	for scenario, covered := range scenarioCovered {
		if !covered {
			t.Errorf("Required scenario '%s' is not covered by any event", scenario)
		}
	}
}

// Test that all event types have events
func TestDefaultEventsHaveAllTypes(t *testing.T) {
	typeCounts := map[EventType]int{
		EventOpportunity: 0,
		EventFate:        0,
		EventChallenge:   0,
		EventSecurity:    0,
	}

	for _, event := range DefaultEvents {
		typeCounts[event.Type]++
	}

	// Verify each type has at least some events
	for eventType, count := range typeCounts {
		if count == 0 {
			t.Errorf("Event type %s has no events", eventType)
		}
		t.Logf("Event type %s: %d events", eventType, count)
	}
}

// Test that events have AWS topics for learning
func TestDefaultEventsHaveAWSTopics(t *testing.T) {
	eventsWithTopics := 0
	for _, event := range DefaultEvents {
		if len(event.AWSTopics) > 0 {
			eventsWithTopics++
		}
	}

	// At least 80% of events should have AWS topics
	minRequired := len(DefaultEvents) * 80 / 100
	if eventsWithTopics < minRequired {
		t.Errorf("At least %d events should have AWS topics, got %d", minRequired, eventsWithTopics)
	}
}

// Test that events have real world references
func TestDefaultEventsHaveRealWorldRefs(t *testing.T) {
	eventsWithRefs := 0
	for _, event := range DefaultEvents {
		if event.RealWorldRef != "" {
			eventsWithRefs++
		}
	}

	// At least 80% of events should have real world references
	minRequired := len(DefaultEvents) * 80 / 100
	if eventsWithRefs < minRequired {
		t.Errorf("At least %d events should have real world references, got %d", minRequired, eventsWithRefs)
	}
}

// Test that each event has both AWS and on-prem options where applicable
func TestDefaultEventsHaveChoiceVariety(t *testing.T) {
	eventsWithBothOptions := 0
	for _, event := range DefaultEvents {
		hasAWS := false
		hasOnPrem := false
		for _, choice := range event.Choices {
			if choice.IsAWS {
				hasAWS = true
			} else {
				hasOnPrem = true
			}
		}
		if hasAWS && hasOnPrem {
			eventsWithBothOptions++
		}
	}

	// At least 80% of events should have both AWS and on-prem options
	minRequired := len(DefaultEvents) * 80 / 100
	if eventsWithBothOptions < minRequired {
		t.Errorf("At least %d events should have both AWS and on-prem options, got %d", minRequired, eventsWithBothOptions)
	}
}
