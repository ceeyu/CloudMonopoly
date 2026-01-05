package awscatalog

import (
	"encoding/json"
	"fmt"
	"os"
)

// ServicesJSON JSON 格式的服務資料結構
type ServicesJSON struct {
	Version       string        `json:"version"`
	Description   string        `json:"description"`
	TotalServices int           `json:"total_services"`
	Services      []*AWSService `json:"services"`
}

// LoadServicesFromJSON 從 JSON 檔案載入服務
func LoadServicesFromJSON(filepath string) ([]*AWSService, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read services file: %w", err)
	}

	var servicesJSON ServicesJSON
	if err := json.Unmarshal(data, &servicesJSON); err != nil {
		return nil, fmt.Errorf("failed to parse services JSON: %w", err)
	}

	return servicesJSON.Services, nil
}

// SaveServicesToJSON 將服務儲存為 JSON 檔案
func SaveServicesToJSON(filepath string, services []*AWSService) error {
	servicesJSON := ServicesJSON{
		Version:       "1.0",
		Description:   "AWS Learning Game 服務目錄資料集",
		TotalServices: len(services),
		Services:      services,
	}

	data, err := json.MarshalIndent(servicesJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal services: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write services file: %w", err)
	}

	return nil
}

// ExportDefaultServicesToJSON 將預設服務匯出為 JSON
func ExportDefaultServicesToJSON(filepath string) error {
	return SaveServicesToJSON(filepath, DefaultAWSServices)
}
