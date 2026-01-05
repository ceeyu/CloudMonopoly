package architecture

import (
	"fmt"
	"strings"

	"github.com/aws-learning-game/internal/company"
	"github.com/aws-learning-game/internal/event"
)

// architectureVisualizerImpl 架構視覺化實作
type architectureVisualizerImpl struct{}

// NewArchitectureVisualizer 建立新的架構視覺化器
func NewArchitectureVisualizer() ArchitectureVisualizer {
	return &architectureVisualizerImpl{}
}

// GenerateDiagram 產生架構圖（使用預設選項）
func (a *architectureVisualizerImpl) GenerateDiagram(choice *event.EventChoice) (string, error) {
	return a.GenerateDiagramWithOptions(choice, DefaultDiagramOptions())
}

// GenerateDiagramWithOptions 產生架構圖（含選項）
func (a *architectureVisualizerImpl) GenerateDiagramWithOptions(choice *event.EventChoice, opts DiagramOptions) (string, error) {
	if choice == nil {
		return "", fmt.Errorf("choice cannot be nil")
	}

	switch opts.Format {
	case FormatASCII:
		return a.generateASCIIDiagram(choice, opts.ShowDetails)
	case FormatMermaid:
		return a.generateMermaidDiagram(choice, opts.ShowDetails)
	default:
		return a.generateMermaidDiagram(choice, opts.ShowDetails)
	}
}

// generateMermaidDiagram 產生 Mermaid 格式架構圖
func (a *architectureVisualizerImpl) generateMermaidDiagram(choice *event.EventChoice, showDetails bool) (string, error) {
	var sb strings.Builder

	sb.WriteString("```mermaid\n")
	sb.WriteString("graph TB\n")

	if choice.IsAWS {
		sb.WriteString("    subgraph AWS[\"AWS Cloud\"]\n")
		for i, svc := range choice.AWSServices {
			nodeID := fmt.Sprintf("svc%d", i)
			sb.WriteString(fmt.Sprintf("        %s[\"%s\"]\n", nodeID, svc))
		}
		sb.WriteString("    end\n")

		// 連接服務
		if len(choice.AWSServices) > 1 {
			for i := 0; i < len(choice.AWSServices)-1; i++ {
				sb.WriteString(fmt.Sprintf("    svc%d --> svc%d\n", i, i+1))
			}
		}

		// 添加用戶端
		sb.WriteString("    User[\"使用者\"] --> svc0\n")
	} else {
		sb.WriteString("    subgraph OnPrem[\"地端機房\"]\n")
		sb.WriteString("        Server[\"伺服器\"]\n")
		sb.WriteString("        Storage[\"儲存設備\"]\n")
		sb.WriteString("        Network[\"網路設備\"]\n")
		sb.WriteString("    end\n")
		sb.WriteString("    User[\"使用者\"] --> Server\n")
		sb.WriteString("    Server --> Storage\n")
		sb.WriteString("    Server --> Network\n")
	}

	if showDetails {
		sb.WriteString("\n    %% 方案資訊\n")
		sb.WriteString(fmt.Sprintf("    %% 標題: %s\n", choice.Title))
		sb.WriteString(fmt.Sprintf("    %% 最低資本需求: %d 萬元\n", choice.Requirements.MinCapital))
		sb.WriteString(fmt.Sprintf("    %% 實施時間: %d 回合\n", choice.Outcomes.TimeToImplement))
	}

	sb.WriteString("```")

	return sb.String(), nil
}

// generateASCIIDiagram 產生 ASCII 格式架構圖
func (a *architectureVisualizerImpl) generateASCIIDiagram(choice *event.EventChoice, showDetails bool) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("=== %s ===\n\n", choice.Title))

	if choice.IsAWS {
		sb.WriteString("┌─────────────────────────────────────┐\n")
		sb.WriteString("│           AWS Cloud                 │\n")
		sb.WriteString("├─────────────────────────────────────┤\n")

		for _, svc := range choice.AWSServices {
			sb.WriteString(fmt.Sprintf("│  [%s]%s│\n", svc, strings.Repeat(" ", 33-len(svc))))
		}

		sb.WriteString("└─────────────────────────────────────┘\n")
		sb.WriteString("              ▲\n")
		sb.WriteString("              │\n")
		sb.WriteString("         [使用者]\n")
	} else {
		sb.WriteString("┌─────────────────────────────────────┐\n")
		sb.WriteString("│          地端機房                   │\n")
		sb.WriteString("├─────────────────────────────────────┤\n")
		sb.WriteString("│  [伺服器] ─── [儲存設備]            │\n")
		sb.WriteString("│      │                              │\n")
		sb.WriteString("│  [網路設備]                         │\n")
		sb.WriteString("└─────────────────────────────────────┘\n")
		sb.WriteString("              ▲\n")
		sb.WriteString("              │\n")
		sb.WriteString("         [使用者]\n")
	}

	if showDetails {
		sb.WriteString("\n--- 方案資訊 ---\n")
		sb.WriteString(fmt.Sprintf("最低資本需求: %d 萬元\n", choice.Requirements.MinCapital))
		sb.WriteString(fmt.Sprintf("最低員工數: %d 人\n", choice.Requirements.MinEmployees))
		sb.WriteString(fmt.Sprintf("實施時間: %d 回合\n", choice.Outcomes.TimeToImplement))
		sb.WriteString(fmt.Sprintf("成功率: %.0f%%\n", choice.Outcomes.SuccessRate*100))
	}

	return sb.String(), nil
}

// GenerateComparisonTable 產生比較表
func (a *architectureVisualizerImpl) GenerateComparisonTable(choices []event.EventChoice) (string, error) {
	if len(choices) == 0 {
		return "", fmt.Errorf("choices cannot be empty")
	}

	var sb strings.Builder

	// 表頭
	sb.WriteString("| 項目 |")
	for _, c := range choices {
		sb.WriteString(fmt.Sprintf(" %s |", truncateString(c.Title, 20)))
	}
	sb.WriteString("\n")

	// 分隔線
	sb.WriteString("|------|")
	for range choices {
		sb.WriteString("----------------------|")
	}
	sb.WriteString("\n")

	// 方案類型
	sb.WriteString("| 類型 |")
	for _, c := range choices {
		if c.IsAWS {
			sb.WriteString(" AWS 雲端 |")
		} else {
			sb.WriteString(" 地端方案 |")
		}
	}
	sb.WriteString("\n")

	// 初始成本
	sb.WriteString("| 初始成本 |")
	for _, c := range choices {
		sb.WriteString(fmt.Sprintf(" %d 萬元 |", c.Requirements.MinCapital))
	}
	sb.WriteString("\n")

	// 員工需求
	sb.WriteString("| 員工需求 |")
	for _, c := range choices {
		sb.WriteString(fmt.Sprintf(" %d 人 |", c.Requirements.MinEmployees))
	}
	sb.WriteString("\n")

	// 資安等級
	sb.WriteString("| 資安等級 |")
	for _, c := range choices {
		sb.WriteString(fmt.Sprintf(" Level %d |", c.Requirements.MinSecurityLevel))
	}
	sb.WriteString("\n")

	// 實施時間
	sb.WriteString("| 實施時間 |")
	for _, c := range choices {
		sb.WriteString(fmt.Sprintf(" %d 回合 |", c.Outcomes.TimeToImplement))
	}
	sb.WriteString("\n")

	// 成功率
	sb.WriteString("| 成功率 |")
	for _, c := range choices {
		sb.WriteString(fmt.Sprintf(" %.0f%% |", c.Outcomes.SuccessRate*100))
	}
	sb.WriteString("\n")

	// 資本變化
	sb.WriteString("| 資本變化 |")
	for _, c := range choices {
		if c.Outcomes.CapitalChange >= 0 {
			sb.WriteString(fmt.Sprintf(" +%d 萬元 |", c.Outcomes.CapitalChange))
		} else {
			sb.WriteString(fmt.Sprintf(" %d 萬元 |", c.Outcomes.CapitalChange))
		}
	}
	sb.WriteString("\n")

	// AWS 服務
	sb.WriteString("| AWS 服務 |")
	for _, c := range choices {
		if c.IsAWS && len(c.AWSServices) > 0 {
			services := strings.Join(c.AWSServices, ", ")
			sb.WriteString(fmt.Sprintf(" %s |", truncateString(services, 18)))
		} else {
			sb.WriteString(" - |")
		}
	}
	sb.WriteString("\n")

	return sb.String(), nil
}

// GenerateCompanyArchitecture 產生公司當前架構
func (a *architectureVisualizerImpl) GenerateCompanyArchitecture(comp *company.Company) (string, error) {
	if comp == nil {
		return "", fmt.Errorf("company cannot be nil")
	}

	var sb strings.Builder

	sb.WriteString("```mermaid\n")
	sb.WriteString("graph TB\n")
	sb.WriteString(fmt.Sprintf("    subgraph Company[\"%s - %s\"]\n", comp.Name, getCompanyTypeName(comp.Type)))

	// 公司基本資訊
	sb.WriteString("        Info[\"公司資訊<br/>")
	sb.WriteString(fmt.Sprintf("資本: %d 萬元<br/>", comp.Capital))
	sb.WriteString(fmt.Sprintf("員工: %d 人<br/>", comp.Employees))
	sb.WriteString(fmt.Sprintf("資安等級: %d<br/>", comp.SecurityLevel))
	sb.WriteString(fmt.Sprintf("雲端採用率: %.0f%%\"]\n", comp.CloudAdoption))

	// 基礎設施
	if len(comp.Infrastructure) > 0 {
		sb.WriteString("        subgraph Infra[\"已部署基礎設施\"]\n")
		for i, infra := range comp.Infrastructure {
			nodeID := fmt.Sprintf("infra%d", i)
			sb.WriteString(fmt.Sprintf("            %s[\"%s\"]\n", nodeID, infra))
		}
		sb.WriteString("        end\n")
		sb.WriteString("        Info --> Infra\n")
	} else {
		sb.WriteString("        NoInfra[\"尚無部署基礎設施\"]\n")
		sb.WriteString("        Info --> NoInfra\n")
	}

	sb.WriteString("    end\n")

	// 產品週期狀態
	sb.WriteString(fmt.Sprintf("    Cycle[\"產品週期: %s\"]\n", getProductCycleName(comp.ProductCycle)))
	sb.WriteString("    Company --> Cycle\n")

	// 技術債指標
	if comp.TechDebt > 0 {
		sb.WriteString(fmt.Sprintf("    TechDebt[\"技術債: %d\"]\n", comp.TechDebt))
		sb.WriteString("    Company --> TechDebt\n")
	}

	sb.WriteString("```")

	return sb.String(), nil
}

// truncateString 截斷字串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// getCompanyTypeName 取得公司類型名稱
func getCompanyTypeName(t company.CompanyType) string {
	switch t {
	case company.Startup:
		return "新創公司"
	case company.Traditional:
		return "傳產公司"
	case company.CloudReseller:
		return "雲端代理商"
	case company.CloudNative:
		return "雲端公司"
	default:
		return "未知類型"
	}
}

// getProductCycleName 取得產品週期名稱
func getProductCycleName(cycle string) string {
	switch cycle {
	case "development":
		return "開發期"
	case "launch":
		return "上市期"
	case "growth":
		return "成長期"
	case "mature":
		return "成熟期"
	default:
		return cycle
	}
}
