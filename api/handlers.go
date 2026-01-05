package api

import (
	"net/http"

	"github.com/aws-learning-game/internal/company"
	"github.com/aws-learning-game/internal/decision"
	"github.com/aws-learning-game/internal/event"
	"github.com/aws-learning-game/internal/game"
	"github.com/aws-learning-game/internal/storage"
	"github.com/gin-gonic/gin"
)

// Handler API 處理器
type Handler struct {
	engine         game.Engine
	storage        storage.GameStorage
	eventSystem    event.EventSystem
	decisionEngine decision.DecisionEngine
}

// NewHandler 建立新的 API 處理器
func NewHandler(engine game.Engine, storage storage.GameStorage) *Handler {
	return &Handler{
		engine:         engine,
		storage:        storage,
		eventSystem:    event.NewEventSystemWithEvents(event.DefaultEvents),
		decisionEngine: decision.NewDecisionEngine(),
	}
}

// --- Request/Response 結構 ---

// CreateGameRequest 建立遊戲請求
type CreateGameRequest struct {
	MaxPlayers      int    `json:"max_players" binding:"required,min=2,max=4"`
	BoardType       string `json:"board_type"`
	DifficultyLevel string `json:"difficulty_level"`
}

// CreateGameResponse 建立遊戲回應
type CreateGameResponse struct {
	GameID  string          `json:"game_id"`
	Status  game.GameStatus `json:"status"`
	Config  game.GameConfig `json:"config"`
	Message string          `json:"message"`
}

// JoinGameRequest 加入遊戲請求
type JoinGameRequest struct {
	PlayerID    string              `json:"player_id" binding:"required"`
	PlayerName  string              `json:"player_name" binding:"required"`
	CompanyType company.CompanyType `json:"company_type" binding:"required"`
}

// JoinGameResponse 加入遊戲回應
type JoinGameResponse struct {
	Message string `json:"message"`
	GameID  string `json:"game_id"`
}

// GameStateResponse 遊戲狀態回應
type GameStateResponse struct {
	GameID          string           `json:"game_id"`
	Status          game.GameStatus  `json:"status"`
	CurrentTurn     int              `json:"current_turn"`
	CurrentPlayerID string           `json:"current_player_id"`
	Players         []PlayerStateDTO `json:"players"`
	BoardSize       int              `json:"board_size"`
}

// PlayerStateDTO 玩家狀態 DTO
type PlayerStateDTO struct {
	PlayerID   string           `json:"player_id"`
	PlayerName string           `json:"player_name"`
	Company    *company.Company `json:"company"`
	Position   int              `json:"position"`
}

// TurnRequest 回合動作請求
type TurnRequest struct {
	PlayerID   string      `json:"player_id" binding:"required"`
	ActionType string      `json:"action_type" binding:"required"`
	Payload    interface{} `json:"payload"`
}

// TurnResponse 回合結果回應
type TurnResponse struct {
	DiceValue        int    `json:"dice_value"`
	OldPosition      int    `json:"old_position"`
	NewPosition      int    `json:"new_position"`
	CapitalChange    int64  `json:"capital_change"`
	EmployeeChange   int    `json:"employee_change"`
	CircuitCompleted bool   `json:"circuit_completed"`
	DecisionRequired bool   `json:"decision_required"`
	CellType         string `json:"cell_type"`
}

// DecisionRequest 決策請求
type DecisionRequest struct {
	PlayerID string `json:"player_id" binding:"required"`
	EventID  string `json:"event_id" binding:"required"`
	ChoiceID int    `json:"choice_id" binding:"required"`
}

// DecisionResponse 決策回應
type DecisionResponse struct {
	Success         bool     `json:"success"`
	Message         string   `json:"message"`
	Explanation     string   `json:"explanation"`
	CapitalChange   int64    `json:"capital_change"`
	EmployeeChange  int      `json:"employee_change"`
	LearningPoints  []string `json:"learning_points"`
	AWSBestPractice string   `json:"aws_best_practice"`
}

// EventResponse 事件回應
type EventResponse struct {
	ID           string           `json:"id"`
	Type         string           `json:"type"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	RealWorldRef string           `json:"real_world_ref"`
	Context      EventContextDTO  `json:"context"`
	Choices      []EventChoiceDTO `json:"choices"`
	AWSTopics    []string         `json:"aws_topics"`
}

// EventContextDTO 事件背景 DTO
type EventContextDTO struct {
	Scenario       string   `json:"scenario"`
	BusinessImpact string   `json:"business_impact"`
	TechnicalNeeds []string `json:"technical_needs"`
	Constraints    []string `json:"constraints"`
}

// EventChoiceDTO 事件選項 DTO
type EventChoiceDTO struct {
	ID                  int                   `json:"id"`
	Title               string                `json:"title"`
	Description         string                `json:"description"`
	IsAWS               bool                  `json:"is_aws"`
	AWSServices         []string              `json:"aws_services"`
	OnPremSolution      string                `json:"on_prem_solution"`
	Requirements        ChoiceRequirementsDTO `json:"requirements"`
	Outcomes            ChoiceOutcomesDTO     `json:"outcomes"`
	ArchitectureDiagram string                `json:"architecture_diagram"`
}

// ChoiceRequirementsDTO 選項需求 DTO
type ChoiceRequirementsDTO struct {
	MinCapital       int64    `json:"min_capital"`
	MinEmployees     int      `json:"min_employees"`
	MinSecurityLevel int      `json:"min_security_level"`
	RequiredInfra    []string `json:"required_infra"`
}

// ChoiceOutcomesDTO 選項結果 DTO
type ChoiceOutcomesDTO struct {
	CapitalChange       int64   `json:"capital_change"`
	EmployeeChange      int     `json:"employee_change"`
	SecurityChange      int     `json:"security_change"`
	CloudAdoptionChange float64 `json:"cloud_adoption_change"`
	SuccessRate         float64 `json:"success_rate"`
	TimeToImplement     int     `json:"time_to_implement"`
}

// SaveGameResponse 儲存遊戲回應
type SaveGameResponse struct {
	Message string `json:"message"`
	GameID  string `json:"game_id"`
}

// ErrorResponse 錯誤回應
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// --- 遊戲管理 API ---

// CreateGame 建立遊戲
// POST /games
// Requirements 9.1: REST API for game operations
func (h *Handler) CreateGame(c *gin.Context) {
	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Code:    "400",
			Message: err.Error(),
		})
		return
	}

	// 設定預設值
	if req.BoardType == "" {
		req.BoardType = "default"
	}
	if req.DifficultyLevel == "" {
		req.DifficultyLevel = "normal"
	}

	config := game.GameConfig{
		MaxPlayers:      req.MaxPlayers,
		BoardType:       req.BoardType,
		DifficultyLevel: req.DifficultyLevel,
	}

	g, err := h.engine.CreateGame(config)
	if err != nil {
		h.handleGameError(c, err)
		return
	}

	c.JSON(http.StatusCreated, CreateGameResponse{
		GameID:  g.ID,
		Status:  g.Status,
		Config:  g.Config,
		Message: "遊戲建立成功",
	})
}

// GetGameState 取得遊戲狀態
// GET /games/:id
// Requirements 9.1: REST API for game operations
func (h *Handler) GetGameState(c *gin.Context) {
	gameID := c.Param("id")

	g, err := h.engine.GetGameState(gameID)
	if err != nil {
		h.handleGameError(c, err)
		return
	}

	// 轉換為 DTO
	players := make([]PlayerStateDTO, len(g.Players))
	for i, p := range g.Players {
		players[i] = PlayerStateDTO{
			PlayerID:   p.PlayerID,
			PlayerName: p.PlayerName,
			Company:    p.Company,
			Position:   p.Position,
		}
	}

	boardSize := 0
	if g.Board != nil {
		boardSize = g.Board.Size
	}

	c.JSON(http.StatusOK, GameStateResponse{
		GameID:          g.ID,
		Status:          g.Status,
		CurrentTurn:     g.CurrentTurn,
		CurrentPlayerID: g.CurrentPlayerID,
		Players:         players,
		BoardSize:       boardSize,
	})
}

// JoinGame 加入遊戲
// POST /games/:id/join
// Requirements 9.1: REST API for game operations
func (h *Handler) JoinGame(c *gin.Context) {
	gameID := c.Param("id")

	var req JoinGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Code:    "400",
			Message: err.Error(),
		})
		return
	}

	// 驗證公司類型
	if !isValidCompanyType(req.CompanyType) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_COMPANY_TYPE",
			Code:    "400",
			Message: "無效的公司類型，請選擇: startup, traditional, cloud_reseller, cloud_native",
		})
		return
	}

	err := h.engine.JoinGame(gameID, req.PlayerID, req.PlayerName, req.CompanyType)
	if err != nil {
		h.handleGameError(c, err)
		return
	}

	c.JSON(http.StatusOK, JoinGameResponse{
		Message: "成功加入遊戲",
		GameID:  gameID,
	})
}

// StartGame 開始遊戲
// POST /games/:id/start
// Requirements 9.1: REST API for game operations
func (h *Handler) StartGame(c *gin.Context) {
	gameID := c.Param("id")

	err := h.engine.StartGame(gameID)
	if err != nil {
		h.handleGameError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "遊戲開始",
		"game_id": gameID,
	})
}

// --- 回合操作 API ---

// ExecuteTurn 執行回合動作
// POST /games/:id/turn
// Requirements 9.1, 9.4: REST API for turn operations
func (h *Handler) ExecuteTurn(c *gin.Context) {
	gameID := c.Param("id")

	var req TurnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Code:    "400",
			Message: err.Error(),
		})
		return
	}

	action := game.TurnAction{
		ActionType: req.ActionType,
		Payload:    req.Payload,
	}

	result, err := h.engine.ExecuteTurn(gameID, req.PlayerID, action)
	if err != nil {
		h.handleGameError(c, err)
		return
	}

	c.JSON(http.StatusOK, TurnResponse{
		DiceValue:        result.DiceValue,
		OldPosition:      result.OldPosition,
		NewPosition:      result.NewPosition,
		CapitalChange:    result.CapitalChange,
		EmployeeChange:   result.EmployeeChange,
		CircuitCompleted: result.CircuitCompleted,
		DecisionRequired: result.DecisionRequired,
		CellType:         result.CellType,
	})
}

// SubmitDecision 提交決策
// POST /games/:id/decision
// Requirements 9.1, 9.4: REST API for decision submission
func (h *Handler) SubmitDecision(c *gin.Context) {
	gameID := c.Param("id")

	var req DecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Code:    "400",
			Message: err.Error(),
		})
		return
	}

	// 取得玩家資訊
	player, err := h.engine.GetPlayer(gameID, req.PlayerID)
	if err != nil {
		h.handleGameError(c, err)
		return
	}

	// 取得事件
	evt, err := h.eventSystem.GetEvent(req.EventID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "EVENT_NOT_FOUND",
			Code:    "404",
			Message: "事件不存在",
		})
		return
	}

	// 找到選擇的選項
	var selectedChoice *event.EventChoice
	for i := range evt.Choices {
		if evt.Choices[i].ID == req.ChoiceID {
			selectedChoice = &evt.Choices[i]
			break
		}
	}
	if selectedChoice == nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_CHOICE",
			Code:    "400",
			Message: "無效的選項",
		})
		return
	}

	// 執行決策
	result, err := h.decisionEngine.ExecuteDecision(player.Company, selectedChoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "DECISION_FAILED",
			Code:    "500",
			Message: err.Error(),
		})
		return
	}

	// 記錄決策
	record := game.DecisionOutcomeRecord{
		DecisionRecord: game.DecisionRecord{
			EventID:  req.EventID,
			ChoiceID: req.ChoiceID,
		},
		Success:        result.Success,
		AWSServices:    selectedChoice.AWSServices,
		LearningPoints: result.LearningPoints,
	}

	err = h.engine.RecordDecision(gameID, req.PlayerID, record)
	if err != nil {
		h.handleGameError(c, err)
		return
	}

	c.JSON(http.StatusOK, DecisionResponse{
		Success:         result.Success,
		Message:         result.Message,
		Explanation:     result.Explanation,
		CapitalChange:   result.ActualOutcome.CapitalChange,
		EmployeeChange:  result.ActualOutcome.EmployeeChange,
		LearningPoints:  result.LearningPoints,
		AWSBestPractice: result.AWSBestPractice,
	})
}

// GetEvent 取得事件
// GET /games/:id/event/:position
// Requirements 9.3: Display events
func (h *Handler) GetEvent(c *gin.Context) {
	gameID := c.Param("id")
	eventID := c.Param("eventId")

	// 如果提供了 eventId，直接取得該事件
	if eventID != "" {
		evt, err := h.eventSystem.GetEvent(eventID)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "EVENT_NOT_FOUND",
				Code:    "404",
				Message: "事件不存在",
			})
			return
		}
		c.JSON(http.StatusOK, h.convertEventToDTO(evt))
		return
	}

	// 否則根據遊戲狀態取得當前玩家位置的事件
	g, err := h.engine.GetGameState(gameID)
	if err != nil {
		h.handleGameError(c, err)
		return
	}

	// 找到當前玩家
	var currentPlayer *game.PlayerState
	for _, p := range g.Players {
		if p.PlayerID == g.CurrentPlayerID {
			currentPlayer = p
			break
		}
	}
	if currentPlayer == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "PLAYER_NOT_FOUND",
			Code:    "404",
			Message: "找不到當前玩家",
		})
		return
	}

	// 取得當前位置的格子類型
	if g.Board == nil || currentPlayer.Position >= len(g.Board.Cells) {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INVALID_BOARD",
			Code:    "500",
			Message: "棋盤狀態無效",
		})
		return
	}

	cell := g.Board.Cells[currentPlayer.Position]

	// 根據格子類型取得對應事件
	evt, err := h.eventSystem.GetEventForCell(cell.Type, currentPlayer.Company)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "NO_EVENT",
			Code:    "404",
			Message: "此格子沒有事件",
		})
		return
	}

	c.JSON(http.StatusOK, h.convertEventToDTO(evt))
}

// GetRandomEvent 取得隨機事件
// GET /games/:id/event/random/:type
func (h *Handler) GetRandomEvent(c *gin.Context) {
	gameID := c.Param("id")
	eventType := c.Param("type")

	// 取得遊戲狀態以獲取玩家公司資訊
	g, err := h.engine.GetGameState(gameID)
	if err != nil {
		h.handleGameError(c, err)
		return
	}

	// 找到當前玩家
	var currentPlayer *game.PlayerState
	for _, p := range g.Players {
		if p.PlayerID == g.CurrentPlayerID {
			currentPlayer = p
			break
		}
	}

	var comp *company.Company
	if currentPlayer != nil {
		comp = currentPlayer.Company
	}

	// 轉換事件類型
	var evtType event.EventType
	switch eventType {
	case "opportunity":
		evtType = event.EventOpportunity
	case "fate":
		evtType = event.EventFate
	case "challenge":
		evtType = event.EventChallenge
	case "security":
		evtType = event.EventSecurity
	default:
		evtType = event.EventOpportunity
	}

	evt, err := h.eventSystem.GetRandomEvent(evtType, comp)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "NO_EVENT",
			Code:    "404",
			Message: "沒有可用的事件",
		})
		return
	}

	c.JSON(http.StatusOK, h.convertEventToDTO(evt))
}

// convertEventToDTO 轉換事件為 DTO
func (h *Handler) convertEventToDTO(evt *event.Event) EventResponse {
	choices := make([]EventChoiceDTO, len(evt.Choices))
	for i, choice := range evt.Choices {
		choices[i] = EventChoiceDTO{
			ID:             choice.ID,
			Title:          choice.Title,
			Description:    choice.Description,
			IsAWS:          choice.IsAWS,
			AWSServices:    choice.AWSServices,
			OnPremSolution: choice.OnPremSolution,
			Requirements: ChoiceRequirementsDTO{
				MinCapital:       choice.Requirements.MinCapital,
				MinEmployees:     choice.Requirements.MinEmployees,
				MinSecurityLevel: choice.Requirements.MinSecurityLevel,
				RequiredInfra:    choice.Requirements.RequiredInfra,
			},
			Outcomes: ChoiceOutcomesDTO{
				CapitalChange:       choice.Outcomes.CapitalChange,
				EmployeeChange:      choice.Outcomes.EmployeeChange,
				SecurityChange:      choice.Outcomes.SecurityChange,
				CloudAdoptionChange: choice.Outcomes.CloudAdoptionChange,
				SuccessRate:         choice.Outcomes.SuccessRate,
				TimeToImplement:     choice.Outcomes.TimeToImplement,
			},
			ArchitectureDiagram: choice.ArchitectureDiagram,
		}
	}

	return EventResponse{
		ID:           evt.ID,
		Type:         string(evt.Type),
		Title:        evt.Title,
		Description:  evt.Description,
		RealWorldRef: evt.RealWorldRef,
		Context: EventContextDTO{
			Scenario:       evt.Context.Scenario,
			BusinessImpact: evt.Context.BusinessImpact,
			TechnicalNeeds: evt.Context.TechnicalNeeds,
			Constraints:    evt.Context.Constraints,
		},
		Choices:   choices,
		AWSTopics: evt.AWSTopics,
	}
}

// --- 存檔 API ---

// SaveGame 儲存遊戲
// POST /games/:id/save
// Requirements 8.1, 8.2: Save game state
func (h *Handler) SaveGame(c *gin.Context) {
	gameID := c.Param("id")

	// 取得遊戲狀態
	g, err := h.engine.GetGameState(gameID)
	if err != nil {
		h.handleGameError(c, err)
		return
	}

	// 儲存到持久化儲存
	if h.storage != nil {
		err = h.storage.SaveGame(c.Request.Context(), g)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "SAVE_FAILED",
				Code:    "500",
				Message: "儲存遊戲失敗: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, SaveGameResponse{
		Message: "遊戲儲存成功",
		GameID:  gameID,
	})
}

// LoadGame 載入遊戲
// GET /games/:id/load
// Requirements 8.1, 8.2: Load game state
func (h *Handler) LoadGame(c *gin.Context) {
	gameID := c.Param("id")

	if h.storage == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "STORAGE_NOT_CONFIGURED",
			Code:    "500",
			Message: "儲存服務未設定",
		})
		return
	}

	// 從持久化儲存載入
	g, err := h.storage.LoadGame(c.Request.Context(), gameID)
	if err != nil {
		h.handleGameError(c, err)
		return
	}

	// 匯入到遊戲引擎
	err = h.engine.ImportGame(g)
	if err != nil {
		h.handleGameError(c, err)
		return
	}

	// 轉換為回應格式
	players := make([]PlayerStateDTO, len(g.Players))
	for i, p := range g.Players {
		players[i] = PlayerStateDTO{
			PlayerID:   p.PlayerID,
			PlayerName: p.PlayerName,
			Company:    p.Company,
			Position:   p.Position,
		}
	}

	boardSize := 0
	if g.Board != nil {
		boardSize = g.Board.Size
	}

	c.JSON(http.StatusOK, GameStateResponse{
		GameID:          g.ID,
		Status:          g.Status,
		CurrentTurn:     g.CurrentTurn,
		CurrentPlayerID: g.CurrentPlayerID,
		Players:         players,
		BoardSize:       boardSize,
	})
}

// --- 輔助函數 ---

// handleGameError 處理遊戲錯誤
func (h *Handler) handleGameError(c *gin.Context, err error) {
	switch err {
	case game.ErrGameNotFound:
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "GAME_NOT_FOUND",
			Code:    "404",
			Message: "遊戲不存在",
		})
	case game.ErrGameFull:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "GAME_FULL",
			Code:    "400",
			Message: "遊戲已滿",
		})
	case game.ErrGameAlreadyStarted:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "GAME_ALREADY_STARTED",
			Code:    "400",
			Message: "遊戲已開始",
		})
	case game.ErrNotYourTurn:
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "NOT_YOUR_TURN",
			Code:    "403",
			Message: "不是你的回合",
		})
	case game.ErrInvalidAction:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_ACTION",
			Code:    "400",
			Message: "無效的動作",
		})
	case game.ErrPlayerNotFound:
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "PLAYER_NOT_FOUND",
			Code:    "404",
			Message: "玩家不存在",
		})
	case game.ErrInsufficientPlayers:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INSUFFICIENT_PLAYERS",
			Code:    "400",
			Message: "玩家人數不足 (需要 2-4 人)",
		})
	case game.ErrInvalidPlayerCount:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_PLAYER_COUNT",
			Code:    "400",
			Message: "無效的玩家數量配置 (需要 2-4 人)",
		})
	case game.ErrSaveFailed:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "SAVE_FAILED",
			Code:    "500",
			Message: "儲存遊戲失敗",
		})
	case game.ErrLoadFailed:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "LOAD_FAILED",
			Code:    "500",
			Message: "載入遊戲失敗",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Code:    "500",
			Message: err.Error(),
		})
	}
}

// isValidCompanyType 驗證公司類型
func isValidCompanyType(ct company.CompanyType) bool {
	switch ct {
	case company.Startup, company.Traditional, company.CloudReseller, company.CloudNative:
		return true
	default:
		return false
	}
}
