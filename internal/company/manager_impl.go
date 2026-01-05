package company

import (
	"fmt"

	"github.com/aws-learning-game/pkg/models"
)

// generateID 產生唯一 ID
func (m *CompanyManager) generateID() string {
	m.idCounter++
	return fmt.Sprintf("company_%d", m.idCounter)
}

// CreateCompany 建立公司
func (m *CompanyManager) CreateCompany(companyType CompanyType) (*Company, error) {
	defaults, ok := CompanyDefaults[companyType]
	if !ok {
		return nil, ErrInvalidCompanyType
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	company := &Company{
		ID:              m.generateID(),
		Name:            string(companyType) + "_company",
		Type:            companyType,
		Capital:         defaults.Capital,
		Employees:       defaults.Employees,
		IsInternational: false,
		ProductCycle:    defaults.ProductCycle,
		TechDebt:        defaults.TechDebt,
		SecurityLevel:   defaults.SecurityLevel,
		CloudAdoption:   defaults.CloudAdoption,
		Infrastructure:  []string{},
	}

	m.companies[company.ID] = company
	return company, nil
}

// UpdateCapital 更新資本
func (m *CompanyManager) UpdateCapital(companyID string, delta int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	company, ok := m.companies[companyID]
	if !ok {
		return ErrCompanyNotFound
	}

	newCapital := company.Capital + delta
	if newCapital < 0 {
		newCapital = 0
	}
	company.Capital = newCapital
	return nil
}

// UpdateEmployees 更新員工數
func (m *CompanyManager) UpdateEmployees(companyID string, delta int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	company, ok := m.companies[companyID]
	if !ok {
		return ErrCompanyNotFound
	}

	newEmployees := company.Employees + delta
	if newEmployees < 0 {
		newEmployees = 0
	}
	company.Employees = newEmployees
	return nil
}

// GetCompanyState 取得公司狀態
func (m *CompanyManager) GetCompanyState(companyID string) (*Company, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	company, ok := m.companies[companyID]
	if !ok {
		return nil, ErrCompanyNotFound
	}

	// 返回副本以避免外部修改
	companyCopy := *company
	companyCopy.Infrastructure = make([]string, len(company.Infrastructure))
	copy(companyCopy.Infrastructure, company.Infrastructure)

	return &companyCopy, nil
}

// CheckDecisionEligibility 檢查是否符合決策條件
func (m *CompanyManager) CheckDecisionEligibility(companyID string, requirements models.ChoiceRequirements) (bool, []string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	company, ok := m.companies[companyID]
	if !ok {
		return false, []string{"company not found"}
	}

	var issues []string

	if company.Capital < requirements.MinCapital {
		issues = append(issues, fmt.Sprintf("insufficient capital: need %d, have %d", requirements.MinCapital, company.Capital))
	}

	if company.Employees < requirements.MinEmployees {
		issues = append(issues, fmt.Sprintf("insufficient employees: need %d, have %d", requirements.MinEmployees, company.Employees))
	}

	if company.SecurityLevel < requirements.MinSecurityLevel {
		issues = append(issues, fmt.Sprintf("insufficient security level: need %d, have %d", requirements.MinSecurityLevel, company.SecurityLevel))
	}

	// 檢查必要基礎設施
	infraSet := make(map[string]bool)
	for _, infra := range company.Infrastructure {
		infraSet[infra] = true
	}
	for _, required := range requirements.RequiredInfra {
		if !infraSet[required] {
			issues = append(issues, fmt.Sprintf("missing required infrastructure: %s", required))
		}
	}

	return len(issues) == 0, issues
}

// SetCompany 設定公司 (用於測試或載入遊戲)
func (m *CompanyManager) SetCompany(company *Company) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.companies[company.ID] = company
}
