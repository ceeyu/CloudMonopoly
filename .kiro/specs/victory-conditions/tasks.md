# Implementation Plan: Victory Conditions

## Overview

實作角色專屬勝利條件機制，包含後端勝利檢查邏輯、前端規則說明與勝利畫面。

## Tasks

- [x] 1. 後端：新增勝利條件型別與檢查邏輯
  - [x] 1.1 在 `internal/game/victory.go` 新增 VictoryCondition 型別和預設條件
    - 定義 VictoryCondition 結構
    - 定義 DefaultVictoryConditions map
    - _Requirements: 1.1, 1.2, 1.3, 1.4_
  - [x] 1.2 實作 CheckVictory 函數
    - 根據公司類型檢查對應勝利條件
    - _Requirements: 1.1, 1.2, 1.3, 1.4_
  - [x] 1.3 實作 CalculateVictoryProgress 函數
    - 計算各公司類型的勝利進度百分比
    - 確保進度上限為 100%
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_
  - [x] 1.4 撰寫勝利條件檢查的單元測試
    - 測試各公司類型的勝利檢查
    - 測試進度計算公式
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 3.1-3.5_

- [x] 2. 後端：整合勝利檢查到遊戲引擎
  - [x] 2.1 擴充 Game 和 PlayerState 結構
    - 新增 MaxTurnsPerPlayer、WinnerID、WinReason 到 Game
    - 新增 VictoryProgress 到 PlayerState
    - _Requirements: 2.1, 2.4_
  - [x] 2.2 修改 ExecuteTurn 函數
    - 每次回合後檢查勝利條件
    - 檢查回合限制 (30 回合)
    - 更新玩家勝利進度
    - _Requirements: 1.5, 2.1, 2.2_
  - [x] 2.3 實作 DetermineWinnerByProgress 函數
    - 回合限制時根據進度判定贏家
    - _Requirements: 2.3_
  - [x] 2.4 撰寫遊戲引擎整合測試
    - 測試完整遊戲流程
    - 測試回合限制邏輯
    - _Requirements: 2.1, 2.2, 2.3_

- [x] 3. 後端：更新 API 回應
  - [x] 3.1 修改 GetGameState handler
    - 回應中包含 winner_id、win_reason、max_turns_per_player
    - 回應中包含各玩家的 victory_progress
    - _Requirements: 2.4, 3.1-3.5_
  - [x] 3.2 修改 ExecuteTurn handler
    - 回應中包含勝利狀態
    - _Requirements: 1.1-1.5_

- [x] 4. Checkpoint - 確保後端測試通過
  - 確保所有測試通過，如有問題請詢問使用者

- [x] 5. 前端：更新 API 型別定義
  - [x] 5.1 更新 types.ts
    - 新增 victory_progress、winner_id、win_reason 等欄位
    - _Requirements: 2.4, 3.1-3.5_

- [x] 6. 前端：建立 GameRulesModal 組件
  - [x] 6.1 建立 GameRulesModal.tsx 和 GameRulesModal.css
    - 顯示遊戲目標和基本規則
    - 顯示四種公司類型的勝利條件
    - 顯示 30 回合限制說明
    - 高亮當前玩家的勝利條件
    - 提供「開始遊戲」按鈕
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6_

- [x] 7. 前端：建立 VictoryModal 組件
  - [x] 7.1 建立 VictoryModal.tsx 和 VictoryModal.css
    - 顯示贏家資訊
    - 顯示所有玩家最終統計
    - 區分「達成條件」和「時間到」兩種結束方式
    - 提供「返回大廳」按鈕
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_

- [x] 8. 前端：擴充 CompanyStatus 組件
  - [x] 8.1 新增勝利進度條顯示
    - 顯示當前進度百分比
    - 顯示目標描述
    - _Requirements: 3.1-3.5_

- [x] 9. 前端：整合到 Game 頁面
  - [x] 9.1 修改 Game.tsx
    - 遊戲開始時顯示 GameRulesModal
    - 勝利時顯示 VictoryModal
    - 傳遞勝利進度給 CompanyStatus
    - _Requirements: 4.1, 5.1_

- [x] 10. Checkpoint - 確保前端功能正常
  - 確保所有功能正常運作，如有問題請詢問使用者

- [x] 11. 部署更新
  - [x] 11.1 重新建置並部署後端到 App Runner
    - 建置 Docker image
    - 推送到 ECR
    - 更新 App Runner 服務
  - [x] 11.2 重新建置並部署前端到 S3/CloudFront
    - npm run build
    - 上傳到 S3
    - 清除 CloudFront 快取

- [x] 12. Final Checkpoint - 確保部署成功
  - 確保所有功能在生產環境正常運作

## Notes

- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties
- Unit tests validate specific examples and edge cases
