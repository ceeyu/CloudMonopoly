package storage

import (
	"testing"
	"time"

	"github.com/aws-learning-game/internal/board"
	"github.com/aws-learning-game/internal/company"
	"github.com/aws-learning-game/internal/game"
	"pgregory.net/rapid"
)

// executeTurnAndHandleDecision 執行回合並在需要時自動處理決策
func executeTurnAndHandleDecision(engine *game.GameEngine, gameID string, playerID string) (*game.TurnResult, error) {
	action := game.TurnAction{ActionType: "roll_dice"}
	result, err := engine.ExecuteTurn(gameID, playerID, action)
	if err != nil {
		return nil, err
	}

	// 如果需要決策，自動提交一個空決策來切換玩家
	if result.DecisionRequired {
		record := game.DecisionOutcomeRecord{
			DecisionRecord: game.DecisionRecord{
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

// Feature: aws-learning-game, Property 22: Save and Load Consistency
// For any saved game, loading it SHALL restore a GameState where all player positions,
// company attributes, and turn number match the saved state.
// **Validates: Requirements 8.1, 8.2**
func TestProperty22_SaveAndLoadConsistency(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Create a game engine
		engine := game.NewGameEngine()

		// Generate random config
		maxPlayers := rapid.IntRange(2, 4).Draw(t, "maxPlayers")
		config := game.GameConfig{
			MaxPlayers:      maxPlayers,
			BoardType:       "default",
			DifficultyLevel: "normal",
		}

		// Create game
		g, err := engine.CreateGame(config)
		if err != nil {
			t.Fatalf("CreateGame failed: %v", err)
		}

		// Generate random number of players (2 to maxPlayers)
		numPlayers := rapid.IntRange(2, maxPlayers).Draw(t, "numPlayers")
		companyTypes := []company.CompanyType{
			company.Startup,
			company.Traditional,
			company.CloudReseller,
			company.CloudNative,
		}

		playerIDs := make([]string, numPlayers)
		for i := 0; i < numPlayers; i++ {
			playerIDs[i] = rapid.StringMatching(`player_[a-z]{5}`).Draw(t, "playerID")
			companyType := companyTypes[rapid.IntRange(0, len(companyTypes)-1).Draw(t, "companyTypeIndex")]
			err := engine.JoinGame(g.ID, playerIDs[i], "Player "+playerIDs[i], companyType)
			if err != nil {
				t.Fatalf("JoinGame failed: %v", err)
			}
		}

		// Start game
		err = engine.StartGame(g.ID)
		if err != nil {
			t.Fatalf("StartGame failed: %v", err)
		}

		// Execute random number of turns with decision handling
		numTurns := rapid.IntRange(1, 20).Draw(t, "numTurns")
		currentPlayerIndex := 0
		for i := 0; i < numTurns; i++ {
			_, err := executeTurnAndHandleDecision(engine, g.ID, playerIDs[currentPlayerIndex])
			if err != nil {
				t.Fatalf("ExecuteTurn failed: %v", err)
			}
			currentPlayerIndex = (currentPlayerIndex + 1) % numPlayers
		}

		// Get game state before save
		gameStateBefore, err := engine.GetGameState(g.ID)
		if err != nil {
			t.Fatalf("GetGameState failed: %v", err)
		}

		// Save game
		savedData, err := engine.SaveGame(g.ID)
		if err != nil {
			t.Fatalf("SaveGame failed: %v", err)
		}

		// Verify saved data is not empty
		if len(savedData) == 0 {
			t.Fatal("Saved data should not be empty")
		}

		// Create a new engine and load the game
		newEngine := game.NewGameEngine()
		loadedGame, err := newEngine.LoadGame(savedData)
		if err != nil {
			t.Fatalf("LoadGame failed: %v", err)
		}

		// Verify loaded game matches saved state
		// Check game ID
		if loadedGame.ID != gameStateBefore.ID {
			t.Errorf("Game ID mismatch: expected %s, got %s", gameStateBefore.ID, loadedGame.ID)
		}

		// Check status
		if loadedGame.Status != gameStateBefore.Status {
			t.Errorf("Status mismatch: expected %s, got %s", gameStateBefore.Status, loadedGame.Status)
		}

		// Check current turn
		if loadedGame.CurrentTurn != gameStateBefore.CurrentTurn {
			t.Errorf("CurrentTurn mismatch: expected %d, got %d", gameStateBefore.CurrentTurn, loadedGame.CurrentTurn)
		}

		// Check current player ID
		if loadedGame.CurrentPlayerID != gameStateBefore.CurrentPlayerID {
			t.Errorf("CurrentPlayerID mismatch: expected %s, got %s", gameStateBefore.CurrentPlayerID, loadedGame.CurrentPlayerID)
		}

		// Check player count
		if len(loadedGame.Players) != len(gameStateBefore.Players) {
			t.Errorf("Player count mismatch: expected %d, got %d", len(gameStateBefore.Players), len(loadedGame.Players))
		}

		// Check each player's state
		for i, playerBefore := range gameStateBefore.Players {
			playerAfter := loadedGame.Players[i]

			// Check player ID
			if playerAfter.PlayerID != playerBefore.PlayerID {
				t.Errorf("Player %d ID mismatch: expected %s, got %s", i, playerBefore.PlayerID, playerAfter.PlayerID)
			}

			// Check position
			if playerAfter.Position != playerBefore.Position {
				t.Errorf("Player %d position mismatch: expected %d, got %d", i, playerBefore.Position, playerAfter.Position)
			}

			// Check turns played
			if playerAfter.TurnsPlayed != playerBefore.TurnsPlayed {
				t.Errorf("Player %d turns played mismatch: expected %d, got %d", i, playerBefore.TurnsPlayed, playerAfter.TurnsPlayed)
			}

			// Check company attributes
			if playerAfter.Company.Capital != playerBefore.Company.Capital {
				t.Errorf("Player %d capital mismatch: expected %d, got %d", i, playerBefore.Company.Capital, playerAfter.Company.Capital)
			}

			if playerAfter.Company.Employees != playerBefore.Company.Employees {
				t.Errorf("Player %d employees mismatch: expected %d, got %d", i, playerBefore.Company.Employees, playerAfter.Company.Employees)
			}

			if playerAfter.Company.Type != playerBefore.Company.Type {
				t.Errorf("Player %d company type mismatch: expected %s, got %s", i, playerBefore.Company.Type, playerAfter.Company.Type)
			}

			if playerAfter.Company.SecurityLevel != playerBefore.Company.SecurityLevel {
				t.Errorf("Player %d security level mismatch: expected %d, got %d", i, playerBefore.Company.SecurityLevel, playerAfter.Company.SecurityLevel)
			}

			if playerAfter.Company.CloudAdoption != playerBefore.Company.CloudAdoption {
				t.Errorf("Player %d cloud adoption mismatch: expected %f, got %f", i, playerBefore.Company.CloudAdoption, playerAfter.Company.CloudAdoption)
			}
		}

		// Check board
		if loadedGame.Board.Size != gameStateBefore.Board.Size {
			t.Errorf("Board size mismatch: expected %d, got %d", gameStateBefore.Board.Size, loadedGame.Board.Size)
		}
	})
}

// Test save and load with decision history
func TestSaveLoadWithDecisionHistory(t *testing.T) {
	engine := game.NewGameEngine()

	config := game.GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	// Create and setup game
	g, _ := engine.CreateGame(config)
	engine.JoinGame(g.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(g.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(g.ID)

	// Execute some turns
	engine.ExecuteTurn(g.ID, "player1", game.TurnAction{ActionType: "roll_dice"})
	engine.ExecuteTurn(g.ID, "player2", game.TurnAction{ActionType: "roll_dice"})

	// Record a decision
	record := game.DecisionOutcomeRecord{
		DecisionRecord: game.DecisionRecord{
			TurnNumber: 1,
			EventID:    "event_test",
			ChoiceID:   1,
			Timestamp:  time.Now(),
		},
		Success:        true,
		AWSServices:    []string{"EC2", "S3"},
		LearningPoints: []string{"了解 EC2 執行個體類型"},
	}
	engine.RecordDecision(g.ID, "player1", record)

	// Get state before save
	stateBefore, _ := engine.GetGameState(g.ID)

	// Save
	savedData, err := engine.SaveGame(g.ID)
	if err != nil {
		t.Fatalf("SaveGame failed: %v", err)
	}

	// Load in new engine
	newEngine := game.NewGameEngine()
	loadedGame, err := newEngine.LoadGame(savedData)
	if err != nil {
		t.Fatalf("LoadGame failed: %v", err)
	}

	// Verify decision history is preserved
	if len(loadedGame.Players[0].DecisionHistory) != len(stateBefore.Players[0].DecisionHistory) {
		t.Errorf("Decision history count mismatch: expected %d, got %d",
			len(stateBefore.Players[0].DecisionHistory),
			len(loadedGame.Players[0].DecisionHistory))
	}

	if len(loadedGame.Players[0].DecisionOutcomeHistory) != len(stateBefore.Players[0].DecisionOutcomeHistory) {
		t.Errorf("Decision outcome history count mismatch: expected %d, got %d",
			len(stateBefore.Players[0].DecisionOutcomeHistory),
			len(loadedGame.Players[0].DecisionOutcomeHistory))
	}

	// Verify decision details
	if len(loadedGame.Players[0].DecisionOutcomeHistory) > 0 {
		loadedOutcome := loadedGame.Players[0].DecisionOutcomeHistory[0]
		originalOutcome := stateBefore.Players[0].DecisionOutcomeHistory[0]

		if loadedOutcome.EventID != originalOutcome.EventID {
			t.Errorf("EventID mismatch: expected %s, got %s", originalOutcome.EventID, loadedOutcome.EventID)
		}

		if loadedOutcome.Success != originalOutcome.Success {
			t.Errorf("Success mismatch: expected %v, got %v", originalOutcome.Success, loadedOutcome.Success)
		}

		if len(loadedOutcome.AWSServices) != len(originalOutcome.AWSServices) {
			t.Errorf("AWSServices count mismatch: expected %d, got %d",
				len(originalOutcome.AWSServices), len(loadedOutcome.AWSServices))
		}
	}
}

// Test save and load preserves game config
func TestSaveLoadPreservesConfig(t *testing.T) {
	engine := game.NewGameEngine()

	config := game.GameConfig{
		MaxPlayers:      3,
		BoardType:       "default",
		DifficultyLevel: "hard",
	}

	g, _ := engine.CreateGame(config)
	engine.JoinGame(g.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(g.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(g.ID)

	// Save
	savedData, _ := engine.SaveGame(g.ID)

	// Load
	newEngine := game.NewGameEngine()
	loadedGame, _ := newEngine.LoadGame(savedData)

	// Verify config
	if loadedGame.Config.MaxPlayers != config.MaxPlayers {
		t.Errorf("MaxPlayers mismatch: expected %d, got %d", config.MaxPlayers, loadedGame.Config.MaxPlayers)
	}

	if loadedGame.Config.BoardType != config.BoardType {
		t.Errorf("BoardType mismatch: expected %s, got %s", config.BoardType, loadedGame.Config.BoardType)
	}

	if loadedGame.Config.DifficultyLevel != config.DifficultyLevel {
		t.Errorf("DifficultyLevel mismatch: expected %s, got %s", config.DifficultyLevel, loadedGame.Config.DifficultyLevel)
	}
}

// Test that loaded game can continue playing
func TestLoadedGameCanContinue(t *testing.T) {
	engine := game.NewGameEngine()

	config := game.GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	g, _ := engine.CreateGame(config)
	engine.JoinGame(g.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(g.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(g.ID)

	// Execute some turns with decision handling
	executeTurnAndHandleDecision(engine, g.ID, "player1")

	// Save
	savedData, _ := engine.SaveGame(g.ID)

	// Load in new engine
	newEngine := game.NewGameEngine()
	loadedGame, _ := newEngine.LoadGame(savedData)

	// Verify we can continue playing
	// Current player should be player2
	if loadedGame.CurrentPlayerID != "player2" {
		t.Errorf("Expected current player to be player2, got %s", loadedGame.CurrentPlayerID)
	}

	// Execute turn on loaded game with decision handling
	result, err := executeTurnAndHandleDecision(newEngine, loadedGame.ID, "player2")
	if err != nil {
		t.Fatalf("ExecuteTurn on loaded game failed: %v", err)
	}

	// Verify turn executed successfully
	if result.DiceValue < 1 || result.DiceValue > 6 {
		t.Errorf("Invalid dice value: %d", result.DiceValue)
	}

	// Verify game state updated
	updatedState, _ := newEngine.GetGameState(loadedGame.ID)
	if updatedState.CurrentPlayerID != "player1" {
		t.Errorf("Expected current player to be player1 after turn, got %s", updatedState.CurrentPlayerID)
	}
}

// Helper function to generate random game for testing
func generateTestGame(t *rapid.T) *game.Game {
	gameID := rapid.StringMatching(`game_[a-z0-9]{8}`).Draw(t, "gameID")
	maxPlayers := rapid.IntRange(2, 4).Draw(t, "maxPlayers")

	status := rapid.SampledFrom([]game.GameStatus{
		game.StatusWaiting,
		game.StatusInProgress,
		game.StatusFinished,
	}).Draw(t, "status")

	numPlayers := rapid.IntRange(1, maxPlayers).Draw(t, "numPlayers")
	players := make([]*game.PlayerState, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = &game.PlayerState{
			PlayerID:   rapid.StringMatching(`player_[a-z]{5}`).Draw(t, "playerID"),
			PlayerName: rapid.StringMatching(`Player_[A-Z][a-z]{4}`).Draw(t, "playerName"),
			Company: &company.Company{
				ID:            rapid.StringMatching(`company_[a-z0-9]{6}`).Draw(t, "companyID"),
				Name:          rapid.StringMatching(`[A-Z][a-z]{3,8} Corp`).Draw(t, "companyName"),
				Type:          company.Startup,
				Capital:       int64(rapid.IntRange(100, 10000).Draw(t, "capital")),
				Employees:     rapid.IntRange(5, 500).Draw(t, "employees"),
				SecurityLevel: rapid.IntRange(1, 5).Draw(t, "securityLevel"),
				CloudAdoption: float64(rapid.IntRange(0, 100).Draw(t, "cloudAdoption")),
			},
			Position:               rapid.IntRange(0, 29).Draw(t, "position"),
			TurnsPlayed:            rapid.IntRange(0, 50).Draw(t, "turnsPlayed"),
			DecisionHistory:        []game.DecisionRecord{},
			DecisionOutcomeHistory: []game.DecisionOutcomeRecord{},
		}
	}

	currentPlayerID := ""
	if numPlayers > 0 {
		currentPlayerID = players[rapid.IntRange(0, numPlayers-1).Draw(t, "currentPlayerIndex")].PlayerID
	}

	boardSize := rapid.IntRange(30, 50).Draw(t, "boardSize")
	cells := make([]board.Cell, boardSize)
	for i := 0; i < boardSize; i++ {
		cells[i] = board.Cell{
			Position:      i,
			Type:          board.CellNormal,
			Name:          rapid.StringMatching(`Cell_[0-9]{2}`).Draw(t, "cellName"),
			BaseCapital:   int64(rapid.IntRange(0, 100).Draw(t, "baseCapital")),
			BaseEmployees: rapid.IntRange(0, 10).Draw(t, "baseEmployees"),
		}
	}
	cells[0].Type = board.CellStart

	return &game.Game{
		ID: gameID,
		Config: game.GameConfig{
			MaxPlayers:      maxPlayers,
			BoardType:       "default",
			DifficultyLevel: "normal",
		},
		Status:          status,
		CurrentTurn:     rapid.IntRange(0, 100).Draw(t, "currentTurn"),
		CurrentPlayerID: currentPlayerID,
		Players:         players,
		Board: &board.Board{
			ID:    rapid.StringMatching(`board_[a-z0-9]{6}`).Draw(t, "boardID"),
			Name:  "Test Board",
			Cells: cells,
			Size:  boardSize,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
