package game

import (
	"time"

	"github.com/aws-learning-game/internal/board"
	"github.com/aws-learning-game/internal/company"
)

// GameStatus 遊戲狀態
type GameStatus string

const (
	StatusWaiting    GameStatus = "waiting"     // 等待玩家加入
	StatusInProgress GameStatus = "in_progress" // 遊戲進行中
	StatusFinished   GameStatus = "finished"    // 遊戲結束
)

// GameConfig 遊戲配置
type GameConfig struct {
	MaxPlayers      int
	BoardType       string
	DifficultyLevel string
}

// DefaultGameConfig 預設遊戲配置
var DefaultGameConfig = GameConfig{
	MaxPlayers:      4,
	BoardType:       "default",
	DifficultyLevel: "normal",
}

// Game 遊戲實體
type Game struct {
	ID              string
	Config          GameConfig
	Status          GameStatus
	CurrentTurn     int
	CurrentPlayerID string
	Players         []*PlayerState
	Board           *board.Board
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PlayerState 玩家狀態
type PlayerState struct {
	PlayerID               string
	PlayerName             string
	Company                *company.Company
	Position               int
	TurnsPlayed            int
	DecisionHistory        []DecisionRecord
	DecisionOutcomeHistory []DecisionOutcomeRecord // 詳細決策結果歷史
}

// DecisionRecord 決策記錄
type DecisionRecord struct {
	TurnNumber int
	EventID    string
	ChoiceID   int
	Timestamp  time.Time
}

// GameSummary 遊戲結束摘要
type GameSummary struct {
	GameID           string          `json:"game_id"`
	TotalTurns       int             `json:"total_turns"`
	Winner           *PlayerSummary  `json:"winner"`
	PlayerSummaries  []PlayerSummary `json:"player_summaries"`
	LearningInsights []string        `json:"learning_insights"`
}

// PlayerSummary 玩家摘要
type PlayerSummary struct {
	PlayerID        string                  `json:"player_id"`
	PlayerName      string                  `json:"player_name"`
	CompanyType     company.CompanyType     `json:"company_type"`
	FinalCapital    int64                   `json:"final_capital"`
	FinalEmployees  int                     `json:"final_employees"`
	TurnsPlayed     int                     `json:"turns_played"`
	DecisionCount   int                     `json:"decision_count"`
	DecisionRecords []DecisionRecordSummary `json:"decision_records"`
}

// DecisionRecordSummary 決策記錄摘要
type DecisionRecordSummary struct {
	TurnNumber int       `json:"turn_number"`
	EventID    string    `json:"event_id"`
	ChoiceID   int       `json:"choice_id"`
	Timestamp  time.Time `json:"timestamp"`
}

// PlayerProgress 玩家進度追蹤
type PlayerProgress struct {
	PlayerID            string         `json:"player_id"`
	TotalDecisions      int            `json:"total_decisions"`
	SuccessfulDecisions int            `json:"successful_decisions"`
	FailedDecisions     int            `json:"failed_decisions"`
	AWSServicesUsed     []string       `json:"aws_services_used"`
	TopicsLearned       []string       `json:"topics_learned"`
	SkillAreas          map[string]int `json:"skill_areas"`       // 技能領域分數
	LearningProgress    float64        `json:"learning_progress"` // 學習進度百分比
	Recommendations     []string       `json:"recommendations"`   // 改進建議
}

// DecisionOutcomeRecord 決策結果記錄 (擴展 DecisionRecord)
type DecisionOutcomeRecord struct {
	DecisionRecord
	Success        bool     `json:"success"`
	AWSServices    []string `json:"aws_services"`
	LearningPoints []string `json:"learning_points"`
}
