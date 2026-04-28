import { PlayerState, WinReason, COMPANY_TYPES, VICTORY_CONDITIONS, CompanyType } from '../api/types'
import './VictoryModal.css'

interface VictoryModalProps {
  isOpen: boolean
  winner: PlayerState | null
  players: PlayerState[]
  winReason: WinReason
  onClose: () => void
}

function VictoryModal({ isOpen, winner, players, winReason, onClose }: VictoryModalProps) {
  if (!isOpen || !winner) return null

  const winnerCompanyType = winner.company?.type as CompanyType | undefined
  const winnerCompanyInfo = winnerCompanyType ? COMPANY_TYPES[winnerCompanyType] : null
  const winnerVictoryInfo = winnerCompanyType ? VICTORY_CONDITIONS[winnerCompanyType] : null

  // 根據勝利進度排序玩家
  const sortedPlayers = [...players].sort((a, b) => b.victory_progress - a.victory_progress)

  return (
    <div className="modal-overlay victory-modal-overlay">
      <div className="modal-content victory-modal-content" onClick={(e) => e.stopPropagation()}>
        {/* 勝利標題 */}
        <div className="modal-header victory-header">
          {winReason === 'condition_met' ? (
            <h2>🏆 勝利！</h2>
          ) : (
            <h2>⏱️ 時間到！</h2>
          )}
        </div>

        <div className="modal-body victory-body">
          {/* 贏家資訊 */}
          <section className="winner-section">
            <div className="winner-card">
              <div className="winner-icon">
                {winnerCompanyInfo?.icon || '🏆'}
              </div>
              <div className="winner-info">
                <h3 className="winner-name">{winner.player_name}</h3>
                <p className="winner-company">
                  {winnerCompanyInfo?.name || '未知公司'}
                </p>
              </div>
            </div>

            {/* 勝利原因說明 */}
            <div className="win-reason-box">
              {winReason === 'condition_met' ? (
                <>
                  <span className="win-reason-label">達成條件</span>
                  <p className="win-reason-text">
                    {winnerVictoryInfo?.target || '達成勝利條件'}
                  </p>
                </>
              ) : (
                <>
                  <span className="win-reason-label">最高進度</span>
                  <p className="win-reason-text">
                    勝利進度 {winner.victory_progress.toFixed(1)}%
                  </p>
                </>
              )}
            </div>
          </section>

          {/* 所有玩家最終統計 */}
          <section className="final-stats-section">
            <h3>📊 最終排名</h3>
            <div className="players-ranking">
              {sortedPlayers.map((player, index) => {
                const companyType = player.company?.type as CompanyType | undefined
                const companyInfo = companyType ? COMPANY_TYPES[companyType] : null
                const isWinner = player.player_id === winner.player_id

                return (
                  <div 
                    key={player.player_id} 
                    className={`ranking-card ${isWinner ? 'winner' : ''}`}
                  >
                    <div className="ranking-position">
                      {index === 0 ? '🥇' : index === 1 ? '🥈' : index === 2 ? '🥉' : `#${index + 1}`}
                    </div>
                    <div className="ranking-player">
                      <span className="ranking-icon">{companyInfo?.icon || '🏢'}</span>
                      <span className="ranking-name">{player.player_name}</span>
                    </div>
                    <div className="ranking-progress">
                      <div className="progress-bar-container">
                        <div 
                          className="progress-bar-fill"
                          style={{ width: `${Math.min(player.victory_progress, 100)}%` }}
                        />
                      </div>
                      <span className="progress-text">{player.victory_progress.toFixed(1)}%</span>
                    </div>
                    <div className="ranking-stats">
                      {player.company && (
                        <>
                          <span className="stat">💰 {player.company.capital}萬</span>
                          <span className="stat">👥 {player.company.employees}人</span>
                        </>
                      )}
                    </div>
                  </div>
                )
              })}
            </div>
          </section>

          {/* 遊戲統計摘要 */}
          <section className="game-summary-section">
            <h3>📈 遊戲統計</h3>
            <div className="summary-grid">
              <div className="summary-item">
                <span className="summary-label">總回合數</span>
                <span className="summary-value">{winner.turns_played}</span>
              </div>
              <div className="summary-item">
                <span className="summary-label">參與玩家</span>
                <span className="summary-value">{players.length} 人</span>
              </div>
            </div>
          </section>
        </div>

        <div className="modal-footer victory-footer">
          <button className="btn btn-primary btn-return-lobby" onClick={onClose}>
            🏠 返回大廳
          </button>
        </div>
      </div>
    </div>
  )
}

export default VictoryModal
