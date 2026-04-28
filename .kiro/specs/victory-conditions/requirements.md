# Requirements Document

## Introduction

本功能為 AWS 學習遊戲新增角色專屬勝利條件機制。每種公司類型有不同的勝利目標，遊戲最多進行 30 回合，先達成目標者獲勝。同時在遊戲開始前顯示規則說明頁面。

## Glossary

- **Victory_Condition_System**: 勝利條件判定系統，負責檢查玩家是否達成勝利目標
- **Game_Rules_Modal**: 遊戲規則彈窗，在遊戲開始前顯示遊戲說明與勝利條件
- **Victory_Progress**: 勝利進度，計算玩家距離勝利目標的完成百分比
- **Max_Turns**: 最大回合數限制，每位玩家最多擲 30 次骰子

## Requirements

### Requirement 1: 角色專屬勝利條件

**User Story:** As a player, I want each company type to have unique victory conditions, so that different roles have different strategies and gameplay experiences.

#### Acceptance Criteria

1. WHEN a Startup player's capital reaches 3000 (萬), THE Victory_Condition_System SHALL declare that player as the winner
2. WHEN a Traditional player's cloud adoption rate reaches 80%, THE Victory_Condition_System SHALL declare that player as the winner
3. WHEN a CloudReseller player's employee count reaches 150, THE Victory_Condition_System SHALL declare that player as the winner
4. WHEN a CloudNative player's capital reaches 2000 (萬) AND security level reaches 5, THE Victory_Condition_System SHALL declare that player as the winner
5. WHEN multiple players achieve their victory conditions in the same turn, THE Victory_Condition_System SHALL declare the player who completed their turn first as the winner

### Requirement 2: 回合限制機制

**User Story:** As a player, I want the game to have a maximum turn limit, so that games don't go on indefinitely.

#### Acceptance Criteria

1. THE Game_Engine SHALL limit each player to a maximum of 30 dice rolls
2. WHEN all players have completed 30 turns without anyone achieving victory, THE Victory_Condition_System SHALL end the game
3. WHEN the game ends due to turn limit, THE Victory_Condition_System SHALL determine the winner based on victory progress percentage
4. THE Victory_Progress SHALL calculate each player's completion percentage toward their victory condition

### Requirement 3: 勝利進度計算

**User Story:** As a player, I want to see my progress toward victory, so that I can track how close I am to winning.

#### Acceptance Criteria

1. FOR Startup players, THE Victory_Progress SHALL calculate as (current_capital / 3000) * 100%
2. FOR Traditional players, THE Victory_Progress SHALL calculate as (current_cloud_adoption / 80) * 100%
3. FOR CloudReseller players, THE Victory_Progress SHALL calculate as (current_employees / 150) * 100%
4. FOR CloudNative players, THE Victory_Progress SHALL calculate as ((capital_progress + security_progress) / 2) * 100%, where capital_progress = min(current_capital / 2000, 1) and security_progress = min(current_security_level / 5, 1)
5. THE Victory_Progress SHALL cap at 100% maximum

### Requirement 4: 遊戲規則說明頁面

**User Story:** As a player, I want to see game rules and victory conditions before starting, so that I understand how to win the game.

#### Acceptance Criteria

1. WHEN the game status changes to "in_progress", THE Game_Rules_Modal SHALL display automatically before the first dice roll
2. THE Game_Rules_Modal SHALL display the game objective and basic rules
3. THE Game_Rules_Modal SHALL display all four company types and their respective victory conditions
4. THE Game_Rules_Modal SHALL display the 30-turn limit rule
5. WHEN the player clicks "開始遊戲" button, THE Game_Rules_Modal SHALL close and allow gameplay to begin
6. THE Game_Rules_Modal SHALL highlight the current player's victory condition

### Requirement 5: 勝利畫面顯示

**User Story:** As a player, I want to see a victory screen when someone wins, so that the game ending is clear and celebratory.

#### Acceptance Criteria

1. WHEN a player achieves their victory condition, THE Game_UI SHALL display a victory modal
2. THE Victory_Modal SHALL show the winner's name, company type, and achieved condition
3. THE Victory_Modal SHALL show all players' final statistics
4. THE Victory_Modal SHALL provide a "返回大廳" button to return to lobby
5. WHEN the game ends due to turn limit, THE Victory_Modal SHALL show "時間到！" and display the winner based on highest victory progress
