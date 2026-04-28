package game

import (
	"testing"
	"time"

	"github.com/aws-learning-game/internal/board"
	"github.com/aws-learning-game/internal/company"
	"pgregory.net/rapid"
)

// Feature: aws-learning-game, Property 21: Game State Serialization Round-Trip
// For any valid GameState, serializing to JSON then deserializing SHALL produce
// a GameState equivalent to the original (all fields match).
// **Validates: Requirements 8.4, 8.5**
func TestProperty21_GameStateSerializationRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random game state
		game := generateRandomGame(t)

		// Serialize to JSON
		jsonData, err := SerializeGameState(game)
		if err != nil {
			t.Fatalf("SerializeGameState failed: %v", err)
		}

		// Verify JSON is not empty
		if len(jsonData) == 0 {
			t.Fatal("Serialized JSON should not be empty")
		}

		// Deserialize back to Game
		restored, err := DeserializeGameState(jsonData)
		if err != nil {
			t.Fatalf("DeserializeGameState failed: %v", err)
		}

		// Verify the restored game equals the original
		if !GameStateEquals(game, restored) {
			t.Errorf("Round-trip failed: original and restored game states are not equal")
			t.Logf("Original ID: %s, Restored ID: %s", game.ID, restored.ID)
			t.Logf("Original Status: %s, Restored Status: %s", game.Status, restored.Status)
			t.Logf("Original Turn: %d, Restored Turn: %d", game.CurrentTurn, restored.CurrentTurn)
			t.Logf("Original Players: %d, Restored Players: %d", len(game.Players), len(restored.Players))
		}
	})
}

// generateRandomGame 產生隨機遊戲狀態用於測試
func generateRandomGame(t *rapid.T) *Game {
	// Generate random game ID
	gameID := rapid.StringMatching(`game_[a-z0-9]{8}`).Draw(t, "gameID")

	// Generate random config
	maxPlayers := rapid.IntRange(2, 4).Draw(t, "maxPlayers")
	boardType := rapid.SampledFrom([]string{"default", "advanced", "simple"}).Draw(t, "boardType")
	difficulty := rapid.SampledFrom([]string{"easy", "normal", "hard"}).Draw(t, "difficulty")

	// Generate random status
	status := rapid.SampledFrom([]GameStatus{StatusWaiting, StatusInProgress, StatusFinished}).Draw(t, "status")

	// Generate random turn
	currentTurn := rapid.IntRange(0, 100).Draw(t, "currentTurn")

	// Generate random number of players (1 to maxPlayers)
	numPlayers := rapid.IntRange(1, maxPlayers).Draw(t, "numPlayers")

	// Generate players
	players := make([]*PlayerState, numPlayers)
	playerIDs := make([]string, numPlayers)
	for i := 0; i < numPlayers; i++ {
		playerIDs[i] = rapid.StringMatching(`player_[a-z]{5}`).Draw(t, "playerID")
		players[i] = generateRandomPlayer(t, playerIDs[i])
	}

	// Select current player
	currentPlayerID := ""
	if numPlayers > 0 {
		currentPlayerIndex := rapid.IntRange(0, numPlayers-1).Draw(t, "currentPlayerIndex")
		currentPlayerID = playerIDs[currentPlayerIndex]
	}

	// Generate board
	gameBoard := generateRandomBoard(t)

	// Generate timestamps
	now := time.Now()
	createdAt := now.Add(-time.Duration(rapid.IntRange(0, 3600).Draw(t, "createdOffset")) * time.Second)
	updatedAt := now.Add(-time.Duration(rapid.IntRange(0, 60).Draw(t, "updatedOffset")) * time.Second)

	return &Game{
		ID: gameID,
		Config: GameConfig{
			MaxPlayers:      maxPlayers,
			BoardType:       boardType,
			DifficultyLevel: difficulty,
		},
		Status:            status,
		CurrentTurn:       currentTurn,
		CurrentPlayerID:   currentPlayerID,
		Players:           players,
		Board:             gameBoard,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		MaxTurnsPerPlayer: DefaultMaxTurnsPerPlayer, // 預設 30 回合
	}
}

// generateRandomPlayer 產生隨機玩家狀態
func generateRandomPlayer(t *rapid.T, playerID string) *PlayerState {
	playerName := rapid.StringMatching(`Player_[A-Z][a-z]{4}`).Draw(t, "playerName")

	// Generate company
	comp := generateRandomCompany(t)

	// Generate position
	position := rapid.IntRange(0, 29).Draw(t, "position")

	// Generate turns played
	turnsPlayed := rapid.IntRange(0, 50).Draw(t, "turnsPlayed")

	// Generate decision history
	numDecisions := rapid.IntRange(0, 10).Draw(t, "numDecisions")
	decisionHistory := make([]DecisionRecord, numDecisions)
	for i := 0; i < numDecisions; i++ {
		decisionHistory[i] = generateRandomDecisionRecord(t, i+1)
	}

	// Generate decision outcome history
	numOutcomes := rapid.IntRange(0, numDecisions).Draw(t, "numOutcomes")
	decisionOutcomeHistory := make([]DecisionOutcomeRecord, numOutcomes)
	for i := 0; i < numOutcomes; i++ {
		decisionOutcomeHistory[i] = generateRandomDecisionOutcomeRecord(t, i+1)
	}

	return &PlayerState{
		PlayerID:               playerID,
		PlayerName:             playerName,
		Company:                comp,
		Position:               position,
		TurnsPlayed:            turnsPlayed,
		DecisionHistory:        decisionHistory,
		DecisionOutcomeHistory: decisionOutcomeHistory,
	}
}

// generateRandomCompany 產生隨機公司
func generateRandomCompany(t *rapid.T) *company.Company {
	companyTypes := []company.CompanyType{
		company.Startup,
		company.Traditional,
		company.CloudReseller,
		company.CloudNative,
	}
	companyType := rapid.SampledFrom(companyTypes).Draw(t, "companyType")

	companyID := rapid.StringMatching(`company_[a-z0-9]{6}`).Draw(t, "companyID")
	companyName := rapid.StringMatching(`[A-Z][a-z]{3,8} Corp`).Draw(t, "companyName")

	capital := int64(rapid.IntRange(100, 10000).Draw(t, "capital"))
	employees := rapid.IntRange(5, 500).Draw(t, "employees")
	isInternational := rapid.Bool().Draw(t, "isInternational")

	productCycles := []string{"development", "launch", "growth", "mature"}
	productCycle := rapid.SampledFrom(productCycles).Draw(t, "productCycle")

	techDebt := rapid.IntRange(0, 100).Draw(t, "techDebt")
	securityLevel := rapid.IntRange(1, 5).Draw(t, "securityLevel")
	cloudAdoption := float64(rapid.IntRange(0, 100).Draw(t, "cloudAdoption"))

	// Generate infrastructure
	infraOptions := []string{"EC2", "S3", "RDS", "Lambda", "DynamoDB", "VPC", "CloudFront"}
	numInfra := rapid.IntRange(0, len(infraOptions)).Draw(t, "numInfra")
	infrastructure := make([]string, numInfra)
	for i := 0; i < numInfra; i++ {
		infrastructure[i] = rapid.SampledFrom(infraOptions).Draw(t, "infra")
	}

	return &company.Company{
		ID:              companyID,
		Name:            companyName,
		Type:            companyType,
		Capital:         capital,
		Employees:       employees,
		IsInternational: isInternational,
		ProductCycle:    productCycle,
		TechDebt:        techDebt,
		SecurityLevel:   securityLevel,
		CloudAdoption:   cloudAdoption,
		Infrastructure:  infrastructure,
	}
}

// generateRandomBoard 產生隨機棋盤
func generateRandomBoard(t *rapid.T) *board.Board {
	boardID := rapid.StringMatching(`board_[a-z0-9]{6}`).Draw(t, "boardID")
	boardName := rapid.SampledFrom([]string{"Default Board", "Advanced Board", "Simple Board"}).Draw(t, "boardName")

	// Board size must be at least 30 per requirements
	boardSize := rapid.IntRange(30, 50).Draw(t, "boardSize")

	cellTypes := []board.CellType{
		board.CellNormal,
		board.CellOpportunity,
		board.CellFate,
		board.CellChallenge,
		board.CellStart,
		board.CellBonus,
	}

	cells := make([]board.Cell, boardSize)
	for i := 0; i < boardSize; i++ {
		cellType := rapid.SampledFrom(cellTypes).Draw(t, "cellType")
		cells[i] = board.Cell{
			Position:      i,
			Type:          cellType,
			Name:          rapid.StringMatching(`Cell_[0-9]{2}`).Draw(t, "cellName"),
			BaseCapital:   int64(rapid.IntRange(0, 100).Draw(t, "baseCapital")),
			BaseEmployees: rapid.IntRange(0, 10).Draw(t, "baseEmployees"),
			EventID:       rapid.StringMatching(`event_[a-z0-9]{4}`).Draw(t, "eventID"),
		}
	}

	// Ensure first cell is start
	cells[0].Type = board.CellStart

	return &board.Board{
		ID:    boardID,
		Name:  boardName,
		Cells: cells,
		Size:  boardSize,
	}
}

// generateRandomDecisionRecord 產生隨機決策記錄
func generateRandomDecisionRecord(t *rapid.T, turnNumber int) DecisionRecord {
	return DecisionRecord{
		TurnNumber: turnNumber,
		EventID:    rapid.StringMatching(`event_[a-z0-9]{6}`).Draw(t, "eventID"),
		ChoiceID:   rapid.IntRange(1, 5).Draw(t, "choiceID"),
		Timestamp:  time.Now().Add(-time.Duration(rapid.IntRange(0, 3600).Draw(t, "timestampOffset")) * time.Second),
	}
}

// generateRandomDecisionOutcomeRecord 產生隨機決策結果記錄
func generateRandomDecisionOutcomeRecord(t *rapid.T, turnNumber int) DecisionOutcomeRecord {
	record := generateRandomDecisionRecord(t, turnNumber)

	// Generate AWS services
	awsServices := []string{"EC2", "S3", "RDS", "Lambda", "DynamoDB", "VPC", "CloudFront", "IAM", "KMS"}
	numServices := rapid.IntRange(0, 3).Draw(t, "numServices")
	services := make([]string, numServices)
	for i := 0; i < numServices; i++ {
		services[i] = rapid.SampledFrom(awsServices).Draw(t, "service")
	}

	// Generate learning points
	learningPoints := []string{
		"了解 EC2 執行個體類型選擇",
		"學習 S3 儲存類別差異",
		"掌握 RDS 多可用區部署",
		"理解 Lambda 無伺服器架構",
		"熟悉 VPC 網路設計",
	}
	numPoints := rapid.IntRange(0, 3).Draw(t, "numPoints")
	points := make([]string, numPoints)
	for i := 0; i < numPoints; i++ {
		points[i] = rapid.SampledFrom(learningPoints).Draw(t, "point")
	}

	return DecisionOutcomeRecord{
		DecisionRecord: record,
		Success:        rapid.Bool().Draw(t, "success"),
		AWSServices:    services,
		LearningPoints: points,
	}
}

// Test serialization with nil game
func TestSerializeGameState_NilGame(t *testing.T) {
	_, err := SerializeGameState(nil)
	if err != ErrGameNotFound {
		t.Errorf("Expected ErrGameNotFound for nil game, got: %v", err)
	}
}

// Test deserialization with empty data
func TestDeserializeGameState_EmptyData(t *testing.T) {
	_, err := DeserializeGameState([]byte{})
	if err != ErrCorruptedState {
		t.Errorf("Expected ErrCorruptedState for empty data, got: %v", err)
	}
}

// Test deserialization with invalid JSON
func TestDeserializeGameState_InvalidJSON(t *testing.T) {
	_, err := DeserializeGameState([]byte("not valid json"))
	if err != ErrCorruptedState {
		t.Errorf("Expected ErrCorruptedState for invalid JSON, got: %v", err)
	}
}

// Test deserialization with missing required fields
func TestDeserializeGameState_MissingFields(t *testing.T) {
	// JSON with empty ID
	invalidJSON := `{"id":"","status":"waiting","config":{"max_players":2},"board":{"id":"b1","name":"test","cells":[],"size":30}}`
	_, err := DeserializeGameState([]byte(invalidJSON))
	if err != ErrCorruptedState {
		t.Errorf("Expected ErrCorruptedState for missing ID, got: %v", err)
	}
}

// Test basic serialization and deserialization
func TestSerializeDeserialize_Basic(t *testing.T) {
	engine := NewGameEngine()

	config := GameConfig{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}

	// Create and setup a game
	game, err := engine.CreateGame(config)
	if err != nil {
		t.Fatalf("CreateGame failed: %v", err)
	}

	engine.JoinGame(game.ID, "player1", "Player 1", company.Startup)
	engine.JoinGame(game.ID, "player2", "Player 2", company.Traditional)
	engine.StartGame(game.ID)

	// Execute some turns
	engine.ExecuteTurn(game.ID, "player1", TurnAction{ActionType: "roll_dice"})
	engine.ExecuteTurn(game.ID, "player2", TurnAction{ActionType: "roll_dice"})

	// Get current game state
	gameState, _ := engine.GetGameState(game.ID)

	// Serialize
	jsonData, err := SerializeGameState(gameState)
	if err != nil {
		t.Fatalf("SerializeGameState failed: %v", err)
	}

	// Deserialize
	restored, err := DeserializeGameState(jsonData)
	if err != nil {
		t.Fatalf("DeserializeGameState failed: %v", err)
	}

	// Verify key fields
	if restored.ID != gameState.ID {
		t.Errorf("ID mismatch: expected %s, got %s", gameState.ID, restored.ID)
	}
	if restored.Status != gameState.Status {
		t.Errorf("Status mismatch: expected %s, got %s", gameState.Status, restored.Status)
	}
	if restored.CurrentTurn != gameState.CurrentTurn {
		t.Errorf("CurrentTurn mismatch: expected %d, got %d", gameState.CurrentTurn, restored.CurrentTurn)
	}
	if len(restored.Players) != len(gameState.Players) {
		t.Errorf("Players count mismatch: expected %d, got %d", len(gameState.Players), len(restored.Players))
	}

	// Verify player data
	for i, player := range gameState.Players {
		restoredPlayer := restored.Players[i]
		if restoredPlayer.PlayerID != player.PlayerID {
			t.Errorf("Player %d ID mismatch: expected %s, got %s", i, player.PlayerID, restoredPlayer.PlayerID)
		}
		if restoredPlayer.Position != player.Position {
			t.Errorf("Player %d Position mismatch: expected %d, got %d", i, player.Position, restoredPlayer.Position)
		}
		if restoredPlayer.Company.Capital != player.Company.Capital {
			t.Errorf("Player %d Capital mismatch: expected %d, got %d", i, player.Company.Capital, restoredPlayer.Company.Capital)
		}
	}
}
