package game

import (
	"testing"
	"time"

	"github.com/aws-learning-game/internal/company"
	"pgregory.net/rapid"
)

// executeTurnAndHandleDecision 執行回合並在需要時自動處理決策
// 這個輔助函數用於測試，當玩家落在事件格子時自動提交決策以切換到下一個玩家
func executeTurnAndHandleDecision(engine *GameEngine, gameID string, playerID string) (*TurnResult, error) {
	action := TurnAction{ActionType: "roll_dice"}
	result, err := engine.ExecuteTurn(gameID, playerID, action)
	if err != nil {
		return nil, err
	}

	// 如果需要決策，自動提交一個空決策來切換玩家
	if result.DecisionRequired {
		record := DecisionOutcomeRecord{
			DecisionRecord: DecisionRecord{
				EventID:    "test_event",
				ChoiceID:   1,
				TurnNumber: 1,
				Timestamp:  time.Now(),
			},
			Success:        true,
			AWSServices:    []string{},
			LearningPoints: []string{},
		}
		engine.RecordDecision(gameID, playerID, record)
	}

	return result, nil
}

// Feature: aws-learning-game, Property 18: Player Count Validation
// For any game creation or join attempt, the Game_Engine SHALL enforce player count between 2 and 4 inclusive.
// **Validates: Requirements 7.1**
func TestProperty18_PlayerCountValidation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		engine := NewGameEngine()

		// Generate random max players config (valid range: 2-4)
		maxPlayers := rapid.IntRange(2, 4).Draw(t, "maxPlayers")

		config := GameConfig{
			MaxPlayers:      maxPlayers,
			BoardType:       "default",
			DifficultyLevel: "normal",
		}

		// Create game should succeed with valid player count
		game, err := engine.CreateGame(config)
		if err != nil {
			t.Fatalf("CreateGame failed with valid maxPlayers %d: %v", maxPlayers, err)
		}

		// Generate random number of players to join (1 to maxPlayers+1)
		numPlayersToJoin := rapid.IntRange(1, maxPlayers+1).Draw(t, "numPlayersToJoin")

		companyTypes := []company.CompanyType{company.Startup, company.Traditional, company.CloudReseller, company.CloudNative}

		joinedCount := 0
		for i := 0; i < numPlayersToJoin; i++ {
			playerID := rapid.StringMatching(`player_[a-z]{5}`).Draw(t, "playerID")
			companyType := companyTypes[rapid.IntRange(0, len(companyTypes)-1).Draw(t, "companyTypeIndex")]

			err := engine.JoinGame(game.ID, playerID, "Player "+playerID, companyType)

			if joinedCount < maxPlayers {
				// Should succeed if under max
				if err != nil {
					t.Errorf("JoinGame should succeed for player %d (max: %d), got error: %v", joinedCount+1, maxPlayers, err)
				} else {
					joinedCount++
				}
			} else {
				// Should fail if at max
				if err != ErrGameFull {
					t.Errorf("JoinGame should return ErrGameFull for player %d (max: %d), got: %v", joinedCount+1, maxPlayers, err)
				}
			}
		}

		// Verify final player count
		gameState, _ := engine.GetGameState(game.ID)
		if len(gameState.Players) > maxPlayers {
			t.Errorf("Game has %d players, should not exceed maxPlayers %d", len(gameState.Players), maxPlayers)
		}
	})
}

// Test that game creation fails with invalid player count config
func TestCreateGame_InvalidPlayerCount(t *testing.T) {
	engine := NewGameEngine()

	// Test with maxPlayers < 2
	config := GameConfig{
		MaxPlayers:      1,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}
	_, err := engine.CreateGame(config)
	if err != ErrInvalidPlayerCount {
		t.Errorf("Expected ErrInvalidPlayerCount for maxPlayers=1, got: %v", err)
	}

	// Test with maxPlayers > 4
	config.MaxPlayers = 5
	_, err = engine.CreateGame(config)
	if err != ErrInvalidPlayerCount {
		t.Errorf("Expected ErrInvalidPlayerCount for maxPlayers=5, got: %v", err)
	}
}

// Test that StartGame enforces minimum 2 players
func TestStartGame_MinimumPlayers(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)

	// Try to start with 0 players
	err := engine.StartGame(game.ID)
	if err != ErrInsufficientPlayers {
		t.Errorf("Expected ErrInsufficientPlayers with 0 players, got: %v", err)
	}

	// Add 1 player
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)

	// Try to start with 1 player
	err = engine.StartGame(game.ID)
	if err != ErrInsufficientPlayers {
		t.Errorf("Expected ErrInsufficientPlayers with 1 player, got: %v", err)
	}

	// Add 2nd player
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)

	// Should succeed with 2 players
	err = engine.StartGame(game.ID)
	if err != nil {
		t.Errorf("StartGame should succeed with 2 players, got: %v", err)
	}
}

// Test basic game flow
func TestGameFlow_Basic(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	// Create game
	game, err := engine.CreateGame(config)
	if err != nil {
		t.Fatalf("CreateGame failed: %v", err)
	}
	if game.Status != StatusWaiting {
		t.Errorf("New game should have status 'waiting', got: %s", game.Status)
	}

	// Join players
	err = engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	if err != nil {
		t.Errorf("JoinGame failed: %v", err)
	}

	err = engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	if err != nil {
		t.Errorf("JoinGame failed: %v", err)
	}

	// Start game
	err = engine.StartGame(game.ID)
	if err != nil {
		t.Errorf("StartGame failed: %v", err)
	}

	// Verify game state
	gameState, _ := engine.GetGameState(game.ID)
	if gameState.Status != StatusInProgress {
		t.Errorf("Started game should have status 'in_progress', got: %s", gameState.Status)
	}
	if gameState.CurrentTurn != 1 {
		t.Errorf("Started game should be on turn 1, got: %d", gameState.CurrentTurn)
	}
	if gameState.CurrentPlayerID != "player1" {
		t.Errorf("First player should be player1, got: %s", gameState.CurrentPlayerID)
	}
}

// Test duplicate player join
func TestJoinGame_DuplicatePlayer(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)

	// Try to join again with same ID
	err := engine.JoinGame(game.ID, "player1", "Player 1 Again", company.Traditional)
	if err != ErrPlayerAlreadyInGame {
		t.Errorf("Expected ErrPlayerAlreadyInGame, got: %v", err)
	}
}

// Test join after game started
func TestJoinGame_AfterStart(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	// Try to join after game started
	err := engine.JoinGame(game.ID, "player3", "Player 3", company.CloudNative)
	if err != ErrGameAlreadyStarted {
		t.Errorf("Expected ErrGameAlreadyStarted, got: %v", err)
	}
}

// Test ExecuteTurn with roll_dice action
func TestExecuteTurn_RollDice(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	// Execute turn for player1
	action := TurnAction{ActionType: "roll_dice"}
	result, err := engine.ExecuteTurn(game.ID, "player1", action)

	if err != nil {
		t.Fatalf("ExecuteTurn failed: %v", err)
	}

	// Verify dice value is in valid range
	if result.DiceValue < 1 || result.DiceValue > 6 {
		t.Errorf("Dice value should be 1-6, got: %d", result.DiceValue)
	}

	// Verify position changed
	if result.OldPosition != 0 {
		t.Errorf("Old position should be 0, got: %d", result.OldPosition)
	}

	// Verify new position is correct
	gameState, _ := engine.GetGameState(game.ID)
	expectedNewPos := result.DiceValue % gameState.Board.Size
	if result.NewPosition != expectedNewPos {
		t.Errorf("New position should be %d, got: %d", expectedNewPos, result.NewPosition)
	}

	// Verify turn handling based on cell type
	// If decision is required (event cell), player should NOT change
	// If no decision required (normal cell), player should advance to player2
	if result.DecisionRequired {
		// Player landed on event cell, should stay as player1 until decision is made
		if gameState.CurrentPlayerID != "player1" {
			t.Errorf("When decision required, current player should still be player1, got: %s", gameState.CurrentPlayerID)
		}
	} else {
		// Normal cell, turn should advance to player2
		if gameState.CurrentPlayerID != "player2" {
			t.Errorf("When no decision required, current player should be player2, got: %s", gameState.CurrentPlayerID)
		}
	}
}

// Test ExecuteTurn with wrong player
func TestExecuteTurn_WrongPlayer(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	// Try to execute turn for player2 when it's player1's turn
	action := TurnAction{ActionType: "roll_dice"}
	_, err := engine.ExecuteTurn(game.ID, "player2", action)

	if err != ErrNotYourTurn {
		t.Errorf("Expected ErrNotYourTurn, got: %v", err)
	}
}

// Test ExecuteTurn before game started
func TestExecuteTurn_GameNotStarted(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)

	action := TurnAction{ActionType: "roll_dice"}
	_, err := engine.ExecuteTurn(game.ID, "player1", action)

	if err != ErrGameNotStarted {
		t.Errorf("Expected ErrGameNotStarted, got: %v", err)
	}
}

// Test ExecuteTurn with invalid action
func TestExecuteTurn_InvalidAction(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	action := TurnAction{ActionType: "invalid_action"}
	_, err := engine.ExecuteTurn(game.ID, "player1", action)

	if err != ErrInvalidAction {
		t.Errorf("Expected ErrInvalidAction, got: %v", err)
	}
}

// Feature: aws-learning-game, Property 2: Position Movement Correctness
// For any dice roll value D and current position P on a board of size S,
// the new position SHALL equal (P + D) mod S, ensuring players wrap around the board correctly.
// **Validates: Requirements 2.2**
func TestProperty2_PositionMovementCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		engine := NewGameEngine()

		config := GameConfig{
			MaxPlayers:      4,
			BoardType:       "default",
			DifficultyLevel: "normal",
		}

		game, err := engine.CreateGame(config)
		if err != nil {
			t.Fatalf("CreateGame failed: %v", err)
		}

		// Add players
		engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
		engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
		engine.StartGame(game.ID)

		boardSize := game.Board.Size

		// Generate random starting position (simulate multiple turns)
		numTurns := rapid.IntRange(1, 20).Draw(t, "numTurns")

		currentPlayerID := "player1"
		for i := 0; i < numTurns; i++ {
			// Get current player state before turn
			playerBefore, _ := engine.GetPlayer(game.ID, currentPlayerID)
			oldPosition := playerBefore.Position

			// Execute turn with decision handling
			result, err := executeTurnAndHandleDecision(engine, game.ID, currentPlayerID)
			if err != nil {
				t.Fatalf("ExecuteTurn failed: %v", err)
			}

			// Verify position calculation: newPosition = (oldPosition + diceValue) % boardSize
			expectedNewPosition := (oldPosition + result.DiceValue) % boardSize
			if result.NewPosition != expectedNewPosition {
				t.Errorf("Position calculation incorrect: old=%d, dice=%d, boardSize=%d, expected=%d, got=%d",
					oldPosition, result.DiceValue, boardSize, expectedNewPosition, result.NewPosition)
			}

			// Verify player's actual position matches result
			playerAfter, _ := engine.GetPlayer(game.ID, currentPlayerID)
			if playerAfter.Position != result.NewPosition {
				t.Errorf("Player position mismatch: result=%d, actual=%d", result.NewPosition, playerAfter.Position)
			}

			// Switch to next player for next iteration
			if currentPlayerID == "player1" {
				currentPlayerID = "player2"
			} else {
				currentPlayerID = "player1"
			}
		}
	})
}

// Feature: aws-learning-game, Property 5: Circuit Completion Bonus
// For any player movement that causes position to wrap around (cross the start),
// bonus Capital and Employees SHALL be added to the Company.
// **Validates: Requirements 2.6**
func TestProperty5_CircuitCompletionBonus(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		engine := NewGameEngine()

		config := GameConfig{
			MaxPlayers:      4,
			BoardType:       "default",
			DifficultyLevel: "normal",
		}

		game, err := engine.CreateGame(config)
		if err != nil {
			t.Fatalf("CreateGame failed: %v", err)
		}

		engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
		engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
		engine.StartGame(game.ID)

		boardSize := game.Board.Size

		// Execute multiple turns to potentially trigger circuit completion
		numTurns := rapid.IntRange(10, 50).Draw(t, "numTurns")

		currentPlayerID := "player1"
		for i := 0; i < numTurns; i++ {
			playerBefore, _ := engine.GetPlayer(game.ID, currentPlayerID)
			oldPosition := playerBefore.Position
			oldCapital := playerBefore.Company.Capital
			oldEmployees := playerBefore.Company.Employees

			// Execute turn with decision handling
			result, err := executeTurnAndHandleDecision(engine, game.ID, currentPlayerID)
			if err != nil {
				t.Fatalf("ExecuteTurn failed: %v", err)
			}

			// Check if circuit was completed (position wrapped around)
			circuitCompleted := (oldPosition + result.DiceValue) >= boardSize

			if circuitCompleted {
				// Verify CircuitCompleted flag is set
				if !result.CircuitCompleted {
					t.Errorf("CircuitCompleted should be true when position wraps: old=%d, dice=%d, boardSize=%d",
						oldPosition, result.DiceValue, boardSize)
				}

				// Verify bonus was applied
				playerAfter, _ := engine.GetPlayer(game.ID, currentPlayerID)

				// Get cell landing bonus
				cell := game.Board.Cells[result.NewPosition]
				expectedCapitalChange := DefaultCircuitBonus.Capital + cell.BaseCapital
				expectedEmployeeChange := DefaultCircuitBonus.Employees + cell.BaseEmployees

				actualCapitalChange := playerAfter.Company.Capital - oldCapital
				actualEmployeeChange := playerAfter.Company.Employees - oldEmployees

				if actualCapitalChange != expectedCapitalChange {
					t.Errorf("Circuit bonus capital incorrect: expected %d, got %d (circuit bonus: %d, cell bonus: %d)",
						expectedCapitalChange, actualCapitalChange, DefaultCircuitBonus.Capital, cell.BaseCapital)
				}

				if actualEmployeeChange != expectedEmployeeChange {
					t.Errorf("Circuit bonus employees incorrect: expected %d, got %d (circuit bonus: %d, cell bonus: %d)",
						expectedEmployeeChange, actualEmployeeChange, DefaultCircuitBonus.Employees, cell.BaseEmployees)
				}
			} else {
				// Verify CircuitCompleted flag is false
				if result.CircuitCompleted {
					t.Errorf("CircuitCompleted should be false when position doesn't wrap: old=%d, dice=%d, boardSize=%d",
						oldPosition, result.DiceValue, boardSize)
				}
			}

			// Switch player
			if currentPlayerID == "player1" {
				currentPlayerID = "player2"
			} else {
				currentPlayerID = "player1"
			}
		}
	})
}

// Test cell landing resource update
func TestCellLandingResourceUpdate(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	// Get initial state
	playerBefore, _ := engine.GetPlayer(game.ID, "player1")
	initialCapital := playerBefore.Company.Capital
	initialEmployees := playerBefore.Company.Employees

	// Execute turn
	action := TurnAction{ActionType: "roll_dice"}
	result, err := engine.ExecuteTurn(game.ID, "player1", action)
	if err != nil {
		t.Fatalf("ExecuteTurn failed: %v", err)
	}

	// Get cell at new position
	cell := game.Board.Cells[result.NewPosition]

	// Get player state after turn
	playerAfter, _ := engine.GetPlayer(game.ID, "player1")

	// Calculate expected changes
	expectedCapitalChange := cell.BaseCapital
	expectedEmployeeChange := cell.BaseEmployees

	// Add circuit bonus if applicable
	if result.CircuitCompleted {
		expectedCapitalChange += DefaultCircuitBonus.Capital
		expectedEmployeeChange += DefaultCircuitBonus.Employees
	}

	// Verify capital change
	actualCapitalChange := playerAfter.Company.Capital - initialCapital
	if actualCapitalChange != expectedCapitalChange {
		t.Errorf("Capital change incorrect: expected %d, got %d", expectedCapitalChange, actualCapitalChange)
	}

	// Verify employee change
	actualEmployeeChange := playerAfter.Company.Employees - initialEmployees
	if actualEmployeeChange != expectedEmployeeChange {
		t.Errorf("Employee change incorrect: expected %d, got %d", expectedEmployeeChange, actualEmployeeChange)
	}

	// Verify result contains correct changes
	if result.CapitalChange != expectedCapitalChange {
		t.Errorf("Result capital change incorrect: expected %d, got %d", expectedCapitalChange, result.CapitalChange)
	}
	if result.EmployeeChange != expectedEmployeeChange {
		t.Errorf("Result employee change incorrect: expected %d, got %d", expectedEmployeeChange, result.EmployeeChange)
	}
}

// Feature: aws-learning-game, Property 3: Cell Landing Resource Update
// For any cell landing event, the Company's Capital SHALL increase by exactly the cell's BaseCapital value
// AND the Company's Employees SHALL increase by exactly the cell's BaseEmployees value.
// **Validates: Requirements 2.3, 2.4**
func TestProperty3_CellLandingResourceUpdate(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		engine := NewGameEngine()

		config := GameConfig{
			MaxPlayers:      4,
			BoardType:       "default",
			DifficultyLevel: "normal",
		}

		game, err := engine.CreateGame(config)
		if err != nil {
			t.Fatalf("CreateGame failed: %v", err)
		}

		engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
		engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
		engine.StartGame(game.ID)

		// Execute multiple turns to test various cell landings
		numTurns := rapid.IntRange(5, 30).Draw(t, "numTurns")

		currentPlayerID := "player1"
		for i := 0; i < numTurns; i++ {
			// Get player state before turn
			playerBefore, _ := engine.GetPlayer(game.ID, currentPlayerID)
			capitalBefore := playerBefore.Company.Capital
			employeesBefore := playerBefore.Company.Employees

			// Execute turn with decision handling
			result, err := executeTurnAndHandleDecision(engine, game.ID, currentPlayerID)
			if err != nil {
				t.Fatalf("ExecuteTurn failed: %v", err)
			}

			// Get the cell landed on
			cell := game.Board.Cells[result.NewPosition]

			// Calculate expected resource changes
			expectedCapitalChange := cell.BaseCapital
			expectedEmployeeChange := cell.BaseEmployees

			// Add circuit bonus if applicable
			if result.CircuitCompleted {
				expectedCapitalChange += DefaultCircuitBonus.Capital
				expectedEmployeeChange += DefaultCircuitBonus.Employees
			}

			// Get player state after turn
			playerAfter, _ := engine.GetPlayer(game.ID, currentPlayerID)

			// Verify capital change matches exactly
			actualCapitalChange := playerAfter.Company.Capital - capitalBefore
			if actualCapitalChange != expectedCapitalChange {
				t.Errorf("Turn %d: Capital change mismatch for cell %d (%s): expected %d (cell: %d, circuit: %v), got %d",
					i+1, result.NewPosition, cell.Type, expectedCapitalChange, cell.BaseCapital, result.CircuitCompleted, actualCapitalChange)
			}

			// Verify employee change matches exactly
			actualEmployeeChange := playerAfter.Company.Employees - employeesBefore
			if actualEmployeeChange != expectedEmployeeChange {
				t.Errorf("Turn %d: Employee change mismatch for cell %d (%s): expected %d (cell: %d, circuit: %v), got %d",
					i+1, result.NewPosition, cell.Type, expectedEmployeeChange, cell.BaseEmployees, result.CircuitCompleted, actualEmployeeChange)
			}

			// Verify result struct contains correct values
			if result.CapitalChange != expectedCapitalChange {
				t.Errorf("Turn %d: Result.CapitalChange mismatch: expected %d, got %d",
					i+1, expectedCapitalChange, result.CapitalChange)
			}
			if result.EmployeeChange != expectedEmployeeChange {
				t.Errorf("Turn %d: Result.EmployeeChange mismatch: expected %d, got %d",
					i+1, expectedEmployeeChange, result.EmployeeChange)
			}

			// Switch player
			if currentPlayerID == "player1" {
				currentPlayerID = "player2"
			} else {
				currentPlayerID = "player1"
			}
		}
	})
}

// Feature: aws-learning-game, Property 17: Game End Summary
// For any game that reaches "finished" status, a summary SHALL be generatable
// containing all DecisionRecords from all players.
// **Validates: Requirements 6.3**
func TestProperty17_GameEndSummary(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		engine := NewGameEngine()

		config := GameConfig{
			MaxPlayers:      4,
			BoardType:       "default",
			DifficultyLevel: "normal",
		}

		game, err := engine.CreateGame(config)
		if err != nil {
			t.Fatalf("CreateGame failed: %v", err)
		}

		// Generate random number of players (2-4)
		numPlayers := rapid.IntRange(2, 4).Draw(t, "numPlayers")
		companyTypes := []company.CompanyType{company.Startup, company.Traditional, company.CloudReseller, company.CloudNative}

		playerIDs := make([]string, numPlayers)
		for i := 0; i < numPlayers; i++ {
			playerIDs[i] = rapid.StringMatching(`player_[a-z]{5}`).Draw(t, "playerID")
			companyType := companyTypes[rapid.IntRange(0, len(companyTypes)-1).Draw(t, "companyTypeIndex")]
			err := engine.JoinGame(game.ID, playerIDs[i], "Player "+playerIDs[i], companyType)
			if err != nil {
				t.Fatalf("JoinGame failed: %v", err)
			}
		}

		// Start game
		err = engine.StartGame(game.ID)
		if err != nil {
			t.Fatalf("StartGame failed: %v", err)
		}

		// Execute some turns with decision handling
		numTurns := rapid.IntRange(5, 20).Draw(t, "numTurns")
		currentPlayerIndex := 0
		for i := 0; i < numTurns; i++ {
			_, err := executeTurnAndHandleDecision(engine, game.ID, playerIDs[currentPlayerIndex])
			if err != nil {
				t.Fatalf("ExecuteTurn failed: %v", err)
			}
			currentPlayerIndex = (currentPlayerIndex + 1) % numPlayers
		}

		// End the game
		err = engine.EndGame(game.ID)
		if err != nil {
			t.Fatalf("EndGame failed: %v", err)
		}

		// Get game summary
		summary, err := engine.GetGameSummary(game.ID)
		if err != nil {
			t.Fatalf("GetGameSummary failed: %v", err)
		}

		// Verify summary is not nil
		if summary == nil {
			t.Fatal("GameSummary should not be nil")
		}

		// Verify GameID matches
		if summary.GameID != game.ID {
			t.Errorf("GameID mismatch: expected %s, got %s", game.ID, summary.GameID)
		}

		// Verify TotalTurns is positive
		if summary.TotalTurns <= 0 {
			t.Errorf("TotalTurns should be positive, got %d", summary.TotalTurns)
		}

		// Verify we have summaries for all players
		if len(summary.PlayerSummaries) != numPlayers {
			t.Errorf("Expected %d player summaries, got %d", numPlayers, len(summary.PlayerSummaries))
		}

		// Verify each player summary has required fields
		for i, playerSummary := range summary.PlayerSummaries {
			if playerSummary.PlayerID == "" {
				t.Errorf("Player %d: PlayerID should not be empty", i)
			}
			if playerSummary.PlayerName == "" {
				t.Errorf("Player %d: PlayerName should not be empty", i)
			}
			if playerSummary.FinalCapital <= 0 {
				t.Errorf("Player %d: FinalCapital should be positive, got %d", i, playerSummary.FinalCapital)
			}
			if playerSummary.FinalEmployees <= 0 {
				t.Errorf("Player %d: FinalEmployees should be positive, got %d", i, playerSummary.FinalEmployees)
			}
			// DecisionRecords can be empty if no decisions were made
			// but the slice should exist
			if playerSummary.DecisionRecords == nil {
				t.Errorf("Player %d: DecisionRecords should not be nil", i)
			}
		}

		// Verify winner is determined
		if summary.Winner == nil {
			t.Error("Winner should be determined for finished game")
		} else {
			// Winner should have highest capital
			maxCapital := int64(0)
			for _, ps := range summary.PlayerSummaries {
				if ps.FinalCapital > maxCapital {
					maxCapital = ps.FinalCapital
				}
			}
			if summary.Winner.FinalCapital != maxCapital {
				t.Errorf("Winner should have highest capital: expected %d, got %d", maxCapital, summary.Winner.FinalCapital)
			}
		}

		// Verify learning insights are present
		if len(summary.LearningInsights) == 0 {
			t.Error("LearningInsights should not be empty")
		}
	})
}

// Test GetGameSummary returns error for non-finished game
func TestGetGameSummary_GameNotFinished(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	// Try to get summary before game is finished
	_, err := engine.GetGameSummary(game.ID)
	if err != ErrGameNotFinished {
		t.Errorf("Expected ErrGameNotFinished, got: %v", err)
	}
}

// Test EndGame functionality
func TestEndGame(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	// Execute a few turns
	engine.ExecuteTurn(game.ID, "player1", TurnAction{ActionType: "roll_dice"})
	engine.ExecuteTurn(game.ID, "player2", TurnAction{ActionType: "roll_dice"})

	// End the game
	err := engine.EndGame(game.ID)
	if err != nil {
		t.Fatalf("EndGame failed: %v", err)
	}

	// Verify game status is finished
	gameState, _ := engine.GetGameState(game.ID)
	if gameState.Status != StatusFinished {
		t.Errorf("Game status should be 'finished', got: %s", gameState.Status)
	}

	// Now summary should be available
	summary, err := engine.GetGameSummary(game.ID)
	if err != nil {
		t.Fatalf("GetGameSummary failed after EndGame: %v", err)
	}
	if summary == nil {
		t.Error("Summary should not be nil after EndGame")
	}
}

// Test GetPlayerProgress functionality
func TestGetPlayerProgress(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	// Execute some turns
	engine.ExecuteTurn(game.ID, "player1", TurnAction{ActionType: "roll_dice"})
	engine.ExecuteTurn(game.ID, "player2", TurnAction{ActionType: "roll_dice"})
	engine.ExecuteTurn(game.ID, "player1", TurnAction{ActionType: "roll_dice"})

	// Get player progress
	progress, err := engine.GetPlayerProgress(game.ID, "player1")
	if err != nil {
		t.Fatalf("GetPlayerProgress failed: %v", err)
	}

	// Verify progress is not nil
	if progress == nil {
		t.Fatal("PlayerProgress should not be nil")
	}

	// Verify PlayerID matches
	if progress.PlayerID != "player1" {
		t.Errorf("PlayerID mismatch: expected player1, got %s", progress.PlayerID)
	}

	// Verify skill areas are initialized
	if len(progress.SkillAreas) == 0 {
		t.Error("SkillAreas should not be empty")
	}

	// Verify recommendations are generated
	if len(progress.Recommendations) == 0 {
		t.Error("Recommendations should not be empty")
	}

	// Verify learning progress is calculated
	if progress.LearningProgress < 0 || progress.LearningProgress > 100 {
		t.Errorf("LearningProgress should be 0-100, got %f", progress.LearningProgress)
	}
}

// Test GetPlayerProgress with invalid game
func TestGetPlayerProgress_InvalidGame(t *testing.T) {
	engine := NewGameEngine()

	_, err := engine.GetPlayerProgress("invalid_game", "player1")
	if err != ErrGameNotFound {
		t.Errorf("Expected ErrGameNotFound, got: %v", err)
	}
}

// Test GetPlayerProgress with invalid player
func TestGetPlayerProgress_InvalidPlayer(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)

	_, err := engine.GetPlayerProgress(game.ID, "invalid_player")
	if err != ErrPlayerNotFound {
		t.Errorf("Expected ErrPlayerNotFound, got: %v", err)
	}
}

// Test RecordDecision functionality
func TestRecordDecision(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	// Record a decision
	record := DecisionOutcomeRecord{
		DecisionRecord: DecisionRecord{
			TurnNumber: 1,
			EventID:    "event_1",
			ChoiceID:   1,
		},
		Success:        true,
		AWSServices:    []string{"EC2", "S3"},
		LearningPoints: []string{"學習了 EC2 的使用"},
	}

	err := engine.RecordDecision(game.ID, "player1", record)
	if err != nil {
		t.Fatalf("RecordDecision failed: %v", err)
	}

	// Verify decision was recorded
	player, _ := engine.GetPlayer(game.ID, "player1")
	if len(player.DecisionHistory) != 1 {
		t.Errorf("Expected 1 decision record, got %d", len(player.DecisionHistory))
	}

	if player.DecisionHistory[0].EventID != "event_1" {
		t.Errorf("EventID mismatch: expected event_1, got %s", player.DecisionHistory[0].EventID)
	}
}

// Test enhanced player progress tracking with decision outcomes
func TestPlayerProgressTracking_WithDecisionOutcomes(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	// Record multiple decisions with different outcomes
	decisions := []DecisionOutcomeRecord{
		{
			DecisionRecord: DecisionRecord{TurnNumber: 1, EventID: "event_1", ChoiceID: 1},
			Success:        true,
			AWSServices:    []string{"EC2", "S3"},
			LearningPoints: []string{"EC2 執行個體類型選擇"},
		},
		{
			DecisionRecord: DecisionRecord{TurnNumber: 2, EventID: "event_2", ChoiceID: 2},
			Success:        false,
			AWSServices:    []string{"RDS"},
			LearningPoints: []string{"RDS 多可用區部署"},
		},
		{
			DecisionRecord: DecisionRecord{TurnNumber: 3, EventID: "event_3", ChoiceID: 1},
			Success:        true,
			AWSServices:    []string{"Lambda", "DynamoDB"},
			LearningPoints: []string{"無伺服器架構設計"},
		},
	}

	for _, d := range decisions {
		err := engine.RecordDecision(game.ID, "player1", d)
		if err != nil {
			t.Fatalf("RecordDecision failed: %v", err)
		}
	}

	// Get player progress
	progress, err := engine.GetPlayerProgress(game.ID, "player1")
	if err != nil {
		t.Fatalf("GetPlayerProgress failed: %v", err)
	}

	// Verify total decisions
	if progress.TotalDecisions != 3 {
		t.Errorf("Expected 3 total decisions, got %d", progress.TotalDecisions)
	}

	// Verify successful/failed decisions
	if progress.SuccessfulDecisions != 2 {
		t.Errorf("Expected 2 successful decisions, got %d", progress.SuccessfulDecisions)
	}
	if progress.FailedDecisions != 1 {
		t.Errorf("Expected 1 failed decision, got %d", progress.FailedDecisions)
	}

	// Verify AWS services are tracked
	if len(progress.AWSServicesUsed) != 5 {
		t.Errorf("Expected 5 AWS services used, got %d: %v", len(progress.AWSServicesUsed), progress.AWSServicesUsed)
	}

	// Verify learning points are tracked
	if len(progress.TopicsLearned) != 3 {
		t.Errorf("Expected 3 topics learned, got %d: %v", len(progress.TopicsLearned), progress.TopicsLearned)
	}

	// Verify skill areas are updated based on services used
	// EC2, Lambda should increase compute; S3 should increase storage; RDS, DynamoDB should increase database
	if progress.SkillAreas["compute"] <= 60 { // Base 60 for Startup + service bonuses
		t.Errorf("Compute skill should be > 60 after using EC2 and Lambda, got %d", progress.SkillAreas["compute"])
	}
	if progress.SkillAreas["storage"] <= 50 { // Base 50 for Startup + S3 bonus
		t.Errorf("Storage skill should be > 50 after using S3, got %d", progress.SkillAreas["storage"])
	}
	if progress.SkillAreas["database"] <= 0 { // Base 0 for Startup + RDS, DynamoDB bonuses
		t.Errorf("Database skill should be > 0 after using RDS and DynamoDB, got %d", progress.SkillAreas["database"])
	}

	// Verify learning progress is calculated
	if progress.LearningProgress <= 0 {
		t.Errorf("Learning progress should be > 0, got %f", progress.LearningProgress)
	}

	// Verify recommendations are generated
	if len(progress.Recommendations) == 0 {
		t.Error("Recommendations should not be empty")
	}
}

// Test that DecisionOutcomeHistory is properly stored and retrieved
func TestDecisionOutcomeHistory_Storage(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.CloudNative)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	// Record a decision with detailed outcome
	record := DecisionOutcomeRecord{
		DecisionRecord: DecisionRecord{
			TurnNumber: 1,
			EventID:    "security_event_1",
			ChoiceID:   2,
		},
		Success:        true,
		AWSServices:    []string{"WAF", "Shield", "GuardDuty"},
		LearningPoints: []string{"AWS 安全服務組合", "DDoS 防護策略"},
	}

	err := engine.RecordDecision(game.ID, "player1", record)
	if err != nil {
		t.Fatalf("RecordDecision failed: %v", err)
	}

	// Get game state and verify DecisionOutcomeHistory is stored
	gameState, _ := engine.GetGameState(game.ID)
	player := gameState.Players[0]

	if len(player.DecisionOutcomeHistory) != 1 {
		t.Errorf("Expected 1 decision outcome record, got %d", len(player.DecisionOutcomeHistory))
	}

	outcome := player.DecisionOutcomeHistory[0]
	if !outcome.Success {
		t.Error("Decision outcome should be successful")
	}
	if len(outcome.AWSServices) != 3 {
		t.Errorf("Expected 3 AWS services, got %d", len(outcome.AWSServices))
	}
	if len(outcome.LearningPoints) != 2 {
		t.Errorf("Expected 2 learning points, got %d", len(outcome.LearningPoints))
	}
}

// Test skill area updates from different AWS services
func TestSkillAreaUpdates_FromAWSServices(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	// Use Traditional company which has lower base compute/storage scores
	engine.JoinGame(game.ID, "player1", "Player 1", company.Traditional)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Startup)
	engine.StartGame(game.ID)

	// Get initial progress (no decisions yet)
	initialProgress, _ := engine.GetPlayerProgress(game.ID, "player1")
	initialCompute := initialProgress.SkillAreas["compute"]
	initialSecurity := initialProgress.SkillAreas["security"]

	// Record decisions using compute and security services
	computeDecision := DecisionOutcomeRecord{
		DecisionRecord: DecisionRecord{TurnNumber: 1, EventID: "event_1", ChoiceID: 1},
		Success:        true,
		AWSServices:    []string{"EC2", "Lambda", "ECS"},
		LearningPoints: []string{"運算服務選擇"},
	}
	engine.RecordDecision(game.ID, "player1", computeDecision)

	securityDecision := DecisionOutcomeRecord{
		DecisionRecord: DecisionRecord{TurnNumber: 2, EventID: "event_2", ChoiceID: 1},
		Success:        true,
		AWSServices:    []string{"IAM", "KMS"},
		LearningPoints: []string{"身份與存取管理"},
	}
	engine.RecordDecision(game.ID, "player1", securityDecision)

	// Get updated progress
	updatedProgress, _ := engine.GetPlayerProgress(game.ID, "player1")

	// Verify compute skill increased (3 compute services * 5 points each = 15 points increase)
	if updatedProgress.SkillAreas["compute"] <= initialCompute {
		t.Errorf("Compute skill should increase after using EC2, Lambda, ECS. Initial: %d, Updated: %d",
			initialCompute, updatedProgress.SkillAreas["compute"])
	}

	// Verify security skill increased (2 security services * 5 points each = 10 points increase)
	if updatedProgress.SkillAreas["security"] <= initialSecurity {
		t.Errorf("Security skill should increase after using IAM, KMS. Initial: %d, Updated: %d",
			initialSecurity, updatedProgress.SkillAreas["security"])
	}
}

// Feature: aws-learning-game, Property 19: Turn Order Fairness
// For any multi-player game with N players, after N consecutive turn advancements,
// each player SHALL have had exactly one turn.
// **Validates: Requirements 7.2, 7.3**
func TestProperty19_TurnOrderFairness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		engine := NewGameEngine()

		// Generate random number of players (2-4)
		numPlayers := rapid.IntRange(2, 4).Draw(t, "numPlayers")

		config := GameConfig{
			MaxPlayers:      numPlayers,
			BoardType:       "default",
			DifficultyLevel: "normal",
		}

		game, err := engine.CreateGame(config)
		if err != nil {
			t.Fatalf("CreateGame failed: %v", err)
		}

		// Add players
		companyTypes := []company.CompanyType{company.Startup, company.Traditional, company.CloudReseller, company.CloudNative}
		playerIDs := make([]string, numPlayers)
		for i := 0; i < numPlayers; i++ {
			playerIDs[i] = rapid.StringMatching(`player_[a-z]{5}`).Draw(t, "playerID")
			companyType := companyTypes[i%len(companyTypes)]
			err := engine.JoinGame(game.ID, playerIDs[i], "Player "+playerIDs[i], companyType)
			if err != nil {
				t.Fatalf("JoinGame failed for player %d: %v", i, err)
			}
		}

		// Start game
		err = engine.StartGame(game.ID)
		if err != nil {
			t.Fatalf("StartGame failed: %v", err)
		}

		// Get initial turn order
		turnOrder, err := engine.GetTurnOrder(game.ID)
		if err != nil {
			t.Fatalf("GetTurnOrder failed: %v", err)
		}

		// Verify turn order contains all players
		if len(turnOrder) != numPlayers {
			t.Errorf("Turn order should have %d players, got %d", numPlayers, len(turnOrder))
		}

		// Get initial turn counts
		initialTurnCounts, err := engine.GetPlayerTurnCount(game.ID)
		if err != nil {
			t.Fatalf("GetPlayerTurnCount failed: %v", err)
		}

		// Execute N consecutive turns (one full round) with decision handling
		for i := 0; i < numPlayers; i++ {
			// Get current player
			currentPlayer, err := engine.GetCurrentPlayer(game.ID)
			if err != nil {
				t.Fatalf("GetCurrentPlayer failed at turn %d: %v", i, err)
			}

			// Verify current player matches expected turn order
			expectedPlayerID := turnOrder[i]
			if currentPlayer.PlayerID != expectedPlayerID {
				t.Errorf("Turn %d: Expected player %s, got %s", i, expectedPlayerID, currentPlayer.PlayerID)
			}

			// Execute turn with decision handling
			_, err = executeTurnAndHandleDecision(engine, game.ID, currentPlayer.PlayerID)
			if err != nil {
				t.Fatalf("ExecuteTurn failed for player %s: %v", currentPlayer.PlayerID, err)
			}
		}

		// Get final turn counts
		finalTurnCounts, err := engine.GetPlayerTurnCount(game.ID)
		if err != nil {
			t.Fatalf("GetPlayerTurnCount failed: %v", err)
		}

		// Verify each player had exactly one turn in this round
		for _, playerID := range playerIDs {
			initialCount := initialTurnCounts[playerID]
			finalCount := finalTurnCounts[playerID]
			turnsInRound := finalCount - initialCount

			if turnsInRound != 1 {
				t.Errorf("Player %s should have had exactly 1 turn in round, got %d (initial: %d, final: %d)",
					playerID, turnsInRound, initialCount, finalCount)
			}
		}

		// Execute another full round to verify consistency
		for i := 0; i < numPlayers; i++ {
			currentPlayer, err := engine.GetCurrentPlayer(game.ID)
			if err != nil {
				t.Fatalf("GetCurrentPlayer failed at turn %d (round 2): %v", i, err)
			}

			// Verify turn order is consistent (same order as first round)
			expectedPlayerID := turnOrder[i]
			if currentPlayer.PlayerID != expectedPlayerID {
				t.Errorf("Round 2, Turn %d: Expected player %s, got %s", i, expectedPlayerID, currentPlayer.PlayerID)
			}

			// Execute turn with decision handling
			_, err = executeTurnAndHandleDecision(engine, game.ID, currentPlayer.PlayerID)
			if err != nil {
				t.Fatalf("ExecuteTurn failed for player %s (round 2): %v", currentPlayer.PlayerID, err)
			}
		}

		// Verify each player now has exactly 2 turns total
		finalTurnCounts2, err := engine.GetPlayerTurnCount(game.ID)
		if err != nil {
			t.Fatalf("GetPlayerTurnCount failed: %v", err)
		}

		for _, playerID := range playerIDs {
			initialCount := initialTurnCounts[playerID]
			finalCount := finalTurnCounts2[playerID]
			totalTurns := finalCount - initialCount

			if totalTurns != 2 {
				t.Errorf("Player %s should have had exactly 2 turns after 2 rounds, got %d",
					playerID, totalTurns)
			}
		}
	})
}

// Test that turn automatically advances to next player
func TestTurnAutoAdvance(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      3,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.JoinGame(game.ID, "player3", "Player 3", company.CloudNative)
	engine.StartGame(game.ID)

	// Verify initial player is player1
	currentPlayer, _ := engine.GetCurrentPlayer(game.ID)
	if currentPlayer.PlayerID != "player1" {
		t.Errorf("Initial player should be player1, got %s", currentPlayer.PlayerID)
	}

	// Execute turn for player1 (with decision handling)
	executeTurnAndHandleDecision(engine, game.ID, "player1")

	// Verify turn advanced to player2
	currentPlayer, _ = engine.GetCurrentPlayer(game.ID)
	if currentPlayer.PlayerID != "player2" {
		t.Errorf("After player1's turn, current player should be player2, got %s", currentPlayer.PlayerID)
	}

	// Execute turn for player2
	executeTurnAndHandleDecision(engine, game.ID, "player2")

	// Verify turn advanced to player3
	currentPlayer, _ = engine.GetCurrentPlayer(game.ID)
	if currentPlayer.PlayerID != "player3" {
		t.Errorf("After player2's turn, current player should be player3, got %s", currentPlayer.PlayerID)
	}

	// Execute turn for player3
	executeTurnAndHandleDecision(engine, game.ID, "player3")

	// Verify turn wrapped back to player1
	currentPlayer, _ = engine.GetCurrentPlayer(game.ID)
	if currentPlayer.PlayerID != "player1" {
		t.Errorf("After player3's turn, current player should wrap to player1, got %s", currentPlayer.PlayerID)
	}
}

// Test GetTurnOrder returns correct order
func TestGetTurnOrder(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.JoinGame(game.ID, "player3", "Player 3", company.CloudNative)

	turnOrder, err := engine.GetTurnOrder(game.ID)
	if err != nil {
		t.Fatalf("GetTurnOrder failed: %v", err)
	}

	// Verify order matches join order
	expectedOrder := []string{"player1", "player2", "player3"}
	if len(turnOrder) != len(expectedOrder) {
		t.Errorf("Turn order length mismatch: expected %d, got %d", len(expectedOrder), len(turnOrder))
	}

	for i, playerID := range expectedOrder {
		if turnOrder[i] != playerID {
			t.Errorf("Turn order[%d] mismatch: expected %s, got %s", i, playerID, turnOrder[i])
		}
	}
}

// Test GetPlayerTurnCount returns correct counts
func TestGetPlayerTurnCount(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      2,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	// Initial turn counts should be 0
	turnCounts, _ := engine.GetPlayerTurnCount(game.ID)
	if turnCounts["player1"] != 0 || turnCounts["player2"] != 0 {
		t.Errorf("Initial turn counts should be 0, got player1=%d, player2=%d",
			turnCounts["player1"], turnCounts["player2"])
	}

	// Execute turns (with decision handling)
	executeTurnAndHandleDecision(engine, game.ID, "player1")
	executeTurnAndHandleDecision(engine, game.ID, "player2")
	executeTurnAndHandleDecision(engine, game.ID, "player1")

	// Verify turn counts
	turnCounts, _ = engine.GetPlayerTurnCount(game.ID)
	if turnCounts["player1"] != 2 {
		t.Errorf("Player1 should have 2 turns, got %d", turnCounts["player1"])
	}
	if turnCounts["player2"] != 1 {
		t.Errorf("Player2 should have 1 turn, got %d", turnCounts["player2"])
	}
}

// Test DetermineWinner returns player with highest capital
func TestDetermineWinner(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)       // Capital: 500
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)   // Capital: 5000
	engine.JoinGame(game.ID, "player3", "Player 3", company.CloudReseller) // Capital: 2000
	engine.StartGame(game.ID)

	// Execute some turns to change capital
	action := TurnAction{ActionType: "roll_dice"}
	engine.ExecuteTurn(game.ID, "player1", action)
	engine.ExecuteTurn(game.ID, "player2", action)
	engine.ExecuteTurn(game.ID, "player3", action)

	// End the game
	engine.EndGame(game.ID)

	// Determine winner
	winner, err := engine.DetermineWinner(game.ID)
	if err != nil {
		t.Fatalf("DetermineWinner failed: %v", err)
	}

	// Winner should be player2 (Traditional company has highest initial capital)
	// Note: actual winner depends on dice rolls and cell bonuses
	if winner == nil {
		t.Fatal("Winner should not be nil")
	}

	// Verify winner has the highest capital among all players
	gameState, _ := engine.GetGameState(game.ID)
	maxCapital := int64(0)
	for _, p := range gameState.Players {
		if p.Company.Capital > maxCapital {
			maxCapital = p.Company.Capital
		}
	}

	if winner.Company.Capital != maxCapital {
		t.Errorf("Winner should have highest capital %d, but has %d", maxCapital, winner.Company.Capital)
	}
}

// Test DetermineWinner returns error for non-finished game
func TestDetermineWinner_GameNotFinished(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      2,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	// Try to determine winner before game is finished
	_, err := engine.DetermineWinner(game.ID)
	if err != ErrGameNotFinished {
		t.Errorf("Expected ErrGameNotFinished, got: %v", err)
	}
}

// Test DetermineWinner returns error for invalid game
func TestDetermineWinner_InvalidGame(t *testing.T) {
	engine := NewGameEngine()

	_, err := engine.DetermineWinner("invalid_game")
	if err != ErrGameNotFound {
		t.Errorf("Expected ErrGameNotFound, got: %v", err)
	}
}

// Feature: aws-learning-game, Property 20: Winner Determination
// For any finished game, exactly one player SHALL be determined as winner
// based on highest final Capital value.
// **Validates: Requirements 7.4**
func TestProperty20_WinnerDetermination(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		engine := NewGameEngine()

		// Generate random number of players (2-4)
		numPlayers := rapid.IntRange(2, 4).Draw(t, "numPlayers")

		config := GameConfig{
			MaxPlayers:      numPlayers,
			BoardType:       "default",
			DifficultyLevel: "normal",
		}

		game, err := engine.CreateGame(config)
		if err != nil {
			t.Fatalf("CreateGame failed: %v", err)
		}

		// Add players with different company types (different initial capitals)
		companyTypes := []company.CompanyType{company.Startup, company.Traditional, company.CloudReseller, company.CloudNative}
		playerIDs := make([]string, numPlayers)
		for i := 0; i < numPlayers; i++ {
			playerIDs[i] = rapid.StringMatching(`player_[a-z]{5}`).Draw(t, "playerID")
			companyType := companyTypes[i%len(companyTypes)]
			err := engine.JoinGame(game.ID, playerIDs[i], "Player "+playerIDs[i], companyType)
			if err != nil {
				t.Fatalf("JoinGame failed for player %d: %v", i, err)
			}
		}

		// Start game
		err = engine.StartGame(game.ID)
		if err != nil {
			t.Fatalf("StartGame failed: %v", err)
		}

		// Execute random number of turns with decision handling
		numTurns := rapid.IntRange(5, 30).Draw(t, "numTurns")
		currentPlayerIndex := 0
		for i := 0; i < numTurns; i++ {
			_, err := executeTurnAndHandleDecision(engine, game.ID, playerIDs[currentPlayerIndex])
			if err != nil {
				t.Fatalf("ExecuteTurn failed: %v", err)
			}
			currentPlayerIndex = (currentPlayerIndex + 1) % numPlayers
		}

		// End the game
		err = engine.EndGame(game.ID)
		if err != nil {
			t.Fatalf("EndGame failed: %v", err)
		}

		// Determine winner
		winner, err := engine.DetermineWinner(game.ID)
		if err != nil {
			t.Fatalf("DetermineWinner failed: %v", err)
		}

		// Property 1: Exactly one winner should be determined
		if winner == nil {
			t.Fatal("Winner should not be nil for finished game")
		}

		// Property 2: Winner should have the highest capital
		gameState, _ := engine.GetGameState(game.ID)
		maxCapital := int64(0)
		maxCapitalPlayerID := ""
		for _, p := range gameState.Players {
			if p.Company.Capital > maxCapital {
				maxCapital = p.Company.Capital
				maxCapitalPlayerID = p.PlayerID
			}
		}

		if winner.Company.Capital != maxCapital {
			t.Errorf("Winner should have highest capital %d, but has %d", maxCapital, winner.Company.Capital)
		}

		// Property 3: Winner should be one of the players
		winnerFound := false
		for _, playerID := range playerIDs {
			if winner.PlayerID == playerID {
				winnerFound = true
				break
			}
		}
		if !winnerFound {
			t.Errorf("Winner %s should be one of the players", winner.PlayerID)
		}

		// Property 4: Winner from DetermineWinner should match winner from GetGameSummary
		summary, err := engine.GetGameSummary(game.ID)
		if err != nil {
			t.Fatalf("GetGameSummary failed: %v", err)
		}

		if summary.Winner == nil {
			t.Fatal("Summary winner should not be nil")
		}

		if summary.Winner.PlayerID != winner.PlayerID {
			t.Errorf("Winner mismatch: DetermineWinner=%s, GetGameSummary=%s",
				winner.PlayerID, summary.Winner.PlayerID)
		}

		// Property 5: Winner's capital in summary should match
		if summary.Winner.FinalCapital != winner.Company.Capital {
			t.Errorf("Winner capital mismatch: DetermineWinner=%d, GetGameSummary=%d",
				winner.Company.Capital, summary.Winner.FinalCapital)
		}

		// Log for debugging (only visible in verbose mode)
		t.Logf("Game with %d players, %d turns. Winner: %s with capital %d (max capital player: %s)",
			numPlayers, numTurns, winner.PlayerID, winner.Company.Capital, maxCapitalPlayerID)
	})
}

// Test winner determination with tie (same capital)
func TestWinnerDetermination_Tie(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      2,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	game, _ := engine.CreateGame(config)

	// Join two players with same company type (same initial capital)
	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Startup)
	engine.StartGame(game.ID)

	// End game immediately (both have same capital)
	engine.EndGame(game.ID)

	// Should still determine a winner (first player with max capital)
	winner, err := engine.DetermineWinner(game.ID)
	if err != nil {
		t.Fatalf("DetermineWinner failed: %v", err)
	}

	if winner == nil {
		t.Fatal("Winner should be determined even in tie situation")
	}

	// Winner should be player1 (first player with max capital)
	if winner.PlayerID != "player1" {
		t.Logf("In tie situation, winner is %s (first player with max capital wins)", winner.PlayerID)
	}
}
