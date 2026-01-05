# Requirements Document

## Introduction

AWS Learning Game 是一款結合大富翁遊戲機制的 AWS SAA 證照學習遊戲。玩家扮演不同類型的公司（新創、傳產、雲端代理商、雲端公司），透過骰子前進、遭遇各種商業情境，學習如何根據公司狀況選擇適當的 AWS 服務或地端方案。遊戲使用 Golang 開發，部署於 AWS 平台。

## Glossary

- **Game_Engine**: 負責遊戲核心邏輯的系統，包含回合管理、骰子機制、玩家移動
- **Company**: 玩家在遊戲中扮演的公司實體，具有資本額、員工數、產品週期等屬性
- **Board**: 遊戲棋盤，包含多個格子，每個格子有不同事件
- **Event_System**: 處理機會、命運、關卡等事件的系統
- **Decision_Engine**: 處理玩家決策並計算結果的系統
- **AWS_Service_Catalog**: 儲存所有可選 AWS 服務及其屬性的資料庫
- **Architecture_Visualizer**: 產生架構圖和方案比較的元件
- **Player**: 參與遊戲的使用者
- **Capital**: 公司資本額，影響可選方案和遊戲勝負
- **Turn**: 遊戲回合，玩家依序執行動作

## Requirements

### Requirement 1: 公司角色系統

**User Story:** As a Player, I want to choose different company types at game start, so that I can experience different decision-making scenarios based on company characteristics.

#### Acceptance Criteria

1. WHEN the game starts, THE Game_Engine SHALL present at least 4 company types: 新創公司、傳產公司、雲端代理商、雲端公司
2. WHEN a Player selects a company type, THE Game_Engine SHALL initialize the Company with predefined attributes including Capital, employee count, product cycle, and international presence
3. THE Company SHALL have attributes that affect available AWS service options and decision outcomes
4. WHEN displaying company information, THE Game_Engine SHALL show current Capital, employee count, product status, and company type

### Requirement 2: 遊戲棋盤與移動機制

**User Story:** As a Player, I want to roll dice and move on the board, so that I can progress through the game and encounter various events.

#### Acceptance Criteria

1. WHEN a Player starts their Turn, THE Game_Engine SHALL allow the Player to roll dice
2. WHEN dice are rolled, THE Game_Engine SHALL move the Player's position by the dice value
3. WHEN a Player lands on a Board cell, THE Game_Engine SHALL increase the Company Capital by the cell's base value
4. WHEN a Player lands on a Board cell, THE Game_Engine SHALL increase the Company employee count by a small incremental value
5. THE Board SHALL contain at least 30 cells with different event types
6. WHEN a Player completes a full circuit of the Board, THE Game_Engine SHALL grant bonus Capital and employee growth

### Requirement 3: 事件系統

**User Story:** As a Player, I want to encounter various events like opportunity, fate, and challenges, so that I can face realistic business scenarios requiring AWS knowledge.

#### Acceptance Criteria

1. WHEN a Player lands on an opportunity cell, THE Event_System SHALL present a beneficial scenario such as expansion or international partnership, potentially causing significant employee growth
2. WHEN a Player lands on a fate cell, THE Event_System SHALL present a random scenario that may be positive or negative, including security incidents causing financial loss or data concerns
3. WHEN a Player lands on a challenge cell, THE Event_System SHALL present a technical challenge requiring AWS architecture decisions
4. WHEN a security incident event occurs, THE Event_System SHALL require the Player to choose mitigation strategies affecting Capital and data integrity
5. THE Event_System SHALL store at least 50 unique events based on real-world business cases and actual industry incidents
6. WHEN presenting an event, THE Event_System SHALL display realistic business context, consequences, and available strategic options
7. THE Event_System SHALL include events covering: 擴廠、跨國合作、資安事件、資料外洩、系統故障、業務成長、法規遵循等真實情境

### Requirement 4: AWS 服務決策機制

**User Story:** As a Player, I want to make decisions between on-premises solutions and AWS services, so that I can learn the trade-offs and appropriate use cases for each service.

#### Acceptance Criteria

1. WHEN an event requires infrastructure decision, THE Decision_Engine SHALL present both on-premises and AWS cloud options based on realistic business logic
2. WHEN presenting AWS options, THE AWS_Service_Catalog SHALL provide relevant services based on the scenario (compute, storage, database, networking, security, etc.)
3. WHEN a Player selects an option, THE Decision_Engine SHALL calculate the outcome based on Company attributes and service characteristics using real-world cost and performance models
4. THE Decision_Engine SHALL consider Capital cost, implementation time, scalability, maintenance requirements, and security implications in outcome calculation
5. IF a Player selects an option incompatible with Company attributes, THEN THE Decision_Engine SHALL apply realistic penalty effects such as budget overrun or implementation delays
6. WHEN a security-related decision is required, THE Decision_Engine SHALL evaluate the choice against industry best practices and compliance requirements
7. THE Decision_Engine SHALL reference actual AWS pricing models and real-world deployment considerations

### Requirement 5: 架構圖與方案比較

**User Story:** As a Player, I want to see architecture diagrams and comparisons for each decision, so that I can understand the technical implications of my choices.

#### Acceptance Criteria

1. WHEN presenting decision options, THE Architecture_Visualizer SHALL display a basic architecture diagram for each option
2. WHEN comparing options, THE Architecture_Visualizer SHALL show a comparison table including cost, scalability, complexity, and AWS exam relevance
3. THE Architecture_Visualizer SHALL use text-based or simple graphical representation suitable for web display
4. WHEN a decision is made, THE Architecture_Visualizer SHALL show the resulting architecture state of the Company

### Requirement 6: 遊戲結果與學習回饋

**User Story:** As a Player, I want to receive feedback on my decisions and learn from the outcomes, so that I can improve my AWS knowledge.

#### Acceptance Criteria

1. WHEN a decision outcome is calculated, THE Game_Engine SHALL display the result with explanation of why the choice was good or suboptimal
2. WHEN displaying decision feedback, THE Game_Engine SHALL reference relevant AWS SAA exam topics and best practices
3. WHEN the game ends, THE Game_Engine SHALL provide a summary of all decisions and learning points
4. THE Game_Engine SHALL track Player progress and highlight areas needing improvement

### Requirement 7: 多人遊戲支援

**User Story:** As a Player, I want to play with other players, so that I can compete and learn together.

#### Acceptance Criteria

1. THE Game_Engine SHALL support 2-4 Players in a single game session
2. WHEN multiple Players are in a game, THE Game_Engine SHALL manage turn order fairly
3. WHEN a Player's Turn ends, THE Game_Engine SHALL automatically advance to the next Player
4. THE Game_Engine SHALL determine a winner based on final Capital and company growth metrics

### Requirement 8: 遊戲狀態持久化

**User Story:** As a Player, I want to save and resume my game, so that I can continue playing later.

#### Acceptance Criteria

1. WHEN a Player requests to save, THE Game_Engine SHALL persist the current game state
2. WHEN a Player requests to load a saved game, THE Game_Engine SHALL restore the complete game state
3. THE Game_Engine SHALL store game state using AWS services (DynamoDB or S3)
4. WHEN serializing game state, THE Game_Engine SHALL encode it using JSON format
5. WHEN deserializing game state, THE Game_Engine SHALL validate and restore the exact previous state

### Requirement 9: Web 介面

**User Story:** As a Player, I want to access the game through a web browser, so that I can play without installing additional software.

#### Acceptance Criteria

1. THE Game_Engine SHALL expose a REST API for all game operations
2. WHEN a Player accesses the game URL, THE System SHALL serve a web-based user interface
3. THE web interface SHALL display the game board, company status, and current events
4. WHEN a Player makes a decision, THE web interface SHALL send the choice to the backend and display the result
5. THE web interface SHALL be responsive and work on both desktop and mobile browsers
