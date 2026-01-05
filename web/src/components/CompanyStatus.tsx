import { PlayerState, COMPANY_TYPES } from '../api/types'
import './CompanyStatus.css'

interface CompanyStatusProps {
  player: PlayerState
  isCurrentPlayer: boolean
  isCurrentTurn: boolean
}

// 產品週期顯示名稱
const PRODUCT_CYCLE_NAMES: Record<string, string> = {
  development: '開發期',
  launch: '上市期',
  growth: '成長期',
  mature: '成熟期',
}

// 資安等級描述
const SECURITY_LEVEL_DESC: Record<number, { name: string; color: string }> = {
  1: { name: '基礎', color: '#ef4444' },
  2: { name: '標準', color: '#f59e0b' },
  3: { name: '進階', color: '#eab308' },
  4: { name: '企業', color: '#22c55e' },
  5: { name: '最高', color: '#10b981' },
}

function CompanyStatus({ player, isCurrentPlayer, isCurrentTurn }: CompanyStatusProps) {
  const company = player.company
  const companyTypeInfo = company?.type ? COMPANY_TYPES[company.type] : null
  const productCycleName = company?.product_cycle ? PRODUCT_CYCLE_NAMES[company.product_cycle] || company.product_cycle : '未知'
  const securityInfo = company?.security_level ? SECURITY_LEVEL_DESC[company.security_level] || { name: `Lv.${company.security_level}`, color: '#6b7280' } : null

  // 計算雲端採用率的顏色
  const getCloudAdoptionColor = (adoption: number) => {
    if (adoption >= 80) return '#10b981'
    if (adoption >= 50) return '#22c55e'
    if (adoption >= 30) return '#eab308'
    return '#f59e0b'
  }

  return (
    <div className={`company-status ${isCurrentPlayer ? 'is-me' : ''} ${isCurrentTurn ? 'is-turn' : ''}`}>
      {/* 玩家標題 */}
      <div className="player-header">
        <div className="player-info">
          <span className="player-name">
            {companyTypeInfo?.icon} {player.player_name}
          </span>
          {isCurrentPlayer && <span className="me-badge">你</span>}
        </div>
        {isCurrentTurn && <span className="turn-indicator" title="輪到此玩家">🎯</span>}
      </div>
      
      {company && (
        <div className="company-details">
          {/* 公司類型 */}
          <div className="company-type-row">
            <span className="company-type-name">{companyTypeInfo?.name || '未知公司'}</span>
            {company.is_international && <span className="international-badge" title="跨國企業">🌍</span>}
          </div>

          {/* 主要數據 */}
          <div className="stats-grid">
            <div className="stat stat-capital">
              <span className="stat-icon">💰</span>
              <div className="stat-content">
                <span className="stat-label">資本額</span>
                <span className="stat-value">{company.capital.toLocaleString()} 萬</span>
              </div>
            </div>
            <div className="stat stat-employees">
              <span className="stat-icon">👥</span>
              <div className="stat-content">
                <span className="stat-label">員工數</span>
                <span className="stat-value">{company.employees} 人</span>
              </div>
            </div>
          </div>

          {/* 次要數據 */}
          <div className="secondary-stats">
            <div className="secondary-stat">
              <span className="stat-label">🔒 資安等級</span>
              <span 
                className="stat-badge"
                style={{ backgroundColor: securityInfo?.color }}
              >
                {securityInfo?.name}
              </span>
            </div>
            <div className="secondary-stat">
              <span className="stat-label">☁️ 雲端採用</span>
              <div className="progress-bar">
                <div 
                  className="progress-fill"
                  style={{ 
                    width: `${company.cloud_adoption}%`,
                    backgroundColor: getCloudAdoptionColor(company.cloud_adoption)
                  }}
                />
                <span className="progress-text">{company.cloud_adoption}%</span>
              </div>
            </div>
            <div className="secondary-stat">
              <span className="stat-label">📦 產品週期</span>
              <span className="stat-text">{productCycleName}</span>
            </div>
            {company.tech_debt > 0 && (
              <div className="secondary-stat">
                <span className="stat-label">⚠️ 技術債</span>
                <span className="stat-text tech-debt">{company.tech_debt}</span>
              </div>
            )}
          </div>

          {/* 已部署基礎設施 */}
          {company.infrastructure && company.infrastructure.length > 0 && (
            <div className="infrastructure-section">
              <span className="stat-label">🏗️ 已部署服務</span>
              <div className="infrastructure-tags">
                {company.infrastructure.slice(0, 4).map((infra, index) => (
                  <span key={index} className="infra-tag">{infra}</span>
                ))}
                {company.infrastructure.length > 4 && (
                  <span className="infra-more">+{company.infrastructure.length - 4}</span>
                )}
              </div>
            </div>
          )}

          {/* 位置資訊 */}
          <div className="position-info">
            <span className="position-icon">📍</span>
            <span>位置: 格子 {player.position + 1}</span>
          </div>
        </div>
      )}

      {!company && (
        <div className="no-company">
          <span>等待選擇公司類型...</span>
        </div>
      )}
    </div>
  )
}

export default CompanyStatus
