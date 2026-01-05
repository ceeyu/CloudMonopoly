package board

import (
	"errors"
	"fmt"
)

var (
	ErrBoardNotFound     = errors.New("board not found")
	ErrInvalidPosition   = errors.New("invalid position")
	ErrInvalidBoardSize  = errors.New("board size must be at least 30")
)

// Manager 棋盤管理介面
type Manager interface {
	// CreateBoard 建立棋盤
	CreateBoard(boardType string) (*Board, error)
	// GetCell 取得格子資訊
	GetCell(position int) (*Cell, error)
	// CalculateNewPosition 計算新位置
	CalculateNewPosition(current int, diceValue int) int
}

// BoardManager 棋盤管理器實作
type BoardManager struct {
	board     *Board
	idCounter int
}

// NewBoardManager 建立新的棋盤管理器
func NewBoardManager() *BoardManager {
	return &BoardManager{
		idCounter: 0,
	}
}

// DefaultBoardConfig 預設棋盤配置
var DefaultBoardConfig = BoardConfig{
	Size:                36,
	OpportunityCells:    6,
	FateCells:           6,
	ChallengeCells:      4,
	BonusCells:          2,
	BaseCapitalPerCell:  50,
	BaseEmployeePerCell: 1,
}

// generateID 產生唯一 ID
func (m *BoardManager) generateID() string {
	m.idCounter++
	return fmt.Sprintf("board_%d", m.idCounter)
}
