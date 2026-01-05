package game

import (
	"encoding/json"
	"time"

	"github.com/aws-learning-game/internal/board"
	"github.com/aws-learning-game/internal/company"
)

// GameStateJSON 用於 JSON 序列化的遊戲狀態結構
// 使用明確的 JSON tags 確保序列化/反序列化一致性
type GameStateJSON struct {
	ID              string            `json:"id"`
	Config          GameConfigJSON    `json:"config"`
	Status          GameStatus        `json:"status"`
	CurrentTurn     int               `json:"current_turn"`
	CurrentPlayerID string            `json:"current_player_id"`
	Players         []PlayerStateJSON `json:"players"`
	Board           *board.Board      `json:"board"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// GameConfigJSON 遊戲配置 JSON 結構
type GameConfigJSON struct {
	MaxPlayers      int    `json:"max_players"`
	BoardType       string `json:"board_type"`
	DifficultyLevel string `json:"difficulty_level"`
}

// PlayerStateJSON 玩家狀態 JSON 結構
type PlayerStateJSON struct {
	PlayerID               string                      `json:"player_id"`
	PlayerName             string                      `json:"player_name"`
	Company                *company.Company            `json:"company"`
	Position               int                         `json:"position"`
	TurnsPlayed            int                         `json:"turns_played"`
	DecisionHistory        []DecisionRecordJSON        `json:"decision_history"`
	DecisionOutcomeHistory []DecisionOutcomeRecordJSON `json:"decision_outcome_history"`
}

// DecisionRecordJSON 決策記錄 JSON 結構
type DecisionRecordJSON struct {
	TurnNumber int       `json:"turn_number"`
	EventID    string    `json:"event_id"`
	ChoiceID   int       `json:"choice_id"`
	Timestamp  time.Time `json:"timestamp"`
}

// DecisionOutcomeRecordJSON 決策結果記錄 JSON 結構
type DecisionOutcomeRecordJSON struct {
	DecisionRecordJSON
	Success        bool     `json:"success"`
	AWSServices    []string `json:"aws_services"`
	LearningPoints []string `json:"learning_points"`
}

// SerializeGameState 將遊戲狀態序列化為 JSON
// Requirements 8.4: 使用 JSON 格式編碼遊戲狀態
func SerializeGameState(game *Game) ([]byte, error) {
	if game == nil {
		return nil, ErrGameNotFound
	}

	// 轉換為 JSON 結構
	jsonState := GameStateJSON{
		ID: game.ID,
		Config: GameConfigJSON{
			MaxPlayers:      game.Config.MaxPlayers,
			BoardType:       game.Config.BoardType,
			DifficultyLevel: game.Config.DifficultyLevel,
		},
		Status:          game.Status,
		CurrentTurn:     game.CurrentTurn,
		CurrentPlayerID: game.CurrentPlayerID,
		Board:           game.Board,
		CreatedAt:       game.CreatedAt,
		UpdatedAt:       game.UpdatedAt,
	}

	// 轉換玩家狀態
	jsonState.Players = make([]PlayerStateJSON, len(game.Players))
	for i, player := range game.Players {
		jsonState.Players[i] = convertPlayerToJSON(player)
	}

	return json.Marshal(jsonState)
}

// DeserializeGameState 將 JSON 反序列化為遊戲狀態
// Requirements 8.5: 驗證並還原完整的遊戲狀態
func DeserializeGameState(data []byte) (*Game, error) {
	if len(data) == 0 {
		return nil, ErrCorruptedState
	}

	var jsonState GameStateJSON
	if err := json.Unmarshal(data, &jsonState); err != nil {
		return nil, ErrCorruptedState
	}

	// 驗證必要欄位
	if err := validateGameStateJSON(&jsonState); err != nil {
		return nil, err
	}

	// 轉換回 Game 結構
	game := &Game{
		ID: jsonState.ID,
		Config: GameConfig{
			MaxPlayers:      jsonState.Config.MaxPlayers,
			BoardType:       jsonState.Config.BoardType,
			DifficultyLevel: jsonState.Config.DifficultyLevel,
		},
		Status:          jsonState.Status,
		CurrentTurn:     jsonState.CurrentTurn,
		CurrentPlayerID: jsonState.CurrentPlayerID,
		Board:           jsonState.Board,
		CreatedAt:       jsonState.CreatedAt,
		UpdatedAt:       jsonState.UpdatedAt,
	}

	// 轉換玩家狀態
	game.Players = make([]*PlayerState, len(jsonState.Players))
	for i, playerJSON := range jsonState.Players {
		game.Players[i] = convertJSONToPlayer(&playerJSON)
	}

	return game, nil
}

// convertPlayerToJSON 將 PlayerState 轉換為 JSON 結構
func convertPlayerToJSON(player *PlayerState) PlayerStateJSON {
	if player == nil {
		return PlayerStateJSON{}
	}

	jsonPlayer := PlayerStateJSON{
		PlayerID:    player.PlayerID,
		PlayerName:  player.PlayerName,
		Company:     player.Company,
		Position:    player.Position,
		TurnsPlayed: player.TurnsPlayed,
	}

	// 轉換決策歷史
	jsonPlayer.DecisionHistory = make([]DecisionRecordJSON, len(player.DecisionHistory))
	for i, record := range player.DecisionHistory {
		jsonPlayer.DecisionHistory[i] = DecisionRecordJSON{
			TurnNumber: record.TurnNumber,
			EventID:    record.EventID,
			ChoiceID:   record.ChoiceID,
			Timestamp:  record.Timestamp,
		}
	}

	// 轉換決策結果歷史
	jsonPlayer.DecisionOutcomeHistory = make([]DecisionOutcomeRecordJSON, len(player.DecisionOutcomeHistory))
	for i, outcome := range player.DecisionOutcomeHistory {
		jsonPlayer.DecisionOutcomeHistory[i] = DecisionOutcomeRecordJSON{
			DecisionRecordJSON: DecisionRecordJSON{
				TurnNumber: outcome.TurnNumber,
				EventID:    outcome.EventID,
				ChoiceID:   outcome.ChoiceID,
				Timestamp:  outcome.Timestamp,
			},
			Success:        outcome.Success,
			AWSServices:    outcome.AWSServices,
			LearningPoints: outcome.LearningPoints,
		}
	}

	return jsonPlayer
}

// convertJSONToPlayer 將 JSON 結構轉換為 PlayerState
func convertJSONToPlayer(jsonPlayer *PlayerStateJSON) *PlayerState {
	if jsonPlayer == nil {
		return nil
	}

	player := &PlayerState{
		PlayerID:    jsonPlayer.PlayerID,
		PlayerName:  jsonPlayer.PlayerName,
		Company:     jsonPlayer.Company,
		Position:    jsonPlayer.Position,
		TurnsPlayed: jsonPlayer.TurnsPlayed,
	}

	// 轉換決策歷史
	player.DecisionHistory = make([]DecisionRecord, len(jsonPlayer.DecisionHistory))
	for i, record := range jsonPlayer.DecisionHistory {
		player.DecisionHistory[i] = DecisionRecord{
			TurnNumber: record.TurnNumber,
			EventID:    record.EventID,
			ChoiceID:   record.ChoiceID,
			Timestamp:  record.Timestamp,
		}
	}

	// 轉換決策結果歷史
	player.DecisionOutcomeHistory = make([]DecisionOutcomeRecord, len(jsonPlayer.DecisionOutcomeHistory))
	for i, outcome := range jsonPlayer.DecisionOutcomeHistory {
		player.DecisionOutcomeHistory[i] = DecisionOutcomeRecord{
			DecisionRecord: DecisionRecord{
				TurnNumber: outcome.TurnNumber,
				EventID:    outcome.EventID,
				ChoiceID:   outcome.ChoiceID,
				Timestamp:  outcome.Timestamp,
			},
			Success:        outcome.Success,
			AWSServices:    outcome.AWSServices,
			LearningPoints: outcome.LearningPoints,
		}
	}

	return player
}

// validateGameStateJSON 驗證 JSON 遊戲狀態的完整性
func validateGameStateJSON(state *GameStateJSON) error {
	if state.ID == "" {
		return ErrCorruptedState
	}

	// 驗證狀態值
	switch state.Status {
	case StatusWaiting, StatusInProgress, StatusFinished:
		// 有效狀態
	default:
		return ErrCorruptedState
	}

	// 驗證棋盤
	if state.Board == nil {
		return ErrCorruptedState
	}

	// 驗證玩家配置
	if state.Config.MaxPlayers < 2 || state.Config.MaxPlayers > 4 {
		return ErrCorruptedState
	}

	return nil
}

// GameStateEquals 比較兩個遊戲狀態是否相等
// 用於測試 round-trip 序列化
func GameStateEquals(a, b *Game) bool {
	if a == nil || b == nil {
		return a == b
	}

	// 比較基本欄位
	if a.ID != b.ID ||
		a.Status != b.Status ||
		a.CurrentTurn != b.CurrentTurn ||
		a.CurrentPlayerID != b.CurrentPlayerID {
		return false
	}

	// 比較配置
	if a.Config.MaxPlayers != b.Config.MaxPlayers ||
		a.Config.BoardType != b.Config.BoardType ||
		a.Config.DifficultyLevel != b.Config.DifficultyLevel {
		return false
	}

	// 比較玩家數量
	if len(a.Players) != len(b.Players) {
		return false
	}

	// 比較每個玩家
	for i := range a.Players {
		if !playerStateEquals(a.Players[i], b.Players[i]) {
			return false
		}
	}

	// 比較棋盤
	if !boardEquals(a.Board, b.Board) {
		return false
	}

	return true
}

// playerStateEquals 比較兩個玩家狀態是否相等
func playerStateEquals(a, b *PlayerState) bool {
	if a == nil || b == nil {
		return a == b
	}

	if a.PlayerID != b.PlayerID ||
		a.PlayerName != b.PlayerName ||
		a.Position != b.Position ||
		a.TurnsPlayed != b.TurnsPlayed {
		return false
	}

	// 比較公司
	if !companyEquals(a.Company, b.Company) {
		return false
	}

	// 比較決策歷史
	if len(a.DecisionHistory) != len(b.DecisionHistory) {
		return false
	}
	for i := range a.DecisionHistory {
		if !decisionRecordEquals(&a.DecisionHistory[i], &b.DecisionHistory[i]) {
			return false
		}
	}

	// 比較決策結果歷史
	if len(a.DecisionOutcomeHistory) != len(b.DecisionOutcomeHistory) {
		return false
	}
	for i := range a.DecisionOutcomeHistory {
		if !decisionOutcomeRecordEquals(&a.DecisionOutcomeHistory[i], &b.DecisionOutcomeHistory[i]) {
			return false
		}
	}

	return true
}

// companyEquals 比較兩個公司是否相等
func companyEquals(a, b *company.Company) bool {
	if a == nil || b == nil {
		return a == b
	}

	if a.ID != b.ID ||
		a.Name != b.Name ||
		a.Type != b.Type ||
		a.Capital != b.Capital ||
		a.Employees != b.Employees ||
		a.IsInternational != b.IsInternational ||
		a.ProductCycle != b.ProductCycle ||
		a.TechDebt != b.TechDebt ||
		a.SecurityLevel != b.SecurityLevel ||
		a.CloudAdoption != b.CloudAdoption {
		return false
	}

	// 比較基礎設施
	if len(a.Infrastructure) != len(b.Infrastructure) {
		return false
	}
	for i := range a.Infrastructure {
		if a.Infrastructure[i] != b.Infrastructure[i] {
			return false
		}
	}

	return true
}

// decisionRecordEquals 比較兩個決策記錄是否相等
func decisionRecordEquals(a, b *DecisionRecord) bool {
	if a == nil || b == nil {
		return a == b
	}

	return a.TurnNumber == b.TurnNumber &&
		a.EventID == b.EventID &&
		a.ChoiceID == b.ChoiceID &&
		a.Timestamp.Equal(b.Timestamp)
}

// decisionOutcomeRecordEquals 比較兩個決策結果記錄是否相等
func decisionOutcomeRecordEquals(a, b *DecisionOutcomeRecord) bool {
	if a == nil || b == nil {
		return a == b
	}

	if !decisionRecordEquals(&a.DecisionRecord, &b.DecisionRecord) {
		return false
	}

	if a.Success != b.Success {
		return false
	}

	// 比較 AWS 服務
	if len(a.AWSServices) != len(b.AWSServices) {
		return false
	}
	for i := range a.AWSServices {
		if a.AWSServices[i] != b.AWSServices[i] {
			return false
		}
	}

	// 比較學習要點
	if len(a.LearningPoints) != len(b.LearningPoints) {
		return false
	}
	for i := range a.LearningPoints {
		if a.LearningPoints[i] != b.LearningPoints[i] {
			return false
		}
	}

	return true
}

// boardEquals 比較兩個棋盤是否相等
func boardEquals(a, b *board.Board) bool {
	if a == nil || b == nil {
		return a == b
	}

	if a.ID != b.ID ||
		a.Name != b.Name ||
		a.Size != b.Size {
		return false
	}

	if len(a.Cells) != len(b.Cells) {
		return false
	}

	for i := range a.Cells {
		if a.Cells[i].Position != b.Cells[i].Position ||
			a.Cells[i].Type != b.Cells[i].Type ||
			a.Cells[i].Name != b.Cells[i].Name ||
			a.Cells[i].BaseCapital != b.Cells[i].BaseCapital ||
			a.Cells[i].BaseEmployees != b.Cells[i].BaseEmployees ||
			a.Cells[i].EventID != b.Cells[i].EventID {
			return false
		}
	}

	return true
}
