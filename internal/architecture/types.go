package architecture

// DiagramFormat 架構圖格式
type DiagramFormat string

const (
	FormatASCII   DiagramFormat = "ascii"   // ASCII 文字格式
	FormatMermaid DiagramFormat = "mermaid" // Mermaid 圖表格式
)

// DiagramOptions 架構圖產生選項
type DiagramOptions struct {
	Format      DiagramFormat
	ShowDetails bool // 是否顯示詳細資訊
}

// DefaultDiagramOptions 預設選項
func DefaultDiagramOptions() DiagramOptions {
	return DiagramOptions{
		Format:      FormatMermaid,
		ShowDetails: true,
	}
}
