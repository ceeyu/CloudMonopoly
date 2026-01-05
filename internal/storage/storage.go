package storage

import (
	"context"

	"github.com/aws-learning-game/internal/game"
)

// GameStorage 遊戲儲存介面
// Requirements 8.1, 8.2, 8.3: 使用 AWS 服務儲存遊戲狀態
type GameStorage interface {
	// SaveGame 儲存遊戲狀態
	SaveGame(ctx context.Context, g *game.Game) error
	// LoadGame 載入遊戲狀態
	LoadGame(ctx context.Context, gameID string) (*game.Game, error)
	// DeleteGame 刪除遊戲狀態
	DeleteGame(ctx context.Context, gameID string) error
	// ListGames 列出所有遊戲
	ListGames(ctx context.Context) ([]*GameMetadata, error)
}

// GameMetadata 遊戲元資料 (用於列表顯示)
type GameMetadata struct {
	GameID      string `json:"game_id"`
	Status      string `json:"status"`
	PlayerCount int    `json:"player_count"`
	CurrentTurn int    `json:"current_turn"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
