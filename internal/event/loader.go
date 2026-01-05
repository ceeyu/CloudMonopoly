package event

import (
	"encoding/json"
	"fmt"
	"os"
)

// EventsJSON JSON 格式的事件資料結構
type EventsJSON struct {
	Version     string   `json:"version"`
	Description string   `json:"description"`
	TotalEvents int      `json:"total_events"`
	Events      []*Event `json:"events"`
}

// LoadEventsFromJSON 從 JSON 檔案載入事件
func LoadEventsFromJSON(filepath string) ([]*Event, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read events file: %w", err)
	}

	var eventsJSON EventsJSON
	if err := json.Unmarshal(data, &eventsJSON); err != nil {
		return nil, fmt.Errorf("failed to parse events JSON: %w", err)
	}

	return eventsJSON.Events, nil
}

// SaveEventsToJSON 將事件儲存為 JSON 檔案
func SaveEventsToJSON(filepath string, events []*Event) error {
	eventsJSON := EventsJSON{
		Version:     "1.0",
		Description: "AWS Learning Game 事件資料集",
		TotalEvents: len(events),
		Events:      events,
	}

	data, err := json.MarshalIndent(eventsJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write events file: %w", err)
	}

	return nil
}
