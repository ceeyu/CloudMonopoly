package api

import (
	"github.com/aws-learning-game/internal/game"
	"github.com/aws-learning-game/internal/storage"
	"github.com/gin-gonic/gin"
)

// RouterConfig 路由配置
type RouterConfig struct {
	Engine  game.Engine
	Storage storage.GameStorage
	Mode    string // "debug", "release", "test"
}

// SetupRouter 設定路由
// Requirements 9.1: REST API for all game operations
func SetupRouter(config RouterConfig) *gin.Engine {
	// 設定 Gin 模式
	if config.Mode != "" {
		gin.SetMode(config.Mode)
	}

	router := gin.New()

	// 中介軟體
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	// 建立處理器
	handler := NewHandler(config.Engine, config.Storage)

	// API 路由群組
	api := router.Group("/api/v1")
	{
		// 健康檢查
		api.GET("/health", HealthCheck)

		// 遊戲管理 API (Requirements 9.1)
		games := api.Group("/games")
		{
			// POST /games - 建立遊戲
			games.POST("", handler.CreateGame)

			// GET /games/:id - 取得遊戲狀態
			games.GET("/:id", handler.GetGameState)

			// POST /games/:id/join - 加入遊戲
			games.POST("/:id/join", handler.JoinGame)

			// POST /games/:id/start - 開始遊戲
			games.POST("/:id/start", handler.StartGame)

			// 回合操作 API (Requirements 9.1, 9.4)
			// POST /games/:id/turn - 執行回合動作
			games.POST("/:id/turn", handler.ExecuteTurn)

			// POST /games/:id/decision - 提交決策
			games.POST("/:id/decision", handler.SubmitDecision)

			// 事件 API (Requirements 9.3)
			// GET /games/:id/event - 取得當前位置事件
			games.GET("/:id/event", handler.GetEvent)

			// GET /games/:id/event/:eventId - 取得指定事件
			games.GET("/:id/event/:eventId", handler.GetEvent)

			// GET /games/:id/event/random/:type - 取得隨機事件
			games.GET("/:id/event/random/:type", handler.GetRandomEvent)

			// 存檔 API (Requirements 8.1, 8.2)
			// POST /games/:id/save - 儲存遊戲
			games.POST("/:id/save", handler.SaveGame)

			// GET /games/:id/load - 載入遊戲
			games.GET("/:id/load", handler.LoadGame)
		}
	}

	return router
}

// HealthCheck 健康檢查端點
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "aws-learning-game",
	})
}

// CORSMiddleware CORS 中介軟體
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
