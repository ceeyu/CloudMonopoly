package game

import (
	"github.com/aws-learning-game/internal/company"
)

// VictoryCondition 勝利條件
type VictoryCondition struct {
	CompanyType         company.CompanyType `json:"company_type"`
	TargetCapital       int64               `json:"target_capital"`        // 目標資本 (0 表示不檢查)
	TargetEmployees     int                 `json:"target_employees"`      // 目標員工數 (0 表示不檢查)
	TargetCloudAdoption float64             `json:"target_cloud_adoption"` // 目標雲端採用率 (0 表示不檢查)
	TargetSecurityLevel int                 `json:"target_security_level"` // 目標資安等級 (0 表示不檢查)
}

// DefaultVictoryConditions 預設勝利條件
// Requirements 1.1, 1.2, 1.3, 1.4
var DefaultVictoryConditions = map[company.CompanyType]VictoryCondition{
	// Startup: 資本達到 3000 萬
	company.Startup: {
		CompanyType:   company.Startup,
		TargetCapital: 3000,
	},
	// Traditional: 雲端採用率達到 80%
	company.Traditional: {
		CompanyType:         company.Traditional,
		TargetCloudAdoption: 80.0,
	},
	// CloudReseller: 員工數達到 150
	company.CloudReseller: {
		CompanyType:     company.CloudReseller,
		TargetEmployees: 150,
	},
	// CloudNative: 資本達到 2000 萬 且 資安等級達到 5
	company.CloudNative: {
		CompanyType:         company.CloudNative,
		TargetCapital:       2000,
		TargetSecurityLevel: 5,
	},
}

// CheckVictory 檢查玩家是否達成勝利條件
// Requirements 1.1, 1.2, 1.3, 1.4
func CheckVictory(player *PlayerState) bool {
	if player == nil || player.Company == nil {
		return false
	}

	condition, ok := DefaultVictoryConditions[player.Company.Type]
	if !ok {
		return false
	}

	return checkCondition(player.Company, condition)
}

// checkCondition 檢查公司是否滿足勝利條件
func checkCondition(comp *company.Company, condition VictoryCondition) bool {
	switch condition.CompanyType {
	case company.Startup:
		// Startup: 資本 >= 3000
		return comp.Capital >= condition.TargetCapital

	case company.Traditional:
		// Traditional: 雲端採用率 >= 80%
		return comp.CloudAdoption >= condition.TargetCloudAdoption

	case company.CloudReseller:
		// CloudReseller: 員工數 >= 150
		return comp.Employees >= condition.TargetEmployees

	case company.CloudNative:
		// CloudNative: 資本 >= 2000 且 資安等級 >= 5
		return comp.Capital >= condition.TargetCapital &&
			comp.SecurityLevel >= condition.TargetSecurityLevel

	default:
		return false
	}
}

// CalculateVictoryProgress 計算勝利進度百分比
// Requirements 3.1, 3.2, 3.3, 3.4, 3.5
func CalculateVictoryProgress(player *PlayerState) float64 {
	if player == nil || player.Company == nil {
		return 0.0
	}

	condition, ok := DefaultVictoryConditions[player.Company.Type]
	if !ok {
		return 0.0
	}

	progress := calculateProgressForType(player.Company, condition)

	// 確保進度上限為 100%
	if progress > 100.0 {
		progress = 100.0
	}
	if progress < 0.0 {
		progress = 0.0
	}

	return progress
}

// calculateProgressForType 根據公司類型計算進度
func calculateProgressForType(comp *company.Company, condition VictoryCondition) float64 {
	switch condition.CompanyType {
	case company.Startup:
		// Startup: (current_capital / 3000) * 100%
		if condition.TargetCapital <= 0 {
			return 0.0
		}
		return (float64(comp.Capital) / float64(condition.TargetCapital)) * 100.0

	case company.Traditional:
		// Traditional: (current_cloud_adoption / 80) * 100%
		if condition.TargetCloudAdoption <= 0 {
			return 0.0
		}
		return (comp.CloudAdoption / condition.TargetCloudAdoption) * 100.0

	case company.CloudReseller:
		// CloudReseller: (current_employees / 150) * 100%
		if condition.TargetEmployees <= 0 {
			return 0.0
		}
		return (float64(comp.Employees) / float64(condition.TargetEmployees)) * 100.0

	case company.CloudNative:
		// CloudNative: ((capital_progress + security_progress) / 2) * 100%
		// capital_progress = min(current_capital / 2000, 1)
		// security_progress = min(current_security_level / 5, 1)
		capitalProgress := float64(comp.Capital) / float64(condition.TargetCapital)
		if capitalProgress > 1.0 {
			capitalProgress = 1.0
		}

		securityProgress := float64(comp.SecurityLevel) / float64(condition.TargetSecurityLevel)
		if securityProgress > 1.0 {
			securityProgress = 1.0
		}

		return ((capitalProgress + securityProgress) / 2.0) * 100.0

	default:
		return 0.0
	}
}

// GetVictoryConditionDescription 取得勝利條件描述
func GetVictoryConditionDescription(companyType company.CompanyType) string {
	switch companyType {
	case company.Startup:
		return "資本達到 3000 萬"
	case company.Traditional:
		return "雲端採用率達到 80%"
	case company.CloudReseller:
		return "員工數達到 150 人"
	case company.CloudNative:
		return "資本達到 2000 萬 且 資安等級達到 5"
	default:
		return "未知勝利條件"
	}
}

// DetermineWinnerByProgress 根據進度判定贏家 (回合限制時使用)
// Requirements 2.3: 回合限制時根據進度判定贏家
func DetermineWinnerByProgress(game *Game) *PlayerState {
	if game == nil || len(game.Players) == 0 {
		return nil
	}

	var winner *PlayerState
	var maxProgress float64 = -1

	for _, player := range game.Players {
		// 計算最新進度
		progress := CalculateVictoryProgress(player)
		player.VictoryProgress = progress

		if progress > maxProgress {
			maxProgress = progress
			winner = player
		}
	}

	return winner
}
