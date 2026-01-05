package game

import (
	"sync"

	"github.com/aws-learning-game/internal/board"
	"github.com/aws-learning-game/internal/company"
)

// Engine 遊戲引擎介面
type Engine interface {
	// CreateGame 建立新遊戲
	CreateGame(config GameConfig) (*Game, error)
	// JoinGame 加入遊戲
	JoinGame(gameID string, playerID string, playerName string, companyType company.CompanyType) error
	// StartGame 開始遊戲
	StartGame(gameID string) error
	// ExecuteTurn 執行回合
	ExecuteTurn(gameID string, playerID string, action TurnAction) (*TurnResult, error)
	// GetGameState 取得遊戲狀態
	GetGameState(gameID string) (*Game, error)
	// GetPlayer 取得玩家狀態
	GetPlayer(gameID string, playerID string) (*PlayerState, error)
	// GetGameSummary 取得遊戲結束摘要
	GetGameSummary(gameID string) (*GameSummary, error)
	// GetPlayerProgress 取得玩家進度
	GetPlayerProgress(gameID string, playerID string) (*PlayerProgress, error)
	// RecordDecision 記錄決策
	RecordDecision(gameID string, playerID string, record DecisionOutcomeRecord) error
	// GetCurrentPlayer 取得當前回合玩家
	GetCurrentPlayer(gameID string) (*PlayerState, error)
	// GetTurnOrder 取得回合順序
	GetTurnOrder(gameID string) ([]string, error)
	// GetPlayerTurnCount 取得各玩家已執行的回合數
	GetPlayerTurnCount(gameID string) (map[string]int, error)
	// DetermineWinner 判定贏家
	DetermineWinner(gameID string) (*PlayerState, error)
	// SaveGame 儲存遊戲狀態 (Requirements 8.1)
	SaveGame(gameID string) ([]byte, error)
	// LoadGame 載入遊戲狀態 (Requirements 8.2)
	LoadGame(data []byte) (*Game, error)
	// ImportGame 匯入遊戲到引擎
	ImportGame(game *Game) error
}

// GameEngine 遊戲引擎實作
type GameEngine struct {
	games          map[string]*Game
	boardManager   *board.BoardManager
	companyManager *company.CompanyManager
	mu             sync.RWMutex
	idCounter      int
}

// NewGameEngine 建立新的遊戲引擎
func NewGameEngine() *GameEngine {
	return &GameEngine{
		games:          make(map[string]*Game),
		boardManager:   board.NewBoardManager(),
		companyManager: company.NewCompanyManager(),
		idCounter:      0,
	}
}
