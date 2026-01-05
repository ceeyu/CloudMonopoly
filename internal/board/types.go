package board

// CellType 格子類型
type CellType string

const (
	CellNormal      CellType = "normal"      // 一般格
	CellOpportunity CellType = "opportunity" // 機會
	CellFate        CellType = "fate"        // 命運
	CellChallenge   CellType = "challenge"   // 關卡
	CellStart       CellType = "start"       // 起點
	CellBonus       CellType = "bonus"       // 獎勵格
)

// Cell 棋盤格子
type Cell struct {
	Position      int      `json:"position"`
	Type          CellType `json:"type"`
	Name          string   `json:"name"`
	BaseCapital   int64    `json:"base_capital"`   // 基礎資本獎勵
	BaseEmployees int      `json:"base_employees"` // 基礎員工獎勵
	EventID       string   `json:"event_id"`       // 關聯事件ID (可選)
}

// Board 遊戲棋盤
type Board struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Cells []Cell `json:"cells"`
	Size  int    `json:"size"`
}

// BoardConfig 棋盤配置
type BoardConfig struct {
	Size                int
	OpportunityCells    int
	FateCells           int
	ChallengeCells      int
	BonusCells          int
	BaseCapitalPerCell  int64
	BaseEmployeePerCell int
}
