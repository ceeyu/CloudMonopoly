package company

import (
	"errors"
	"sync"

	"github.com/aws-learning-game/pkg/models"
)

var (
	ErrCompanyNotFound    = errors.New("company not found")
	ErrInvalidCompanyType = errors.New("invalid company type")
	ErrInsufficientFunds  = errors.New("insufficient funds")
)

// Manager 公司管理介面
type Manager interface {
	// CreateCompany 建立公司
	CreateCompany(companyType CompanyType) (*Company, error)
	// UpdateCapital 更新資本
	UpdateCapital(companyID string, delta int64) error
	// UpdateEmployees 更新員工數
	UpdateEmployees(companyID string, delta int) error
	// GetCompanyState 取得公司狀態
	GetCompanyState(companyID string) (*Company, error)
	// CheckDecisionEligibility 檢查是否符合決策條件
	CheckDecisionEligibility(companyID string, requirements models.ChoiceRequirements) (bool, []string)
}

// CompanyManager 公司管理器實作
type CompanyManager struct {
	companies map[string]*Company
	mu        sync.RWMutex
	idCounter int
}

// NewCompanyManager 建立新的公司管理器
func NewCompanyManager() *CompanyManager {
	return &CompanyManager{
		companies: make(map[string]*Company),
		idCounter: 0,
	}
}
