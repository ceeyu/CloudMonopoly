package board

import "fmt"

// CreateBoard 建立棋盤
func (m *BoardManager) CreateBoard(boardType string) (*Board, error) {
	config := DefaultBoardConfig

	if config.Size < 30 {
		return nil, ErrInvalidBoardSize
	}

	cells := m.generateCells(config)

	board := &Board{
		ID:    m.generateID(),
		Name:  boardType,
		Cells: cells,
		Size:  len(cells),
	}

	m.board = board
	return board, nil
}

// CreateBoardWithConfig 使用自訂配置建立棋盤
func (m *BoardManager) CreateBoardWithConfig(boardType string, config BoardConfig) (*Board, error) {
	if config.Size < 30 {
		return nil, ErrInvalidBoardSize
	}

	cells := m.generateCells(config)

	board := &Board{
		ID:    m.generateID(),
		Name:  boardType,
		Cells: cells,
		Size:  len(cells),
	}

	m.board = board
	return board, nil
}

// generateCells 產生棋盤格子
func (m *BoardManager) generateCells(config BoardConfig) []Cell {
	cells := make([]Cell, config.Size)

	// 起點
	cells[0] = Cell{
		Position:      0,
		Type:          CellStart,
		Name:          "起點",
		BaseCapital:   config.BaseCapitalPerCell * 2,
		BaseEmployees: config.BaseEmployeePerCell * 2,
	}

	// 分配特殊格子位置
	specialPositions := m.distributeSpecialCells(config)

	for i := 1; i < config.Size; i++ {
		cellType, ok := specialPositions[i]
		if !ok {
			cellType = CellNormal
		}

		cells[i] = Cell{
			Position:      i,
			Type:          cellType,
			Name:          m.getCellName(cellType, i),
			BaseCapital:   m.getBaseCapital(cellType, config),
			BaseEmployees: m.getBaseEmployees(cellType, config),
		}
	}

	return cells
}

// distributeSpecialCells 分配特殊格子
func (m *BoardManager) distributeSpecialCells(config BoardConfig) map[int]CellType {
	positions := make(map[int]CellType)

	// 均勻分配特殊格子
	totalSpecial := config.OpportunityCells + config.FateCells + config.ChallengeCells + config.BonusCells
	interval := (config.Size - 1) / totalSpecial
	if interval < 1 {
		interval = 1
	}

	pos := interval
	count := 0

	// 機會格
	for i := 0; i < config.OpportunityCells && pos < config.Size; i++ {
		positions[pos] = CellOpportunity
		pos += interval
		count++
	}

	// 命運格
	for i := 0; i < config.FateCells && pos < config.Size; i++ {
		positions[pos] = CellFate
		pos += interval
		count++
	}

	// 關卡格
	for i := 0; i < config.ChallengeCells && pos < config.Size; i++ {
		positions[pos] = CellChallenge
		pos += interval
		count++
	}

	// 獎勵格
	for i := 0; i < config.BonusCells && pos < config.Size; i++ {
		positions[pos] = CellBonus
		pos += interval
		count++
	}

	return positions
}

// getCellName 取得格子名稱
func (m *BoardManager) getCellName(cellType CellType, position int) string {
	switch cellType {
	case CellOpportunity:
		return fmt.Sprintf("機會 %d", position)
	case CellFate:
		return fmt.Sprintf("命運 %d", position)
	case CellChallenge:
		return fmt.Sprintf("關卡 %d", position)
	case CellBonus:
		return fmt.Sprintf("獎勵 %d", position)
	default:
		return fmt.Sprintf("格子 %d", position)
	}
}

// getBaseCapital 取得基礎資本獎勵
func (m *BoardManager) getBaseCapital(cellType CellType, config BoardConfig) int64 {
	switch cellType {
	case CellOpportunity:
		return config.BaseCapitalPerCell * 2
	case CellBonus:
		return config.BaseCapitalPerCell * 3
	case CellChallenge:
		return config.BaseCapitalPerCell / 2
	default:
		return config.BaseCapitalPerCell
	}
}

// getBaseEmployees 取得基礎員工獎勵
func (m *BoardManager) getBaseEmployees(cellType CellType, config BoardConfig) int {
	switch cellType {
	case CellOpportunity:
		return config.BaseEmployeePerCell * 3
	case CellBonus:
		return config.BaseEmployeePerCell * 2
	default:
		return config.BaseEmployeePerCell
	}
}

// GetCell 取得格子資訊
func (m *BoardManager) GetCell(position int) (*Cell, error) {
	if m.board == nil {
		return nil, ErrBoardNotFound
	}

	if position < 0 || position >= m.board.Size {
		return nil, ErrInvalidPosition
	}

	cell := m.board.Cells[position]
	return &cell, nil
}

// CalculateNewPosition 計算新位置 (含繞圈處理)
func (m *BoardManager) CalculateNewPosition(current int, diceValue int) int {
	if m.board == nil {
		return current
	}
	return (current + diceValue) % m.board.Size
}

// GetBoard 取得當前棋盤
func (m *BoardManager) GetBoard() *Board {
	return m.board
}

// SetBoard 設定棋盤 (用於測試或載入遊戲)
func (m *BoardManager) SetBoard(board *Board) {
	m.board = board
}
