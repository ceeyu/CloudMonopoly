package architecture

import (
	"github.com/aws-learning-game/internal/company"
	"github.com/aws-learning-game/internal/event"
)

// ArchitectureVisualizer 架構視覺化介面
type ArchitectureVisualizer interface {
	// GenerateDiagram 產生架構圖
	GenerateDiagram(choice *event.EventChoice) (string, error)
	// GenerateDiagramWithOptions 產生架構圖（含選項）
	GenerateDiagramWithOptions(choice *event.EventChoice, opts DiagramOptions) (string, error)
	// GenerateComparisonTable 產生比較表
	GenerateComparisonTable(choices []event.EventChoice) (string, error)
	// GenerateCompanyArchitecture 產生公司當前架構
	GenerateCompanyArchitecture(comp *company.Company) (string, error)
}
