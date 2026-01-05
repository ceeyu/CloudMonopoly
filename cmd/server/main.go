package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws-learning-game/api"
	"github.com/aws-learning-game/internal/game"
	"github.com/aws-learning-game/internal/storage"
)

func main() {
	fmt.Println("AWS Learning Game Server")

	// 建立遊戲引擎
	engine := game.NewGameEngine()

	// 設定儲存 (可選)
	var gameStorage storage.GameStorage
	if os.Getenv("DYNAMODB_ENDPOINT") != "" || os.Getenv("AWS_REGION") != "" {
		ctx := context.Background()
		cfg := storage.DynamoDBConfig{
			TableName: os.Getenv("DYNAMODB_TABLE"),
			Endpoint:  os.Getenv("DYNAMODB_ENDPOINT"),
		}
		var err error
		gameStorage, err = storage.NewDynamoDBStorage(ctx, cfg)
		if err != nil {
			log.Printf("Warning: Failed to initialize DynamoDB storage: %v", err)
			log.Println("Running without persistent storage")
		}
	}

	// 設定路由
	routerConfig := api.RouterConfig{
		Engine:  engine,
		Storage: gameStorage,
		Mode:    os.Getenv("GIN_MODE"),
	}

	router := api.SetupRouter(routerConfig)

	// 取得埠號
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 啟動伺服器
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
