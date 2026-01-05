package event

import (
	"github.com/aws-learning-game/internal/board"
	"github.com/aws-learning-game/internal/company"
)

// EventSystem 事件系統介面
type EventSystem interface {
	// GetEvent 取得事件
	GetEvent(eventID string) (*Event, error)
	// GetRandomEvent 隨機取得事件 (依類型)
	GetRandomEvent(eventType EventType, companyContext *company.Company) (*Event, error)
	// GetEventForCell 根據格子類型取得對應事件
	GetEventForCell(cellType board.CellType, companyContext *company.Company) (*Event, error)
	// ProcessEventOutcome 處理事件結果
	ProcessEventOutcome(event *Event, choiceID int, comp *company.Company) (*EventOutcome, error)
}

// CellTypeToEventType 格子類型對應事件類型
func CellTypeToEventType(cellType board.CellType) EventType {
	switch cellType {
	case board.CellOpportunity:
		return EventOpportunity
	case board.CellFate:
		return EventFate
	case board.CellChallenge:
		return EventChallenge
	default:
		return EventOpportunity // 預設為機會事件
	}
}
