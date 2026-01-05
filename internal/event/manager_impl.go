package event

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/aws-learning-game/internal/board"
	"github.com/aws-learning-game/internal/company"
)

var (
	ErrEventNotFound   = errors.New("event not found")
	ErrInvalidChoice   = errors.New("invalid choice")
	ErrNoEventsForType = errors.New("no events available for this type")
)

// eventSystemImpl 事件系統實作
type eventSystemImpl struct {
	events       map[string]*Event
	eventsByType map[EventType][]*Event
	mu           sync.RWMutex
}

// NewEventSystem 建立事件系統
func NewEventSystem() EventSystem {
	es := &eventSystemImpl{
		events:       make(map[string]*Event),
		eventsByType: make(map[EventType][]*Event),
	}
	return es
}

// NewEventSystemWithEvents 建立事件系統並載入事件
func NewEventSystemWithEvents(events []*Event) EventSystem {
	es := &eventSystemImpl{
		events:       make(map[string]*Event),
		eventsByType: make(map[EventType][]*Event),
	}
	for _, e := range events {
		es.events[e.ID] = e
		es.eventsByType[e.Type] = append(es.eventsByType[e.Type], e)
	}
	return es
}

// GetEvent 取得事件
func (es *eventSystemImpl) GetEvent(eventID string) (*Event, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	event, ok := es.events[eventID]
	if !ok {
		return nil, ErrEventNotFound
	}
	return event, nil
}

// GetRandomEvent 隨機取得事件 (依類型)
func (es *eventSystemImpl) GetRandomEvent(eventType EventType, companyContext *company.Company) (*Event, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	events, ok := es.eventsByType[eventType]
	if !ok || len(events) == 0 {
		return nil, ErrNoEventsForType
	}

	// 隨機選擇一個事件
	idx := rand.Intn(len(events))
	return events[idx], nil
}

// GetEventForCell 根據格子類型取得對應事件
func (es *eventSystemImpl) GetEventForCell(cellType board.CellType, companyContext *company.Company) (*Event, error) {
	eventType := CellTypeToEventType(cellType)
	return es.GetRandomEvent(eventType, companyContext)
}

// ProcessEventOutcome 處理事件結果
func (es *eventSystemImpl) ProcessEventOutcome(event *Event, choiceID int, comp *company.Company) (*EventOutcome, error) {
	// 找到選擇的選項
	var choice *EventChoice
	for i := range event.Choices {
		if event.Choices[i].ID == choiceID {
			choice = &event.Choices[i]
			break
		}
	}
	if choice == nil {
		return nil, ErrInvalidChoice
	}

	// 計算成功率 (基於公司屬性)
	successRate := choice.Outcomes.SuccessRate
	success := rand.Float64() < successRate

	outcome := &EventOutcome{
		Success:        success,
		LearningPoints: event.AWSTopics,
	}

	if success {
		outcome.CapitalDelta = choice.Outcomes.CapitalChange
		outcome.EmployeeDelta = choice.Outcomes.EmployeeChange
		outcome.Message = "決策成功！" + choice.Title
		if choice.IsAWS && len(choice.AWSServices) > 0 {
			outcome.AWSBestPractice = "使用 AWS 服務: " + choice.AWSServices[0]
		}
	} else {
		// 失敗時減半收益或產生損失
		outcome.CapitalDelta = choice.Outcomes.CapitalChange / 2
		outcome.EmployeeDelta = choice.Outcomes.EmployeeChange / 2
		outcome.Message = "決策未達預期效果: " + choice.Title
	}

	return outcome, nil
}

// AddEvent 新增事件 (用於測試或動態載入)
func (es *eventSystemImpl) AddEvent(event *Event) {
	es.mu.Lock()
	defer es.mu.Unlock()

	es.events[event.ID] = event
	es.eventsByType[event.Type] = append(es.eventsByType[event.Type], event)
}

// GetSecurityEvent 取得資安事件 (確保包含緩解策略)
func (es *eventSystemImpl) GetSecurityEvent(companyContext *company.Company) (*Event, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	// 先嘗試取得 security 類型事件
	events, ok := es.eventsByType[EventSecurity]
	if ok && len(events) > 0 {
		idx := rand.Intn(len(events))
		return events[idx], nil
	}

	// 如果沒有專門的 security 事件，從 fate 事件中找資安相關的
	fateEvents, ok := es.eventsByType[EventFate]
	if !ok || len(fateEvents) == 0 {
		return nil, ErrNoEventsForType
	}

	// 過濾出資安相關事件
	var securityRelated []*Event
	for _, e := range fateEvents {
		if isSecurityRelatedEvent(e) {
			securityRelated = append(securityRelated, e)
		}
	}

	if len(securityRelated) == 0 {
		return nil, ErrNoEventsForType
	}

	idx := rand.Intn(len(securityRelated))
	return securityRelated[idx], nil
}

// isSecurityRelatedEvent 判斷是否為資安相關事件
func isSecurityRelatedEvent(e *Event) bool {
	// 檢查標題或描述是否包含資安關鍵字
	securityKeywords := []string{"資安", "安全", "駭客", "攻擊", "漏洞", "外洩", "security", "breach", "hack"}
	for _, keyword := range securityKeywords {
		if containsKeyword(e.Title, keyword) || containsKeyword(e.Description, keyword) {
			return true
		}
	}
	return false
}

// containsKeyword 檢查字串是否包含關鍵字
func containsKeyword(s, keyword string) bool {
	return len(s) > 0 && len(keyword) > 0 &&
		(len(s) >= len(keyword) && (s == keyword ||
			findSubstring(s, keyword)))
}

// findSubstring 簡單的子字串搜尋
func findSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// ValidateSecurityEvent 驗證資安事件是否包含至少2個緩解策略
func ValidateSecurityEvent(event *Event) bool {
	if event == nil {
		return false
	}
	if event.Type != EventSecurity {
		return true // 非資安事件不需要驗證
	}
	return len(event.Choices) >= 2
}

// CreateSecurityEvent 建立資安事件 (確保包含緩解策略)
func CreateSecurityEvent(id, title, description string, choices []EventChoice) *Event {
	// 確保至少有2個選項
	if len(choices) < 2 {
		// 如果選項不足，添加預設選項
		defaultChoices := []EventChoice{
			{
				ID:          len(choices) + 1,
				Title:       "緊急修補",
				Description: "立即進行系統修補",
				IsAWS:       true,
				AWSServices: []string{"AWS WAF", "AWS Shield"},
				Outcomes: ChoiceOutcomes{
					CapitalChange:  -100,
					SecurityChange: 2,
					SuccessRate:    0.8,
				},
			},
			{
				ID:          len(choices) + 2,
				Title:       "風險評估後處理",
				Description: "先進行風險評估再決定處理方式",
				IsAWS:       false,
				Outcomes: ChoiceOutcomes{
					CapitalChange:  -50,
					SecurityChange: 1,
					SuccessRate:    0.6,
				},
			},
		}
		choices = append(choices, defaultChoices[:2-len(choices)]...)
	}

	return &Event{
		ID:          id,
		Type:        EventSecurity,
		Title:       title,
		Description: description,
		Context: EventContext{
			Scenario:       "資安事件發生",
			BusinessImpact: "可能造成資料外洩或服務中斷",
		},
		Choices:   choices,
		AWSTopics: []string{"AWS WAF", "AWS Shield", "Security Hub", "GuardDuty"},
	}
}

// EventDisplay 事件顯示資訊
type EventDisplay struct {
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Scenario       string          `json:"scenario"`
	BusinessImpact string          `json:"business_impact"`
	Choices        []ChoiceDisplay `json:"choices"`
	AWSTopics      []string        `json:"aws_topics"`
}

// ChoiceDisplay 選項顯示資訊
type ChoiceDisplay struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	IsAWS       bool     `json:"is_aws"`
	AWSServices []string `json:"aws_services,omitempty"`
}

// GetEventDisplay 取得事件顯示資訊
func GetEventDisplay(event *Event) *EventDisplay {
	if event == nil {
		return nil
	}

	display := &EventDisplay{
		Title:          event.Title,
		Description:    event.Description,
		Scenario:       event.Context.Scenario,
		BusinessImpact: event.Context.BusinessImpact,
		AWSTopics:      event.AWSTopics,
		Choices:        make([]ChoiceDisplay, len(event.Choices)),
	}

	for i, choice := range event.Choices {
		display.Choices[i] = ChoiceDisplay{
			ID:          choice.ID,
			Title:       choice.Title,
			Description: choice.Description,
			IsAWS:       choice.IsAWS,
			AWSServices: choice.AWSServices,
		}
	}

	return display
}

// ValidateEventDisplay 驗證事件顯示完整性
// 確保事件包含: Title, Description, Context (Scenario, BusinessImpact), 至少一個 Choice
func ValidateEventDisplay(event *Event) bool {
	if event == nil {
		return false
	}

	// 檢查 Title
	if event.Title == "" {
		return false
	}

	// 檢查 Description
	if event.Description == "" {
		return false
	}

	// 檢查 Context
	if event.Context.Scenario == "" {
		return false
	}
	if event.Context.BusinessImpact == "" {
		return false
	}

	// 檢查至少有一個 Choice
	if len(event.Choices) < 1 {
		return false
	}

	return true
}

// FormatEventForDisplay 格式化事件為顯示字串
func FormatEventForDisplay(event *Event) string {
	if event == nil {
		return ""
	}

	result := "=== " + event.Title + " ===\n"
	result += event.Description + "\n\n"
	result += "情境: " + event.Context.Scenario + "\n"
	result += "商業影響: " + event.Context.BusinessImpact + "\n\n"
	result += "可選方案:\n"

	for _, choice := range event.Choices {
		result += "  [" + string(rune('0'+choice.ID)) + "] " + choice.Title
		if choice.IsAWS {
			result += " (AWS)"
		}
		result += "\n"
	}

	return result
}
