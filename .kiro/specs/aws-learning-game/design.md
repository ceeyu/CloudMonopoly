# Design Document: AWS Learning Game

## Overview

AWS Learning Game 是一款使用 Golang 開發、部署於 AWS 的大富翁風格學習遊戲。系統採用前後端分離架構，後端提供 REST API，前端為 Web 應用。遊戲核心邏輯在後端處理，確保遊戲狀態一致性和防作弊。

### 技術選型
- **後端**: Golang (Gin/Echo framework)
- **前端**: React/Vue.js (SPA)
- **資料庫**: AWS DynamoDB (遊戲狀態) + S3 (靜態資源)
- **部署**: AWS Lambda + API Gateway 或 ECS/Fargate
- **即時通訊**: WebSocket (多人遊戲同步)

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Client Layer                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │   Web UI    │  │  Mobile UI  │  │   API CLI   │              │
│  │  (React)    │  │  (Future)   │  │  (Testing)  │              │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘              │
└─────────┼────────────────┼────────────────┼─────────────────────┘
          │                │                │
          ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway Layer                           │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │              AWS API Gateway / ALB                       │    │
│  │         (REST API + WebSocket endpoints)                 │    │
│  └─────────────────────────┬───────────────────────────────┘    │
└─────────────────────────────┼───────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Application Layer                            │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                  Golang Backend Service                  │    │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌──────────┐│    │
│  │  │   Game    │ │  Company  │ │   Event   │ │ Decision ││    │
│  │  │  Engine   │ │  Manager  │ │  System   │ │  Engine  ││    │
│  │  └───────────┘ └───────────┘ └───────────┘ └──────────┘│    │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐            │    │
│  │  │   Board   │ │   AWS     │ │Architecture│            │    │
│  │  │  Manager  │ │  Catalog  │ │ Visualizer│            │    │
│  │  └───────────┘ └───────────┘ └───────────┘            │    │
│  └─────────────────────────┬───────────────────────────────┘    │
└─────────────────────────────┼───────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Data Layer                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │  DynamoDB   │  │     S3      │  │  ElastiCache│              │
│  │ (Game State)│  │  (Assets)   │  │  (Sessions) │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### 1. Game Engine (遊戲引擎)

負責遊戲核心流程控制。

```go
// GameEngine 遊戲引擎介面
type GameEngine interface {
    // 建立新遊戲
    CreateGame(config GameConfig) (*Game, error)
    // 加入遊戲
    JoinGame(gameID string, player Player) error
    // 開始遊戲
    StartGame(gameID string) error
    // 執行回合
    ExecuteTurn(gameID string, playerID string, action TurnAction) (*TurnResult, error)
    // 取得遊戲狀態
    GetGameState(gameID string) (*GameState, error)
    // 儲存遊戲
    SaveGame(gameID string) error
    // 載入遊戲
    LoadGame(gameID string) (*Game, error)
}

type GameConfig struct {
    MaxPlayers    int
    BoardType     string
    DifficultyLevel string
}

type TurnAction struct {
    ActionType string      // "roll_dice", "make_decision", "use_item"
    Payload    interface{}
}

type TurnResult struct {
    DiceValue      int
    NewPosition    int
    CapitalChange  int64
    EmployeeChange int
    Event          *Event
    DecisionRequired bool
}
```

### 2. Company Manager (公司管理器)

管理玩家公司的所有屬性和狀態變化。

```go
// CompanyManager 公司管理介面
type CompanyManager interface {
    // 建立公司
    CreateCompany(companyType CompanyType) (*Company, error)
    // 更新資本
    UpdateCapital(companyID string, delta int64) error
    // 更新員工數
    UpdateEmployees(companyID string, delta int) error
    // 取得公司狀態
    GetCompanyState(companyID string) (*Company, error)
    // 檢查是否符合決策條件
    CheckDecisionEligibility(companyID string, decision Decision) (bool, []string)
}

type CompanyType string

const (
    Startup       CompanyType = "startup"        // 新創公司
    Traditional   CompanyType = "traditional"    // 傳產公司
    CloudReseller CompanyType = "cloud_reseller" // 雲端代理商
    CloudNative   CompanyType = "cloud_native"   // 雲端公司
)

type Company struct {
    ID                string
    Name              string
    Type              CompanyType
    Capital           int64       // 資本額 (萬元)
    Employees         int         // 員工人數
    IsInternational   bool        // 是否跨國企業
    ProductCycle      string      // 產品週期: "development", "launch", "growth", "mature"
    TechDebt          int         // 技術債 (影響決策)
    SecurityLevel     int         // 資安等級 1-5
    CloudAdoption     float64     // 雲端採用率 0-100%
    Infrastructure    []string    // 已部署的基礎設施
}

// 各公司類型初始屬性
var CompanyDefaults = map[CompanyType]Company{
    Startup: {
        Capital: 500, Employees: 10, SecurityLevel: 2, CloudAdoption: 30,
    },
    Traditional: {
        Capital: 5000, Employees: 200, SecurityLevel: 3, CloudAdoption: 10,
    },
    CloudReseller: {
        Capital: 2000, Employees: 50, SecurityLevel: 4, CloudAdoption: 80,
    },
    CloudNative: {
        Capital: 1000, Employees: 30, SecurityLevel: 4, CloudAdoption: 95,
    },
}
```

### 3. Board Manager (棋盤管理器)

管理遊戲棋盤和格子。

```go
// BoardManager 棋盤管理介面
type BoardManager interface {
    // 建立棋盤
    CreateBoard(boardType string) (*Board, error)
    // 取得格子資訊
    GetCell(position int) (*Cell, error)
    // 計算新位置
    CalculateNewPosition(current int, diceValue int) int
    // 取得格子事件
    GetCellEvent(position int) (*Event, error)
}

type Board struct {
    ID       string
    Name     string
    Cells    []Cell
    Size     int
}

type CellType string

const (
    CellNormal     CellType = "normal"      // 一般格
    CellOpportunity CellType = "opportunity" // 機會
    CellFate       CellType = "fate"        // 命運
    CellChallenge  CellType = "challenge"   // 關卡
    CellStart      CellType = "start"       // 起點
    CellBonus      CellType = "bonus"       // 獎勵格
)

type Cell struct {
    Position       int
    Type           CellType
    Name           string
    BaseCapital    int64  // 基礎資本獎勵
    BaseEmployees  int    // 基礎員工獎勵
    EventID        string // 關聯事件ID (可選)
}
```

### 4. Event System (事件系統)

處理各種遊戲事件，基於真實商業案例。

```go
// EventSystem 事件系統介面
type EventSystem interface {
    // 取得事件
    GetEvent(eventID string) (*Event, error)
    // 隨機取得事件 (依類型)
    GetRandomEvent(eventType EventType, companyContext *Company) (*Event, error)
    // 處理事件結果
    ProcessEventOutcome(event *Event, choice int, company *Company) (*EventOutcome, error)
}

type EventType string

const (
    EventOpportunity EventType = "opportunity" // 機會事件
    EventFate        EventType = "fate"        // 命運事件
    EventChallenge   EventType = "challenge"   // 關卡事件
    EventSecurity    EventType = "security"    // 資安事件
)

type Event struct {
    ID              string
    Type            EventType
    Title           string
    Description     string        // 事件描述 (真實情境)
    RealWorldRef    string        // 真實案例參考
    Context         EventContext  // 事件背景
    Choices         []EventChoice // 可選方案
    AWSTopics       []string      // 相關 AWS SAA 考點
}

type EventContext struct {
    Scenario        string   // 情境說明
    BusinessImpact  string   // 商業影響
    TechnicalNeeds  []string // 技術需求
    Constraints     []string // 限制條件
}

type EventChoice struct {
    ID              int
    Title           string
    Description     string
    IsAWS           bool            // 是否為 AWS 方案
    AWSServices     []string        // 使用的 AWS 服務
    OnPremSolution  string          // 地端方案描述
    Requirements    ChoiceRequirements
    Outcomes        ChoiceOutcomes
    ArchitectureDiagram string      // 架構圖 (ASCII/Mermaid)
}

type ChoiceRequirements struct {
    MinCapital      int64
    MinEmployees    int
    MinSecurityLevel int
    RequiredInfra   []string
}

type ChoiceOutcomes struct {
    CapitalChange   int64
    EmployeeChange  int
    SecurityChange  int
    CloudAdoptionChange float64
    SuccessRate     float64  // 成功機率 (基於公司屬性調整)
    TimeToImplement int      // 實施時間 (回合數)
}

type EventOutcome struct {
    Success         bool
    Message         string
    CapitalDelta    int64
    EmployeeDelta   int
    LearningPoints  []string  // 學習要點
    AWSBestPractice string    // AWS 最佳實踐說明
}
```

### 5. Decision Engine (決策引擎)

計算決策結果，參考真實 AWS 定價和商業邏輯。

```go
// DecisionEngine 決策引擎介面
type DecisionEngine interface {
    // 評估決策
    EvaluateDecision(company *Company, choice *EventChoice) (*DecisionEvaluation, error)
    // 執行決策
    ExecuteDecision(company *Company, choice *EventChoice) (*DecisionResult, error)
    // 取得方案比較
    GetComparison(choices []EventChoice, company *Company) (*Comparison, error)
}

type DecisionEvaluation struct {
    IsEligible      bool
    EligibilityIssues []string
    RiskLevel       string    // "low", "medium", "high"
    ExpectedROI     float64
    ImplementationTime int
    Recommendation  string
}

type DecisionResult struct {
    Success         bool
    ActualOutcome   ChoiceOutcomes
    Explanation     string
    Penalties       []Penalty
    Rewards         []Reward
}

type Penalty struct {
    Type        string  // "budget_overrun", "delay", "security_breach"
    Description string
    Impact      int64
}

type Reward struct {
    Type        string  // "efficiency", "growth", "innovation"
    Description string
    Bonus       int64
}

type Comparison struct {
    Choices         []ChoiceComparison
    Recommendation  int  // 推薦選項 index
    ReasoningSteps  []string
}

type ChoiceComparison struct {
    ChoiceID        int
    CostAnalysis    CostAnalysis
    ScalabilityScore int
    ComplexityScore int
    SecurityScore   int
    AWSExamRelevance []string
}

type CostAnalysis struct {
    InitialCost     int64
    MonthlyCost     int64
    ThreeYearTCO    int64
    CostBreakdown   map[string]int64
}
```

### 6. AWS Service Catalog (AWS 服務目錄)

儲存 AWS 服務資訊，用於決策和學習。

```go
// AWSServiceCatalog AWS 服務目錄介面
type AWSServiceCatalog interface {
    // 取得服務
    GetService(serviceID string) (*AWSService, error)
    // 依類別取得服務
    GetServicesByCategory(category string) ([]AWSService, error)
    // 取得適合的服務建議
    GetRecommendedServices(scenario string, company *Company) ([]AWSService, error)
}

type AWSService struct {
    ID              string
    Name            string
    Category        string    // "compute", "storage", "database", "networking", "security"
    Description     string
    UseCases        []string
    PricingModel    PricingModel
    SAAExamTopics   []string  // SAA 考試相關主題
    BestPractices   []string
}

type PricingModel struct {
    Type            string    // "on_demand", "reserved", "spot", "savings_plan"
    BasePrice       float64
    Unit            string    // "hour", "GB", "request"
    FreeTeir        string
}
```

### 7. Architecture Visualizer (架構視覺化)

產生架構圖和比較表。

```go
// ArchitectureVisualizer 架構視覺化介面
type ArchitectureVisualizer interface {
    // 產生架構圖
    GenerateDiagram(choice *EventChoice) (string, error)
    // 產生比較表
    GenerateComparisonTable(choices []EventChoice) (string, error)
    // 產生公司當前架構
    GenerateCompanyArchitecture(company *Company) (string, error)
}
```

## Data Models

### Game State (遊戲狀態)

```go
type GameState struct {
    GameID          string
    Status          string    // "waiting", "in_progress", "finished"
    CurrentTurn     int
    CurrentPlayerID string
    Players         []PlayerState
    Board           *Board
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type PlayerState struct {
    PlayerID        string
    PlayerName      string
    Company         *Company
    Position        int
    TurnsPlayed     int
    DecisionHistory []DecisionRecord
}

type DecisionRecord struct {
    TurnNumber      int
    EventID         string
    ChoiceID        int
    Outcome         *DecisionResult
    Timestamp       time.Time
}
```

### DynamoDB Schema

```
Table: aws-learning-game-sessions
- PK: GAME#{gameID}
- SK: STATE
- Attributes: GameState JSON

Table: aws-learning-game-events  
- PK: EVENT#{eventID}
- SK: TYPE#{eventType}
- Attributes: Event JSON

Table: aws-learning-game-services
- PK: SERVICE#{serviceID}
- SK: CATEGORY#{category}
- Attributes: AWSService JSON
```



## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Company Initialization Completeness
*For any* valid company type selection, the created Company SHALL have all predefined attributes (Capital, Employees, SecurityLevel, CloudAdoption) set to non-zero default values matching the CompanyDefaults configuration.

**Validates: Requirements 1.2**

### Property 2: Position Movement Correctness
*For any* dice roll value D and current position P on a board of size S, the new position SHALL equal (P + D) mod S, ensuring players wrap around the board correctly.

**Validates: Requirements 2.2**

### Property 3: Cell Landing Resource Update
*For any* cell landing event, the Company's Capital SHALL increase by exactly the cell's BaseCapital value AND the Company's Employees SHALL increase by exactly the cell's BaseEmployees value.

**Validates: Requirements 2.3, 2.4**

### Property 4: Board Size Constraint
*For any* created Board, the number of cells SHALL be at least 30.

**Validates: Requirements 2.5**

### Property 5: Circuit Completion Bonus
*For any* player movement that causes position to wrap around (cross the start), bonus Capital and Employees SHALL be added to the Company.

**Validates: Requirements 2.6**

### Property 6: Event Type Cell Matching
*For any* cell landing, the returned Event type SHALL match the cell type: opportunity cells return opportunity events, fate cells return fate events, challenge cells return challenge events.

**Validates: Requirements 3.1, 3.2, 3.3**

### Property 7: Security Event Mitigation Choices
*For any* security incident event, the Event SHALL contain at least 2 mitigation strategy choices.

**Validates: Requirements 3.4**

### Property 8: Event Display Completeness
*For any* Event presentation, the output SHALL include: Title, Description, Context (Scenario, BusinessImpact), and at least one Choice.

**Validates: Requirements 3.6**

### Property 9: Infrastructure Decision Options
*For any* infrastructure decision event, the available choices SHALL include at least one on-premises option AND at least one AWS cloud option.

**Validates: Requirements 4.1**

### Property 10: AWS Service Category Relevance
*For any* scenario-based service query, all returned AWSService items SHALL have a Category matching the scenario's technical requirements.

**Validates: Requirements 4.2**

### Property 11: Decision Outcome Calculation
*For any* decision execution, the DecisionResult SHALL contain CapitalDelta, EmployeeDelta, and Success status calculated based on Company attributes and ChoiceRequirements.

**Validates: Requirements 4.3**

### Property 12: Decision Evaluation Completeness
*For any* decision evaluation, the DecisionEvaluation SHALL include assessments for: cost (ExpectedROI), time (ImplementationTime), and risk (RiskLevel).

**Validates: Requirements 4.4**

### Property 13: Incompatible Decision Penalty
*For any* decision where Company attributes do not meet ChoiceRequirements (e.g., Capital < MinCapital), the DecisionResult SHALL contain at least one Penalty.

**Validates: Requirements 4.5**

### Property 14: Architecture Diagram Generation
*For any* EventChoice, the ArchitectureVisualizer SHALL produce a non-empty diagram string.

**Validates: Requirements 5.1**

### Property 15: Comparison Table Generation
*For any* set of 2 or more EventChoices, the generated comparison SHALL include CostAnalysis, ScalabilityScore, ComplexityScore, and SecurityScore for each choice.

**Validates: Requirements 5.2**

### Property 16: Decision Feedback Completeness
*For any* DecisionResult, the feedback SHALL include an Explanation string AND at least one AWSBestPractice or LearningPoint.

**Validates: Requirements 6.1, 6.2**

### Property 17: Game End Summary
*For any* game that reaches "finished" status, a summary SHALL be generatable containing all DecisionRecords from all players.

**Validates: Requirements 6.3**

### Property 18: Player Count Validation
*For any* game creation or join attempt, the Game_Engine SHALL enforce player count between 2 and 4 inclusive.

**Validates: Requirements 7.1**

### Property 19: Turn Order Fairness
*For any* multi-player game with N players, after N consecutive turn advancements, each player SHALL have had exactly one turn.

**Validates: Requirements 7.2, 7.3**

### Property 20: Winner Determination
*For any* finished game, exactly one player SHALL be determined as winner based on highest final Capital value.

**Validates: Requirements 7.4**

### Property 21: Game State Serialization Round-Trip
*For any* valid GameState, serializing to JSON then deserializing SHALL produce a GameState equivalent to the original (all fields match).

**Validates: Requirements 8.4, 8.5**

### Property 22: Save and Load Consistency
*For any* saved game, loading it SHALL restore a GameState where all player positions, company attributes, and turn number match the saved state.

**Validates: Requirements 8.1, 8.2**

## Error Handling

### Game Engine Errors

| Error Code | Condition | Response |
|------------|-----------|----------|
| `GAME_NOT_FOUND` | 遊戲 ID 不存在 | 返回 404，提示重新建立遊戲 |
| `GAME_FULL` | 遊戲已達最大玩家數 | 返回 400，提示加入其他遊戲 |
| `NOT_YOUR_TURN` | 非該玩家回合 | 返回 403，提示等待 |
| `INVALID_ACTION` | 無效的回合動作 | 返回 400，列出有效動作 |
| `GAME_ALREADY_STARTED` | 遊戲已開始 | 返回 400，無法加入 |

### Decision Engine Errors

| Error Code | Condition | Response |
|------------|-----------|----------|
| `INSUFFICIENT_CAPITAL` | 資本不足 | 返回決策結果但標記為失敗，套用懲罰 |
| `MISSING_PREREQUISITES` | 缺少前置條件 | 返回不符合條件清單 |
| `INVALID_CHOICE` | 選項不存在 | 返回 400，列出有效選項 |

### Persistence Errors

| Error Code | Condition | Response |
|------------|-----------|----------|
| `SAVE_FAILED` | 儲存失敗 | 返回 500，建議重試 |
| `LOAD_FAILED` | 載入失敗 | 返回 404 或 500，提示原因 |
| `CORRUPTED_STATE` | 狀態損壞 | 返回 500，建議重新開始 |

## Testing Strategy

### Unit Tests
- 測試各元件的獨立功能
- 測試邊界條件（如棋盤邊界、資本為零）
- 測試錯誤處理路徑

### Property-Based Tests
使用 `gopter` 或 `rapid` 進行屬性測試：

- **Property 1-4**: 公司和棋盤初始化測試
- **Property 5-8**: 事件系統測試
- **Property 9-16**: 決策引擎測試
- **Property 17-20**: 多人遊戲邏輯測試
- **Property 21-22**: 序列化/反序列化 round-trip 測試

每個屬性測試至少執行 100 次迭代。

### Integration Tests
- API 端點整合測試
- 完整遊戲流程測試
- 多玩家同步測試

### Test Configuration
```go
// 使用 rapid 進行屬性測試
import "pgregory.net/rapid"

func TestProperty21_GameStateRoundTrip(t *testing.T) {
    rapid.Check(t, func(t *rapid.T) {
        // Feature: aws-learning-game, Property 21: Game State Serialization Round-Trip
        state := generateRandomGameState(t)
        json, err := SerializeGameState(state)
        require.NoError(t, err)
        restored, err := DeserializeGameState(json)
        require.NoError(t, err)
        assert.Equal(t, state, restored)
    })
}
```
