# Implementation Plan: AWS Learning Game

## Overview

使用 Golang 開發 AWS 學習遊戲，採用漸進式實作方式。先建立核心遊戲邏輯，再逐步加入事件系統、決策引擎，最後整合 Web API 和前端。

## Implementation Status

✅ **所有任務已完成** - 專案已完整實作，包含所有核心功能、屬性測試、REST API 和 Web 前端。

## Tasks

- [x] 1. 專案初始化與核心資料結構
  - [x] 1.1 建立 Go module 和專案目錄結構
    - 建立 `cmd/`, `internal/`, `pkg/`, `api/` 目錄
    - 初始化 go.mod
    - _Requirements: 9.1_
  - [x] 1.2 實作 Company 資料結構和 CompanyManager
    - 定義 Company struct 和 CompanyType
    - 實作 CreateCompany, UpdateCapital, UpdateEmployees
    - 實作 CompanyDefaults 初始值
    - _Requirements: 1.1, 1.2, 1.3, 1.4_
  - [x] 1.3 撰寫 Company 屬性測試
    - **Property 1: Company Initialization Completeness**
    - **Validates: Requirements 1.2**
  - [x] 1.4 實作 Board 和 Cell 資料結構
    - 定義 Board, Cell, CellType structs
    - 實作 CreateBoard, GetCell
    - _Requirements: 2.5_
  - [x] 1.5 撰寫 Board 屬性測試
    - **Property 4: Board Size Constraint**
    - **Validates: Requirements 2.5**

- [x] 2. 遊戲引擎核心邏輯
  - [x] 2.1 實作 GameEngine 基礎功能
    - 實作 CreateGame, JoinGame, StartGame
    - 實作 GetGameState
    - _Requirements: 7.1_
  - [x] 2.2 撰寫玩家數量驗證屬性測試
    - **Property 18: Player Count Validation**
    - **Validates: Requirements 7.1**
  - [x] 2.3 實作骰子和移動邏輯
    - 實作 ExecuteTurn 的骰子擲骰
    - 實作位置計算 (含繞圈處理)
    - _Requirements: 2.1, 2.2, 2.6_
  - [x] 2.4 撰寫位置移動屬性測試
    - **Property 2: Position Movement Correctness**
    - **Property 5: Circuit Completion Bonus**
    - **Validates: Requirements 2.2, 2.6**
  - [x] 2.5 實作格子著陸資源更新
    - 更新 Capital 和 Employees
    - _Requirements: 2.3, 2.4_
  - [x] 2.6 撰寫資源更新屬性測試
    - **Property 3: Cell Landing Resource Update**
    - **Validates: Requirements 2.3, 2.4**

- [x] 3. Checkpoint - 核心遊戲邏輯驗證
  - 確保所有測試通過，如有問題請詢問使用者

- [x] 4. 事件系統
  - [x] 4.1 實作 Event 資料結構
    - 定義 Event, EventType, EventChoice, EventOutcome structs
    - 定義 EventContext, ChoiceRequirements, ChoiceOutcomes
    - _Requirements: 3.1, 3.2, 3.3, 3.4_
  - [x] 4.2 實作 EventSystem 介面
    - 實作 GetEvent, GetRandomEvent
    - 實作事件類型與格子類型對應
    - _Requirements: 3.1, 3.2, 3.3_
  - [x] 4.3 撰寫事件類型匹配屬性測試
    - **Property 6: Event Type Cell Matching**
    - **Validates: Requirements 3.1, 3.2, 3.3**
  - [x] 4.4 實作資安事件處理
    - 確保資安事件包含緩解策略選項
    - _Requirements: 3.4_
  - [x] 4.5 撰寫資安事件屬性測試
    - **Property 7: Security Event Mitigation Choices**
    - **Validates: Requirements 3.4**
  - [x] 4.6 實作事件顯示功能
    - 確保事件包含完整資訊
    - _Requirements: 3.6_
  - [x] 4.7 撰寫事件顯示完整性屬性測試
    - **Property 8: Event Display Completeness**
    - **Validates: Requirements 3.6**
  - [x] 4.8 建立初始事件資料集
    - 建立至少 50 個事件 (JSON 格式)
    - 涵蓋擴廠、跨國合作、資安事件等真實情境
    - _Requirements: 3.5, 3.7_

- [x] 5. 決策引擎
  - [x] 5.1 實作 DecisionEngine 介面
    - 實作 EvaluateDecision, ExecuteDecision
    - _Requirements: 4.1, 4.3_
  - [x] 5.2 撰寫基礎設施決策選項屬性測試
    - **Property 9: Infrastructure Decision Options**
    - **Validates: Requirements 4.1**
  - [x] 5.3 實作決策結果計算
    - 基於公司屬性計算結果
    - 實作成功/失敗邏輯
    - _Requirements: 4.3, 4.4_
  - [x] 5.4 撰寫決策結果計算屬性測試
    - **Property 11: Decision Outcome Calculation**
    - **Property 12: Decision Evaluation Completeness**
    - **Validates: Requirements 4.3, 4.4**
  - [x] 5.5 實作不相容決策懲罰
    - 當公司屬性不符合要求時套用懲罰
    - _Requirements: 4.5_
  - [x] 5.6 撰寫懲罰機制屬性測試
    - **Property 13: Incompatible Decision Penalty**
    - **Validates: Requirements 4.5**
  - [x] 5.7 實作方案比較功能
    - 實作 GetComparison
    - _Requirements: 5.2_
  - [x] 5.8 撰寫比較表屬性測試
    - **Property 15: Comparison Table Generation**
    - **Validates: Requirements 5.2**

- [x] 6. AWS 服務目錄
  - [x] 6.1 實作 AWSServiceCatalog 介面
    - 實作 GetService, GetServicesByCategory
    - 實作 GetRecommendedServices
    - _Requirements: 4.2_
  - [x] 6.2 撰寫服務類別相關性屬性測試
    - **Property 10: AWS Service Category Relevance**
    - **Validates: Requirements 4.2**
  - [x] 6.3 建立 AWS 服務資料集
    - 建立常用 AWS 服務資料 (JSON 格式)
    - 包含 SAA 考試相關主題
    - _Requirements: 4.2, 4.7_

- [x] 7. Checkpoint - 事件與決策系統驗證
  - 確保所有測試通過，如有問題請詢問使用者

- [x] 8. 架構視覺化
  - [x] 8.1 實作 ArchitectureVisualizer 介面
    - 實作 GenerateDiagram (ASCII/Mermaid 格式)
    - _Requirements: 5.1, 5.3_
  - [x] 8.2 撰寫架構圖產生屬性測試
    - **Property 14: Architecture Diagram Generation**
    - **Validates: Requirements 5.1**
  - [x] 8.3 實作公司架構狀態顯示
    - 實作 GenerateCompanyArchitecture
    - _Requirements: 5.4_

- [x] 9. 學習回饋系統
  - [x] 9.1 實作決策回饋功能
    - 在 DecisionResult 中加入解釋和學習要點
    - 加入 AWS SAA 考點參考
    - _Requirements: 6.1, 6.2_
  - [x] 9.2 撰寫決策回饋完整性屬性測試
    - **Property 16: Decision Feedback Completeness**
    - **Validates: Requirements 6.1, 6.2**
  - [x] 9.3 實作遊戲結束摘要
    - 產生所有決策記錄摘要
    - _Requirements: 6.3_
  - [x] 9.4 撰寫遊戲結束摘要屬性測試
    - **Property 17: Game End Summary**
    - **Validates: Requirements 6.3**
  - [x] 9.5 實作玩家進度追蹤
    - 追蹤決策歷史和學習進度
    - _Requirements: 6.4_

- [x] 10. 多人遊戲邏輯
  - [x] 10.1 實作回合管理
    - 實作回合順序和自動切換
    - _Requirements: 7.2, 7.3_
  - [x] 10.2 撰寫回合公平性屬性測試
    - **Property 19: Turn Order Fairness**
    - **Validates: Requirements 7.2, 7.3**
  - [x] 10.3 實作勝負判定
    - 基於最終資本判定贏家
    - _Requirements: 7.4_
  - [x] 10.4 撰寫勝負判定屬性測試
    - **Property 20: Winner Determination**
    - **Validates: Requirements 7.4**

- [x] 11. Checkpoint - 完整遊戲邏輯驗證
  - 確保所有測試通過，如有問題請詢問使用者

- [x] 12. 遊戲狀態持久化
  - [x] 12.1 實作 GameState 序列化
    - 實作 JSON 序列化/反序列化
    - _Requirements: 8.4_
  - [x] 12.2 撰寫序列化 round-trip 屬性測試
    - **Property 21: Game State Serialization Round-Trip**
    - **Validates: Requirements 8.4, 8.5**
  - [x] 12.3 實作 DynamoDB 儲存層
    - 實作 SaveGame, LoadGame
    - 設定 DynamoDB table schema
    - _Requirements: 8.1, 8.2, 8.3_
  - [x] 12.4 撰寫存取一致性屬性測試
    - **Property 22: Save and Load Consistency**
    - **Validates: Requirements 8.1, 8.2**

- [x] 13. REST API 層
  - [x] 13.1 設定 API 框架
    - 使用 Gin 或 Echo 框架
    - 設定路由和中介軟體
    - _Requirements: 9.1_
  - [x] 13.2 實作遊戲管理 API
    - POST /games - 建立遊戲
    - GET /games/{id} - 取得遊戲狀態
    - POST /games/{id}/join - 加入遊戲
    - POST /games/{id}/start - 開始遊戲
    - _Requirements: 9.1_
  - [x] 13.3 實作回合操作 API
    - POST /games/{id}/turn - 執行回合動作
    - POST /games/{id}/decision - 提交決策
    - _Requirements: 9.1, 9.4_
  - [x] 13.4 實作存檔 API
    - POST /games/{id}/save - 儲存遊戲
    - GET /games/{id}/load - 載入遊戲
    - _Requirements: 8.1, 8.2_

- [x] 14. Checkpoint - API 整合驗證
  - 確保所有 API 端點正常運作，如有問題請詢問使用者

- [x] 15. Web 前端基礎
  - [x] 15.1 建立前端專案
    - 使用 React 或 Vue.js
    - 設定與後端 API 連接
    - _Requirements: 9.2_
  - [x] 15.2 實作遊戲大廳頁面
    - 建立/加入遊戲功能
    - 選擇公司類型
    - _Requirements: 1.1, 9.2_
  - [x] 15.3 實作遊戲棋盤介面
    - 顯示棋盤和玩家位置
    - 顯示公司狀態
    - _Requirements: 9.3_
  - [x] 15.4 實作事件和決策介面
    - 顯示事件內容
    - 顯示決策選項和架構圖
    - 提交決策並顯示結果
    - _Requirements: 9.3, 9.4_
  - [x] 15.5 實作響應式設計
    - 支援桌面和行動裝置
    - _Requirements: 9.5_

- [x] 16. Final Checkpoint - 完整系統驗證
  - 確保所有測試通過，進行端對端測試，如有問題請詢問使用者

## Notes

- 所有任務都必須執行，包含屬性測試
- 每個任務都參考具體需求以確保可追溯性
- Checkpoint 用於確保階段性驗證
- 屬性測試驗證通用正確性屬性
- 單元測試驗證特定範例和邊界情況
