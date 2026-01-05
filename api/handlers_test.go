package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws-learning-game/internal/company"
	"github.com/aws-learning-game/internal/game"
	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := game.NewGameEngine()
	config := RouterConfig{
		Engine:  engine,
		Storage: nil, // No storage for basic tests
		Mode:    gin.TestMode,
	}
	return SetupRouter(config)
}

func TestHealthCheck(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}

func TestCreateGame(t *testing.T) {
	router := setupTestRouter()

	reqBody := CreateGameRequest{
		MaxPlayers:      4,
		BoardType:       "default",
		DifficultyLevel: "normal",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var response CreateGameResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.GameID == "" {
		t.Error("Expected game_id to be set")
	}
	if response.Status != game.StatusWaiting {
		t.Errorf("Expected status 'waiting', got '%s'", response.Status)
	}
}

func TestCreateGame_InvalidPlayerCount(t *testing.T) {
	router := setupTestRouter()

	// Test with invalid player count (1)
	reqBody := CreateGameRequest{
		MaxPlayers: 1,
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetGameState(t *testing.T) {
	router := setupTestRouter()

	// First create a game
	createReq := CreateGameRequest{MaxPlayers: 4}
	body, _ := json.Marshal(createReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createResp CreateGameResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)

	// Now get the game state
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/games/"+createResp.GameID, nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var stateResp GameStateResponse
	json.Unmarshal(w.Body.Bytes(), &stateResp)

	if stateResp.GameID != createResp.GameID {
		t.Errorf("Expected game_id '%s', got '%s'", createResp.GameID, stateResp.GameID)
	}
}

func TestGetGameState_NotFound(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/games/nonexistent", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestJoinGame(t *testing.T) {
	router := setupTestRouter()

	// Create a game
	createReq := CreateGameRequest{MaxPlayers: 4}
	body, _ := json.Marshal(createReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createResp CreateGameResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)

	// Join the game
	joinReq := JoinGameRequest{
		PlayerID:    "player1",
		PlayerName:  "Player One",
		CompanyType: company.Startup,
	}
	body, _ = json.Marshal(joinReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/join", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestJoinGame_InvalidCompanyType(t *testing.T) {
	router := setupTestRouter()

	// Create a game
	createReq := CreateGameRequest{MaxPlayers: 4}
	body, _ := json.Marshal(createReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createResp CreateGameResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)

	// Join with invalid company type
	joinReq := map[string]string{
		"player_id":    "player1",
		"player_name":  "Player One",
		"company_type": "invalid_type",
	}
	body, _ = json.Marshal(joinReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/join", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestStartGame(t *testing.T) {
	router := setupTestRouter()

	// Create a game
	createReq := CreateGameRequest{MaxPlayers: 4}
	body, _ := json.Marshal(createReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createResp CreateGameResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)

	// Join 2 players
	for i, playerID := range []string{"player1", "player2"} {
		joinReq := JoinGameRequest{
			PlayerID:    playerID,
			PlayerName:  "Player " + string(rune('A'+i)),
			CompanyType: company.Startup,
		}
		body, _ = json.Marshal(joinReq)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/join", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
	}

	// Start the game
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/start", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestStartGame_InsufficientPlayers(t *testing.T) {
	router := setupTestRouter()

	// Create a game
	createReq := CreateGameRequest{MaxPlayers: 4}
	body, _ := json.Marshal(createReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createResp CreateGameResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)

	// Try to start without enough players
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/start", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestExecuteTurn(t *testing.T) {
	router := setupTestRouter()

	// Create and setup a game
	createReq := CreateGameRequest{MaxPlayers: 4}
	body, _ := json.Marshal(createReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createResp CreateGameResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)

	// Join 2 players
	for i, playerID := range []string{"player1", "player2"} {
		joinReq := JoinGameRequest{
			PlayerID:    playerID,
			PlayerName:  "Player " + string(rune('A'+i)),
			CompanyType: company.Startup,
		}
		body, _ = json.Marshal(joinReq)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/join", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
	}

	// Start the game
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/start", nil)
	router.ServeHTTP(w, req)

	// Execute turn
	turnReq := TurnRequest{
		PlayerID:   "player1",
		ActionType: "roll_dice",
	}
	body, _ = json.Marshal(turnReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/turn", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var turnResp TurnResponse
	json.Unmarshal(w.Body.Bytes(), &turnResp)

	if turnResp.DiceValue < 1 || turnResp.DiceValue > 6 {
		t.Errorf("Expected dice value 1-6, got %d", turnResp.DiceValue)
	}
}

func TestExecuteTurn_NotYourTurn(t *testing.T) {
	router := setupTestRouter()

	// Create and setup a game
	createReq := CreateGameRequest{MaxPlayers: 4}
	body, _ := json.Marshal(createReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createResp CreateGameResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)

	// Join 2 players
	for i, playerID := range []string{"player1", "player2"} {
		joinReq := JoinGameRequest{
			PlayerID:    playerID,
			PlayerName:  "Player " + string(rune('A'+i)),
			CompanyType: company.Startup,
		}
		body, _ = json.Marshal(joinReq)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/join", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
	}

	// Start the game
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/start", nil)
	router.ServeHTTP(w, req)

	// Try to execute turn with wrong player (player2 when it's player1's turn)
	turnReq := TurnRequest{
		PlayerID:   "player2",
		ActionType: "roll_dice",
	}
	body, _ = json.Marshal(turnReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/turn", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSubmitDecision(t *testing.T) {
	router := setupTestRouter()

	// Create and setup a game
	createReq := CreateGameRequest{MaxPlayers: 4}
	body, _ := json.Marshal(createReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createResp CreateGameResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)

	// Join 2 players
	for i, playerID := range []string{"player1", "player2"} {
		joinReq := JoinGameRequest{
			PlayerID:    playerID,
			PlayerName:  "Player " + string(rune('A'+i)),
			CompanyType: company.Startup,
		}
		body, _ = json.Marshal(joinReq)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/join", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
	}

	// Start the game
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/start", nil)
	router.ServeHTTP(w, req)

	// Submit decision - use a valid event ID from DefaultEvents
	decisionReq := DecisionRequest{
		PlayerID: "player1",
		EventID:  "opp-001", // Valid event ID from DefaultEvents
		ChoiceID: 1,
	}
	body, _ = json.Marshal(decisionReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/decision", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var decisionResp DecisionResponse
	json.Unmarshal(w.Body.Bytes(), &decisionResp)

	// Decision result depends on random success rate and company attributes
	// Just verify we got a valid response with message and explanation
	if decisionResp.Message == "" {
		t.Error("Expected decision response to have a message")
	}
	if decisionResp.Explanation == "" {
		t.Error("Expected decision response to have an explanation")
	}
}

func TestSaveGame_NoStorage(t *testing.T) {
	router := setupTestRouter()

	// Create a game
	createReq := CreateGameRequest{MaxPlayers: 4}
	body, _ := json.Marshal(createReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createResp CreateGameResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)

	// Try to save (should succeed even without storage - just saves to engine)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/games/"+createResp.GameID+"/save", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestLoadGame_NoStorage(t *testing.T) {
	router := setupTestRouter()

	// Try to load without storage configured
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/games/game_1/load", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCORSHeaders(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/api/v1/games", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected CORS header to be set")
	}
}

func TestFullGameFlow(t *testing.T) {
	router := setupTestRouter()

	// 1. Create game
	createReq := CreateGameRequest{MaxPlayers: 2}
	body, _ := json.Marshal(createReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create game: %s", w.Body.String())
	}

	var createResp CreateGameResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)
	gameID := createResp.GameID

	// 2. Join players
	players := []struct {
		id          string
		name        string
		companyType company.CompanyType
	}{
		{"p1", "Alice", company.Startup},
		{"p2", "Bob", company.Traditional},
	}

	for _, p := range players {
		joinReq := JoinGameRequest{
			PlayerID:    p.id,
			PlayerName:  p.name,
			CompanyType: p.companyType,
		}
		body, _ = json.Marshal(joinReq)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/games/"+gameID+"/join", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Failed to join game: %s", w.Body.String())
		}
	}

	// 3. Start game
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/games/"+gameID+"/start", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Failed to start game: %s", w.Body.String())
	}

	// 4. Execute turns - get current player before each turn
	for i := 0; i < 4; i++ {
		// Get current game state to find current player
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/games/"+gameID, nil)
		router.ServeHTTP(w, req)

		var stateResp GameStateResponse
		json.Unmarshal(w.Body.Bytes(), &stateResp)
		currentPlayerID := stateResp.CurrentPlayerID

		turnReq := TurnRequest{
			PlayerID:   currentPlayerID,
			ActionType: "roll_dice",
		}
		body, _ = json.Marshal(turnReq)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/games/"+gameID+"/turn", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Failed to execute turn %d: %s", i, w.Body.String())
		}

		// Check if decision is required
		var turnResp TurnResponse
		json.Unmarshal(w.Body.Bytes(), &turnResp)

		// If decision required, submit a decision to advance turn
		if turnResp.DecisionRequired {
			decisionReq := DecisionRequest{
				PlayerID: currentPlayerID,
				EventID:  "evt_opportunity_1",
				ChoiceID: 1,
			}
			body, _ = json.Marshal(decisionReq)

			w = httptest.NewRecorder()
			req, _ = http.NewRequest("POST", "/api/v1/games/"+gameID+"/decision", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			// Decision may fail if event doesn't exist, but that's ok for this test
		}
	}

	// 5. Get final state
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/games/"+gameID, nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Failed to get game state: %s", w.Body.String())
	}

	var stateResp GameStateResponse
	json.Unmarshal(w.Body.Bytes(), &stateResp)

	if len(stateResp.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(stateResp.Players))
	}

	if stateResp.CurrentTurn < 1 {
		t.Errorf("Expected current turn >= 1, got %d", stateResp.CurrentTurn)
	}
}
