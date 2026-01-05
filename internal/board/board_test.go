package board

import (
	"testing"

	"pgregory.net/rapid"
)

// Feature: aws-learning-game, Property 4: Board Size Constraint
// For any created Board, the number of cells SHALL be at least 30.
// **Validates: Requirements 2.5**
func TestProperty4_BoardSizeConstraint(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random board configurations with size >= 30
		size := rapid.IntRange(30, 100).Draw(t, "size")

		config := BoardConfig{
			Size:                size,
			OpportunityCells:    rapid.IntRange(1, size/6).Draw(t, "opportunityCells"),
			FateCells:           rapid.IntRange(1, size/6).Draw(t, "fateCells"),
			ChallengeCells:      rapid.IntRange(1, size/6).Draw(t, "challengeCells"),
			BonusCells:          rapid.IntRange(1, size/6).Draw(t, "bonusCells"),
			BaseCapitalPerCell:  int64(rapid.IntRange(10, 100).Draw(t, "baseCapital")),
			BaseEmployeePerCell: rapid.IntRange(1, 5).Draw(t, "baseEmployee"),
		}

		manager := NewBoardManager()
		board, err := manager.CreateBoardWithConfig("test", config)

		if err != nil {
			t.Fatalf("CreateBoardWithConfig failed: %v", err)
		}

		// Verify board size is at least 30
		if board.Size < 30 {
			t.Errorf("Board size should be at least 30, got %d", board.Size)
		}

		// Verify cells count matches size
		if len(board.Cells) != board.Size {
			t.Errorf("Cells count mismatch: expected %d, got %d", board.Size, len(board.Cells))
		}

		// Verify first cell is start
		if board.Cells[0].Type != CellStart {
			t.Errorf("First cell should be start, got %s", board.Cells[0].Type)
		}
	})
}

// Test board creation with default config
func TestCreateBoard_Default(t *testing.T) {
	manager := NewBoardManager()
	board, err := manager.CreateBoard("default")

	if err != nil {
		t.Fatalf("CreateBoard failed: %v", err)
	}

	if board.Size < 30 {
		t.Errorf("Board size should be at least 30, got %d", board.Size)
	}

	if board.Cells[0].Type != CellStart {
		t.Errorf("First cell should be start, got %s", board.Cells[0].Type)
	}
}

// Test board creation with invalid size
func TestCreateBoard_InvalidSize(t *testing.T) {
	manager := NewBoardManager()

	config := BoardConfig{
		Size:                20, // Less than 30
		OpportunityCells:    2,
		FateCells:           2,
		ChallengeCells:      2,
		BonusCells:          1,
		BaseCapitalPerCell:  50,
		BaseEmployeePerCell: 1,
	}

	_, err := manager.CreateBoardWithConfig("test", config)
	if err != ErrInvalidBoardSize {
		t.Errorf("Expected ErrInvalidBoardSize, got %v", err)
	}
}

// Test GetCell
func TestGetCell(t *testing.T) {
	manager := NewBoardManager()
	board, _ := manager.CreateBoard("default")

	// Test valid position
	cell, err := manager.GetCell(0)
	if err != nil {
		t.Errorf("GetCell failed: %v", err)
	}
	if cell.Position != 0 {
		t.Errorf("Expected position 0, got %d", cell.Position)
	}

	// Test invalid position (negative)
	_, err = manager.GetCell(-1)
	if err != ErrInvalidPosition {
		t.Errorf("Expected ErrInvalidPosition, got %v", err)
	}

	// Test invalid position (out of bounds)
	_, err = manager.GetCell(board.Size + 1)
	if err != ErrInvalidPosition {
		t.Errorf("Expected ErrInvalidPosition, got %v", err)
	}
}

// Test CalculateNewPosition with wrap-around
func TestCalculateNewPosition(t *testing.T) {
	manager := NewBoardManager()
	manager.CreateBoard("default")
	board := manager.GetBoard()

	testCases := []struct {
		current   int
		diceValue int
		expected  int
	}{
		{0, 5, 5},
		{board.Size - 1, 1, 0},                         // Wrap around
		{board.Size - 3, 5, 2},                         // Wrap around
		{10, 0, 10},                                    // No movement
		{0, board.Size, 0},                             // Full circuit
		{5, board.Size + 3, (5 + board.Size + 3) % board.Size}, // More than full circuit
	}

	for _, tc := range testCases {
		result := manager.CalculateNewPosition(tc.current, tc.diceValue)
		if result != tc.expected {
			t.Errorf("CalculateNewPosition(%d, %d) = %d, expected %d",
				tc.current, tc.diceValue, result, tc.expected)
		}
	}
}

// Test GetCell without board
func TestGetCell_NoBoard(t *testing.T) {
	manager := NewBoardManager()

	_, err := manager.GetCell(0)
	if err != ErrBoardNotFound {
		t.Errorf("Expected ErrBoardNotFound, got %v", err)
	}
}
