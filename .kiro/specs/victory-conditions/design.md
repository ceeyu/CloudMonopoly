# Design Document: Victory Conditions

## Overview

本設計實作角色專屬勝利條件機制，包含四種公司類型的不同勝利目標、30 回合限制、勝利進度計算，以及遊戲規則說明頁面。

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Frontend (React)                        │
├─────────────────────────────────────────────────────────────┤
│  GameRulesModal  │  VictoryModal  │  CompanyStatus (進度條)  │
└────────┬─────────┴───────┬────────┴──────────┬──────────────┘
         │                 │                   │
         ▼                 ▼                   ▼
┌─────────────────────────────────────────────────────────────┐
│                      API Layer (Go)                          │
├─────────────────────────────────────────────────────────────┤
│  GET /games/:id (含 victory_progress, winner_id)             │
│  POST /games/:id/turn (檢查勝利條件)                          │
└────────┬─────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│                   Game Engine (Go)                           │
├─────────────────────────────────────────────────────────────┤
│  VictoryChecker  │  ProgressCalculator  │  TurnLimiter       │
└─────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### Backend Components

#### 1. VictoryCondition (新增型別)

```go
// internal/game/victory.go

type VictoryCondition struct {
    CompanyType     company.CompanyType
    TargetCapital   int64   // 目標資本 (0 表示不檢查)
    TargetEmployees int     // 目標員工數 (0 表示不檢查)
    TargetCloudAdoption float64 // 目標雲端採用率 (0 表示不檢查)
    TargetSecurityLevel int    // 目標資安等級 (0 表示不檢查)
}

// 預設勝利條件
var DefaultVictoryConditions = map[company.CompanyType]VictoryCondition{
    company.Startup: {
        CompanyType:   company.Startup,
        TargetCapital: 3000,
    },
    company.Traditional: {
        CompanyType:         company.Traditional,
        TargetCloudAdoption: 80.0,
    },
    company.CloudReseller: {
        CompanyType:     company.CloudReseller,
        TargetEmployees: 150,
    },
    company.CloudNative: {
        CompanyType:         company.CloudNative,
        TargetCapital:       2000,
        TargetSecurityLevel: 5,
    },
}
```

#### 2. VictoryChecker Interface

```go
// CheckVictory 檢查玩家是否達成勝利條件
func (e *GameEngine) CheckVictory(player *PlayerState) bool

// CalculateVictoryProgress 計算勝利進度百分比
func (e *GameEngine) CalculateVictoryProgress(player *PlayerState) float64

// DetermineWinnerByProgress 根據進度判定贏家 (回合限制時使用)
func (e *GameEngine) DetermineWinnerByProgress(game *Game) *PlayerState
```

#### 3. Game 結構擴充

```go
// internal/game/types.go 擴充

type Game struct {
    // ... 現有欄位 ...
    MaxTurnsPerPlayer int      // 每位玩家最大回合數 (預設 30)
    WinnerID          string   // 贏家 ID (遊戲結束時設定)
    WinReason         string   // 勝利原因: "condition_met" 或 "turn_limit"
}

type PlayerState struct {
    // ... 現有欄位 ...
    VictoryProgress float64 // 勝利進度百分比 (0-100)
}
```

### Frontend Components

#### 1. GameRulesModal

```typescript
// web/src/components/GameRulesModal.tsx

interface GameRulesModalProps {
    isOpen: boolean;
    onStart: () => void;
    currentPlayerCompanyType: string;
}

// 顯示內容:
// - 遊戲目標說明
// - 四種公司類型的勝利條件
// - 30 回合限制說明
// - 當前玩家的勝利條件高亮顯示
```

#### 2. VictoryModal

```typescript
// web/src/components/VictoryModal.tsx

interface VictoryModalProps {
    isOpen: boolean;
    winner: PlayerState;
    players: PlayerState[];
    winReason: 'condition_met' | 'turn_limit';
    onClose: () => void;
}
```

#### 3. CompanyStatus 擴充 (進度條)

```typescript
// 在 CompanyStatus 組件中新增勝利進度條
interface VictoryProgressProps {
    progress: number;  // 0-100
    companyType: string;
    targetDescription: string;
}
```

## Data Models

### API Response 擴充

```typescript
// web/src/api/types.ts 擴充

interface GameStateResponse {
    // ... 現有欄位 ...
    max_turns_per_player: number;
    winner_id?: string;
    win_reason?: 'condition_met' | 'turn_limit';
}

interface PlayerState {
    // ... 現有欄位 ...
    victory_progress: number;  // 0-100
}

interface VictoryConditionInfo {
    company_type: string;
    target_description: string;
    current_value: string;
    target_value: string;
    progress: number;
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Startup Victory Detection
*For any* Startup player with capital >= 3000, the CheckVictory function SHALL return true.
**Validates: Requirements 1.1**

### Property 2: Traditional Victory Detection
*For any* Traditional player with cloud_adoption >= 80, the CheckVictory function SHALL return true.
**Validates: Requirements 1.2**

### Property 3: CloudReseller Victory Detection
*For any* CloudReseller player with employees >= 150, the CheckVictory function SHALL return true.
**Validates: Requirements 1.3**

### Property 4: CloudNative Victory Detection
*For any* CloudNative player with capital >= 2000 AND security_level >= 5, the CheckVictory function SHALL return true.
**Validates: Requirements 1.4**

### Property 5: Victory Progress Calculation
*For any* player, the CalculateVictoryProgress function SHALL return a value between 0 and 100, correctly calculated based on company type.
**Validates: Requirements 3.1, 3.2, 3.3, 3.4, 3.5**

### Property 6: Turn Limit Enforcement
*For any* game, no player SHALL be allowed to roll dice more than 30 times.
**Validates: Requirements 2.1**

### Property 7: Winner By Progress
*For any* game that ends due to turn limit, the winner SHALL be the player with the highest victory progress percentage.
**Validates: Requirements 2.3**

## Error Handling

| 錯誤情況 | 處理方式 |
|---------|---------|
| 玩家嘗試在非自己回合擲骰 | 返回 ErrNotYourTurn |
| 玩家已達 30 回合限制 | 返回 ErrTurnLimitReached |
| 遊戲已結束時嘗試操作 | 返回 ErrGameFinished |
| 計算進度時公司類型無效 | 返回 0% 進度 |

## Testing Strategy

### Unit Tests
- 測試各公司類型的勝利條件檢查
- 測試勝利進度計算公式
- 測試回合限制邏輯
- 測試同回合多人達成時的優先順序

### Property-Based Tests
- 使用 Go 的 testing/quick 或 gopter 進行屬性測試
- 生成隨機玩家狀態，驗證勝利檢查的正確性
- 生成隨機進度值，驗證計算公式的正確性

### Integration Tests
- 測試完整遊戲流程直到勝利
- 測試 30 回合限制後的勝利判定
- 測試前端 Modal 的顯示邏輯
