package models

// ChoiceRequirements 選項需求條件
type ChoiceRequirements struct {
	MinCapital       int64    `json:"min_capital"`
	MinEmployees     int      `json:"min_employees"`
	MinSecurityLevel int      `json:"min_security_level"`
	RequiredInfra    []string `json:"required_infra"`
}
