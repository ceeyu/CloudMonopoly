package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aws-learning-game/internal/company"
)

// generateID 產生唯一 ID
func (e *GameEngine) generateID() string {
	e.idCounter++
	return fmt.Sprintf("game_%d", e.idCounter)
}

// CreateGame 建立新遊戲
func (e *GameEngine) CreateGame(config GameConfig) (*Game, error) {
	// 驗證玩家數量配置
	if config.MaxPlayers < 2 || config.MaxPlayers > 4 {
		return nil, ErrInvalidPlayerCount
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// 建立棋盤
	gameBoard, err := e.boardManager.CreateBoard(config.BoardType)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	game := &Game{
		ID:          e.generateID(),
		Config:      config,
		Status:      StatusWaiting,
		CurrentTurn: 0,
		Players:     make([]*PlayerState, 0, config.MaxPlayers),
		Board:       gameBoard,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	e.games[game.ID] = game
	return game, nil
}

// JoinGame 加入遊戲
func (e *GameEngine) JoinGame(gameID string, playerID string, playerName string, companyType company.CompanyType) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, ok := e.games[gameID]
	if !ok {
		return ErrGameNotFound
	}

	if game.Status != StatusWaiting {
		return ErrGameAlreadyStarted
	}

	if len(game.Players) >= game.Config.MaxPlayers {
		return ErrGameFull
	}

	// 檢查玩家是否已在遊戲中
	for _, p := range game.Players {
		if p.PlayerID == playerID {
			return ErrPlayerAlreadyInGame
		}
	}

	// 建立公司
	comp, err := e.companyManager.CreateCompany(companyType)
	if err != nil {
		return err
	}

	player := &PlayerState{
		PlayerID:               playerID,
		PlayerName:             playerName,
		Company:                comp,
		Position:               0,
		TurnsPlayed:            0,
		DecisionHistory:        make([]DecisionRecord, 0),
		DecisionOutcomeHistory: make([]DecisionOutcomeRecord, 0),
	}

	game.Players = append(game.Players, player)
	game.UpdatedAt = time.Now()

	return nil
}

// StartGame 開始遊戲
func (e *GameEngine) StartGame(gameID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, ok := e.games[gameID]
	if !ok {
		return ErrGameNotFound
	}

	if game.Status != StatusWaiting {
		return ErrGameAlreadyStarted
	}

	// 驗證玩家數量 (2-4人)
	if len(game.Players) < 2 {
		return ErrInsufficientPlayers
	}
	if len(game.Players) > 4 {
		return ErrTooManyPlayers
	}

	game.Status = StatusInProgress
	game.CurrentTurn = 1
	game.CurrentPlayerID = game.Players[0].PlayerID
	game.UpdatedAt = time.Now()

	return nil
}

// GetGameState 取得遊戲狀態
func (e *GameEngine) GetGameState(gameID string) (*Game, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, ok := e.games[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}

	// 返回副本以避免外部修改
	gameCopy := *game
	gameCopy.Players = make([]*PlayerState, len(game.Players))
	for i, p := range game.Players {
		playerCopy := *p
		if p.Company != nil {
			companyCopy := *p.Company
			playerCopy.Company = &companyCopy
		}
		playerCopy.DecisionHistory = make([]DecisionRecord, len(p.DecisionHistory))
		copy(playerCopy.DecisionHistory, p.DecisionHistory)
		playerCopy.DecisionOutcomeHistory = make([]DecisionOutcomeRecord, len(p.DecisionOutcomeHistory))
		copy(playerCopy.DecisionOutcomeHistory, p.DecisionOutcomeHistory)
		gameCopy.Players[i] = &playerCopy
	}

	return &gameCopy, nil
}

// ExecuteTurn 執行回合 (基礎實作，後續任務會擴充)
func (e *GameEngine) ExecuteTurn(gameID string, playerID string, action TurnAction) (*TurnResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, ok := e.games[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}

	if game.Status != StatusInProgress {
		return nil, ErrGameNotStarted
	}

	if game.CurrentPlayerID != playerID {
		return nil, ErrNotYourTurn
	}

	// 找到當前玩家
	var currentPlayer *PlayerState
	for _, p := range game.Players {
		if p.PlayerID == playerID {
			currentPlayer = p
			break
		}
	}
	if currentPlayer == nil {
		return nil, ErrPlayerNotFound
	}

	result := &TurnResult{}

	switch action.ActionType {
	case "roll_dice":
		result = e.handleRollDice(game, currentPlayer)
	default:
		return nil, ErrInvalidAction
	}

	game.UpdatedAt = time.Now()
	return result, nil
}

// handleRollDice 處理擲骰子動作
func (e *GameEngine) handleRollDice(game *Game, player *PlayerState) *TurnResult {
	// 擲骰子 (1-6)
	diceValue := rand.Intn(6) + 1

	oldPosition := player.Position
	boardSize := game.Board.Size

	// 計算新位置
	newPosition := (oldPosition + diceValue) % boardSize

	// 檢查是否繞圈 (經過起點)
	circuitCompleted := (oldPosition + diceValue) >= boardSize

	result := &TurnResult{
		DiceValue:        diceValue,
		OldPosition:      oldPosition,
		NewPosition:      newPosition,
		CircuitCompleted: circuitCompleted,
	}

	// 更新玩家位置
	player.Position = newPosition
	player.TurnsPlayed++

	// 繞圈獎勵
	if circuitCompleted {
		bonus := DefaultCircuitBonus
		player.Company.Capital += bonus.Capital
		player.Company.Employees += bonus.Employees
		result.CapitalChange += bonus.Capital
		result.EmployeeChange += bonus.Employees
	}

	// 格子著陸獎勵 - 更新 Capital 和 Employees (Requirements 2.3, 2.4)
	cell := game.Board.Cells[newPosition]
	player.Company.Capital += cell.BaseCapital
	player.Company.Employees += cell.BaseEmployees
	result.CapitalChange += cell.BaseCapital
	result.EmployeeChange += cell.BaseEmployees
	result.CellType = string(cell.Type)

	// 檢查是否為事件格子 (機會、命運、關卡)
	// Requirements 3.1, 3.2, 3.3: 觸發對應事件
	switch cell.Type {
	case "opportunity", "fate", "challenge":
		result.DecisionRequired = true
		// 不切換玩家，等待玩家做出決策
		return result
	}

	// 切換到下一個玩家
	e.advanceToNextPlayer(game)

	return result
}

// advanceToNextPlayer 切換到下一個玩家
func (e *GameEngine) advanceToNextPlayer(game *Game) {
	currentIndex := -1
	for i, p := range game.Players {
		if p.PlayerID == game.CurrentPlayerID {
			currentIndex = i
			break
		}
	}

	if currentIndex >= 0 {
		nextIndex := (currentIndex + 1) % len(game.Players)
		game.CurrentPlayerID = game.Players[nextIndex].PlayerID

		// 如果回到第一個玩家，增加回合數
		if nextIndex == 0 {
			game.CurrentTurn++
		}
	}
}

// GetPlayer 取得玩家狀態
func (e *GameEngine) GetPlayer(gameID string, playerID string) (*PlayerState, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, ok := e.games[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}

	for _, p := range game.Players {
		if p.PlayerID == playerID {
			playerCopy := *p
			if p.Company != nil {
				companyCopy := *p.Company
				playerCopy.Company = &companyCopy
			}
			return &playerCopy, nil
		}
	}

	return nil, ErrPlayerNotFound
}

// GetGameSummary 取得遊戲結束摘要
func (e *GameEngine) GetGameSummary(gameID string) (*GameSummary, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, ok := e.games[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}

	// 遊戲必須已結束才能產生摘要
	if game.Status != StatusFinished {
		return nil, ErrGameNotFinished
	}

	summary := &GameSummary{
		GameID:           game.ID,
		TotalTurns:       game.CurrentTurn,
		PlayerSummaries:  make([]PlayerSummary, len(game.Players)),
		LearningInsights: []string{},
	}

	// 找出贏家 (資本最高者)
	var winner *PlayerState
	var maxCapital int64 = -1

	for i, player := range game.Players {
		// 建立玩家摘要
		playerSummary := PlayerSummary{
			PlayerID:        player.PlayerID,
			PlayerName:      player.PlayerName,
			CompanyType:     player.Company.Type,
			FinalCapital:    player.Company.Capital,
			FinalEmployees:  player.Company.Employees,
			TurnsPlayed:     player.TurnsPlayed,
			DecisionCount:   len(player.DecisionHistory),
			DecisionRecords: make([]DecisionRecordSummary, len(player.DecisionHistory)),
		}

		// 複製決策記錄
		for j, record := range player.DecisionHistory {
			playerSummary.DecisionRecords[j] = DecisionRecordSummary{
				TurnNumber: record.TurnNumber,
				EventID:    record.EventID,
				ChoiceID:   record.ChoiceID,
				Timestamp:  record.Timestamp,
			}
		}

		summary.PlayerSummaries[i] = playerSummary

		// 追蹤最高資本
		if player.Company.Capital > maxCapital {
			maxCapital = player.Company.Capital
			winner = player
		}
	}

	// 設定贏家
	if winner != nil {
		winnerSummary := PlayerSummary{
			PlayerID:       winner.PlayerID,
			PlayerName:     winner.PlayerName,
			CompanyType:    winner.Company.Type,
			FinalCapital:   winner.Company.Capital,
			FinalEmployees: winner.Company.Employees,
			TurnsPlayed:    winner.TurnsPlayed,
			DecisionCount:  len(winner.DecisionHistory),
		}
		summary.Winner = &winnerSummary
	}

	// 產生學習洞察
	summary.LearningInsights = e.generateLearningInsights(game)

	return summary, nil
}

// generateLearningInsights 產生學習洞察
func (e *GameEngine) generateLearningInsights(game *Game) []string {
	insights := []string{}

	// 基本遊戲統計洞察
	insights = append(insights, fmt.Sprintf("遊戲共進行了 %d 回合", game.CurrentTurn))

	// 分析各玩家表現
	totalDecisions := 0
	for _, player := range game.Players {
		totalDecisions += len(player.DecisionHistory)
	}
	insights = append(insights, fmt.Sprintf("所有玩家共做出 %d 個決策", totalDecisions))

	// 根據公司類型提供洞察
	companyTypeCount := make(map[company.CompanyType]int)
	for _, player := range game.Players {
		companyTypeCount[player.Company.Type]++
	}

	for compType, count := range companyTypeCount {
		switch compType {
		case company.Startup:
			if count > 0 {
				insights = append(insights, "新創公司在雲端採用上通常更具彈性，但需要注意成本控制")
			}
		case company.Traditional:
			if count > 0 {
				insights = append(insights, "傳統企業的雲端轉型需要循序漸進，建議採用混合雲策略")
			}
		case company.CloudReseller:
			if count > 0 {
				insights = append(insights, "雲端代理商應深入了解各種 AWS 服務，以便為客戶提供最佳建議")
			}
		case company.CloudNative:
			if count > 0 {
				insights = append(insights, "雲端原生公司應持續優化架構，善用 AWS 最新服務和功能")
			}
		}
	}

	// AWS SAA 考試相關洞察
	insights = append(insights, "AWS SAA 考試重點：了解各服務的使用場景、成本優化策略、高可用性設計")
	insights = append(insights, "建議複習：EC2 執行個體類型、S3 儲存類別、RDS 多可用區部署、VPC 網路設計")

	return insights
}

// EndGame 結束遊戲
func (e *GameEngine) EndGame(gameID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, ok := e.games[gameID]
	if !ok {
		return ErrGameNotFound
	}

	if game.Status != StatusInProgress {
		return ErrGameNotStarted
	}

	game.Status = StatusFinished
	game.UpdatedAt = time.Now()

	return nil
}

// GetPlayerProgress 取得玩家進度
func (e *GameEngine) GetPlayerProgress(gameID string, playerID string) (*PlayerProgress, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, ok := e.games[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}

	// 找到玩家
	var player *PlayerState
	for _, p := range game.Players {
		if p.PlayerID == playerID {
			player = p
			break
		}
	}
	if player == nil {
		return nil, ErrPlayerNotFound
	}

	progress := &PlayerProgress{
		PlayerID:            playerID,
		TotalDecisions:      len(player.DecisionHistory),
		SuccessfulDecisions: 0,
		FailedDecisions:     0,
		AWSServicesUsed:     []string{},
		TopicsLearned:       []string{},
		SkillAreas:          make(map[string]int),
		Recommendations:     []string{},
	}

	// 分析詳細決策歷史
	awsServicesSet := make(map[string]bool)
	topicsSet := make(map[string]bool)

	// 初始化技能領域
	skillAreas := []string{"compute", "storage", "database", "networking", "security", "cost_optimization"}
	for _, area := range skillAreas {
		progress.SkillAreas[area] = 0
	}

	// 分析 DecisionOutcomeHistory 以獲取詳細進度
	for _, outcome := range player.DecisionOutcomeHistory {
		if outcome.Success {
			progress.SuccessfulDecisions++
		} else {
			progress.FailedDecisions++
		}

		// 收集使用的 AWS 服務
		for _, service := range outcome.AWSServices {
			awsServicesSet[service] = true
			// 根據服務類型增加技能分數
			e.updateSkillAreaFromService(progress.SkillAreas, service)
		}

		// 收集學習要點
		for _, point := range outcome.LearningPoints {
			topicsSet[point] = true
		}
	}

	// 如果沒有詳細記錄，使用基本決策數量估算
	if len(player.DecisionOutcomeHistory) == 0 {
		progress.SuccessfulDecisions = len(player.DecisionHistory)
	}

	// 計算學習進度 (基於決策數量、成功率和公司成長)
	baseProgress := float64(progress.TotalDecisions) * 5.0 // 每個決策 5%
	if baseProgress > 100 {
		baseProgress = 100
	}

	// 成功率加成
	if progress.TotalDecisions > 0 {
		successRate := float64(progress.SuccessfulDecisions) / float64(progress.TotalDecisions)
		baseProgress += successRate * 10.0 // 成功率加成最多 10%
	}

	// 根據公司屬性調整進度
	if player.Company != nil {
		// 雲端採用率影響學習進度
		cloudBonus := player.Company.CloudAdoption * 0.2
		// 資安等級影響學習進度
		securityBonus := float64(player.Company.SecurityLevel) * 2.0

		progress.LearningProgress = baseProgress + cloudBonus + securityBonus
		if progress.LearningProgress > 100 {
			progress.LearningProgress = 100
		}

		// 根據公司類型設定基礎技能領域分數
		switch player.Company.Type {
		case company.Startup:
			progress.SkillAreas["compute"] += 60
			progress.SkillAreas["cost_optimization"] += 70
			progress.SkillAreas["storage"] += 50
		case company.Traditional:
			progress.SkillAreas["security"] += 60
			progress.SkillAreas["networking"] += 55
			progress.SkillAreas["database"] += 50
		case company.CloudReseller:
			progress.SkillAreas["compute"] += 70
			progress.SkillAreas["storage"] += 65
			progress.SkillAreas["networking"] += 60
		case company.CloudNative:
			progress.SkillAreas["compute"] += 80
			progress.SkillAreas["storage"] += 75
			progress.SkillAreas["database"] += 70
			progress.SkillAreas["security"] += 65
		}

		// 根據雲端採用率調整技能分數
		adoptionMultiplier := 1.0 + (player.Company.CloudAdoption / 200.0)
		for area := range progress.SkillAreas {
			progress.SkillAreas[area] = int(float64(progress.SkillAreas[area]) * adoptionMultiplier)
			if progress.SkillAreas[area] > 100 {
				progress.SkillAreas[area] = 100
			}
		}
	}

	// 收集已使用的 AWS 服務
	for service := range awsServicesSet {
		progress.AWSServicesUsed = append(progress.AWSServicesUsed, service)
	}

	// 收集已學習的主題
	for topic := range topicsSet {
		progress.TopicsLearned = append(progress.TopicsLearned, topic)
	}

	// 產生改進建議
	progress.Recommendations = e.generateProgressRecommendations(progress, player)

	return progress, nil
}

// updateSkillAreaFromService 根據 AWS 服務更新技能領域分數
func (e *GameEngine) updateSkillAreaFromService(skillAreas map[string]int, service string) {
	// 根據服務名稱判斷技能領域
	serviceSkillMap := map[string]string{
		"EC2":         "compute",
		"Lambda":      "compute",
		"ECS":         "compute",
		"EKS":         "compute",
		"Fargate":     "compute",
		"S3":          "storage",
		"EBS":         "storage",
		"EFS":         "storage",
		"Glacier":     "storage",
		"RDS":         "database",
		"DynamoDB":    "database",
		"Aurora":      "database",
		"ElastiCache": "database",
		"Redshift":    "database",
		"VPC":         "networking",
		"Route53":     "networking",
		"CloudFront":  "networking",
		"ELB":         "networking",
		"ALB":         "networking",
		"NLB":         "networking",
		"IAM":         "security",
		"KMS":         "security",
		"WAF":         "security",
		"Shield":      "security",
		"GuardDuty":   "security",
		"Inspector":   "security",
	}

	if area, ok := serviceSkillMap[service]; ok {
		skillAreas[area] += 5 // 每使用一個服務增加 5 分
		if skillAreas[area] > 100 {
			skillAreas[area] = 100
		}
	}
}

// generateProgressRecommendations 產生進度改進建議
func (e *GameEngine) generateProgressRecommendations(progress *PlayerProgress, player *PlayerState) []string {
	recommendations := []string{}

	// 根據技能領域分數給出建議
	for area, score := range progress.SkillAreas {
		if score < 50 {
			switch area {
			case "compute":
				recommendations = append(recommendations, "建議加強 EC2、Lambda 等運算服務的學習")
			case "storage":
				recommendations = append(recommendations, "建議加強 S3、EBS、EFS 等儲存服務的學習")
			case "database":
				recommendations = append(recommendations, "建議加強 RDS、DynamoDB、Aurora 等資料庫服務的學習")
			case "networking":
				recommendations = append(recommendations, "建議加強 VPC、Route 53、CloudFront 等網路服務的學習")
			case "security":
				recommendations = append(recommendations, "建議加強 IAM、KMS、WAF 等安全服務的學習")
			case "cost_optimization":
				recommendations = append(recommendations, "建議學習 AWS 成本優化策略，如 Reserved Instances 和 Savings Plans")
			}
		}
	}

	// 根據決策數量給出建議
	if progress.TotalDecisions < 5 {
		recommendations = append(recommendations, "建議多參與遊戲決策，累積更多實戰經驗")
	}

	// 根據學習進度給出建議
	if progress.LearningProgress < 30 {
		recommendations = append(recommendations, "學習進度較低，建議持續遊玩並關注每次決策的回饋")
	} else if progress.LearningProgress < 60 {
		recommendations = append(recommendations, "學習進度良好，建議嘗試更具挑戰性的決策")
	} else {
		recommendations = append(recommendations, "學習進度優秀！建議開始準備 AWS SAA 認證考試")
	}

	// 根據公司類型給出特定建議
	if player.Company != nil {
		switch player.Company.Type {
		case company.Startup:
			recommendations = append(recommendations, "作為新創公司，建議關注成本效益和快速部署的 AWS 服務")
		case company.Traditional:
			recommendations = append(recommendations, "作為傳統企業，建議學習混合雲架構和漸進式遷移策略")
		case company.CloudReseller:
			recommendations = append(recommendations, "作為雲端代理商，建議深入了解各種 AWS 服務的最佳實踐")
		case company.CloudNative:
			recommendations = append(recommendations, "作為雲端原生公司，建議探索 AWS 最新服務和無伺服器架構")
		}
	}

	return recommendations
}

// RecordDecision 記錄決策
func (e *GameEngine) RecordDecision(gameID string, playerID string, record DecisionOutcomeRecord) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, ok := e.games[gameID]
	if !ok {
		return ErrGameNotFound
	}

	// 找到玩家
	var player *PlayerState
	for _, p := range game.Players {
		if p.PlayerID == playerID {
			player = p
			break
		}
	}
	if player == nil {
		return ErrPlayerNotFound
	}

	// 記錄基本決策
	player.DecisionHistory = append(player.DecisionHistory, record.DecisionRecord)
	// 記錄詳細決策結果 (用於進度追蹤)
	player.DecisionOutcomeHistory = append(player.DecisionOutcomeHistory, record)
	game.UpdatedAt = time.Now()

	// 決策完成後，切換到下一個玩家
	e.advanceToNextPlayer(game)

	return nil
}

// GetCurrentPlayer 取得當前回合玩家
// Requirements 7.2, 7.3: 管理回合順序
func (e *GameEngine) GetCurrentPlayer(gameID string) (*PlayerState, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, ok := e.games[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}

	if game.Status != StatusInProgress {
		return nil, ErrGameNotStarted
	}

	for _, p := range game.Players {
		if p.PlayerID == game.CurrentPlayerID {
			playerCopy := *p
			if p.Company != nil {
				companyCopy := *p.Company
				playerCopy.Company = &companyCopy
			}
			return &playerCopy, nil
		}
	}

	return nil, ErrPlayerNotFound
}

// GetTurnOrder 取得回合順序
// Requirements 7.2: 公平管理回合順序
func (e *GameEngine) GetTurnOrder(gameID string) ([]string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, ok := e.games[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}

	order := make([]string, len(game.Players))
	for i, p := range game.Players {
		order[i] = p.PlayerID
	}

	return order, nil
}

// GetPlayerTurnCount 取得各玩家已執行的回合數
// Requirements 7.2, 7.3: 確保回合公平性
func (e *GameEngine) GetPlayerTurnCount(gameID string) (map[string]int, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, ok := e.games[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}

	turnCounts := make(map[string]int)
	for _, p := range game.Players {
		turnCounts[p.PlayerID] = p.TurnsPlayed
	}

	return turnCounts, nil
}

// DetermineWinner 判定贏家
// Requirements 7.4: 基於最終資本判定贏家
func (e *GameEngine) DetermineWinner(gameID string) (*PlayerState, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, ok := e.games[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}

	if game.Status != StatusFinished {
		return nil, ErrGameNotFinished
	}

	if len(game.Players) == 0 {
		return nil, ErrPlayerNotFound
	}

	// 找出資本最高的玩家
	var winner *PlayerState
	var maxCapital int64 = -1

	for _, p := range game.Players {
		if p.Company != nil && p.Company.Capital > maxCapital {
			maxCapital = p.Company.Capital
			winner = p
		}
	}

	if winner == nil {
		return nil, ErrPlayerNotFound
	}

	// 返回副本
	winnerCopy := *winner
	if winner.Company != nil {
		companyCopy := *winner.Company
		winnerCopy.Company = &companyCopy
	}

	return &winnerCopy, nil
}

// SaveGame 儲存遊戲狀態
// Requirements 8.1: 儲存遊戲狀態
// Requirements 8.4: 使用 JSON 格式編碼
func (e *GameEngine) SaveGame(gameID string) ([]byte, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, ok := e.games[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}

	// 序列化遊戲狀態
	data, err := SerializeGameState(game)
	if err != nil {
		return nil, ErrSaveFailed
	}

	return data, nil
}

// LoadGame 載入遊戲狀態
// Requirements 8.2: 載入遊戲狀態
// Requirements 8.5: 驗證並還原完整的遊戲狀態
func (e *GameEngine) LoadGame(data []byte) (*Game, error) {
	// 反序列化遊戲狀態
	game, err := DeserializeGameState(data)
	if err != nil {
		return nil, ErrLoadFailed
	}

	// 將遊戲加入引擎
	e.mu.Lock()
	defer e.mu.Unlock()

	e.games[game.ID] = game

	return game, nil
}

// ImportGame 匯入遊戲到引擎
func (e *GameEngine) ImportGame(game *Game) error {
	if game == nil {
		return ErrGameNotFound
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.games[game.ID] = game
	return nil
}
