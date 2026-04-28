import { CompanyType, COMPANY_TYPES, VICTORY_CONDITIONS } from '../api/types'
import './GameRulesModal.css'

interface GameRulesModalProps {
  isOpen: boolean
  onStart: () => void
  currentPlayerCompanyType: CompanyType | null
}

// 公司類型順序
const COMPANY_TYPE_ORDER: CompanyType[] = ['startup', 'traditional', 'cloud_reseller', 'cloud_native']

function GameRulesModal({ isOpen, onStart, currentPlayerCompanyType }: GameRulesModalProps) {
  if (!isOpen) return null

  return (
    <div className="modal-overlay rules-modal-overlay">
      <div className="modal-content rules-modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header rules-header">
          <h2>🎮 遊戲規則說明</h2>
        </div>

        <div className="modal-body rules-body">
          {/* 遊戲目標 */}
          <section className="rules-section">
            <h3>🎯 遊戲目標</h3>
            <p className="rules-description">
              在 AWS 雲端學習之旅中，你將扮演一家公司的決策者。透過擲骰子前進、做出商業決策、
              學習 AWS 服務，最終達成你的公司專屬勝利條件！
            </p>
          </section>

          {/* 基本規則 */}
          <section className="rules-section">
            <h3>📋 基本規則</h3>
            <ul className="rules-list">
              <li>每位玩家輪流擲骰子，根據點數在棋盤上前進</li>
              <li>踏上不同格子會觸發不同事件：機會、命運、挑戰等</li>
              <li>做出正確的 AWS 架構決策可以獲得更多資源</li>
              <li>繞行棋盤一圈可獲得額外獎勵</li>
            </ul>
          </section>

          {/* 回合限制 */}
          <section className="rules-section turn-limit-section">
            <h3>⏱️ 回合限制</h3>
            <div className="turn-limit-info">
              <span className="turn-limit-number">30</span>
              <span className="turn-limit-text">每位玩家最多 30 回合</span>
            </div>
            <p className="rules-note">
              若無人在 30 回合內達成勝利條件，將根據勝利進度百分比判定贏家。
            </p>
          </section>

          {/* 勝利條件 */}
          <section className="rules-section victory-section">
            <h3>🏆 各公司類型勝利條件</h3>
            <div className="victory-conditions-grid">
              {COMPANY_TYPE_ORDER.map((type) => {
                const companyInfo = COMPANY_TYPES[type]
                const victoryInfo = VICTORY_CONDITIONS[type]
                const isCurrentPlayer = type === currentPlayerCompanyType

                return (
                  <div 
                    key={type} 
                    className={`victory-condition-card ${isCurrentPlayer ? 'highlighted' : ''}`}
                  >
                    {isCurrentPlayer && (
                      <div className="your-company-badge">👈 你的公司</div>
                    )}
                    <div className="victory-card-header">
                      <span className="company-icon">{companyInfo.icon}</span>
                      <span className="company-name">{companyInfo.name}</span>
                    </div>
                    <div className="victory-card-body">
                      <p className="victory-description">{victoryInfo.description}</p>
                      <div className="victory-target">
                        <span className="target-label">目標：</span>
                        <span className="target-value">{victoryInfo.target}</span>
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>
          </section>

          {/* 提示 */}
          <section className="rules-section tips-section">
            <h3>💡 小提示</h3>
            <ul className="tips-list">
              <li>注意你的公司類型，專注於達成專屬勝利條件</li>
              <li>善用 AWS 服務可以加速達成目標</li>
              <li>留意其他玩家的進度，適時調整策略</li>
            </ul>
          </section>
        </div>

        <div className="modal-footer rules-footer">
          <button className="btn btn-primary btn-start-game" onClick={onStart}>
            🚀 開始遊戲
          </button>
        </div>
      </div>
    </div>
  )
}

export default GameRulesModal
