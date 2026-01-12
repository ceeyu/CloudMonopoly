import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { createGame, joinGame } from '../api/client'
import { CompanyType, COMPANY_TYPES } from '../api/types'
import './Lobby.css'

function Lobby() {
  const navigate = useNavigate()
  const [mode, setMode] = useState<'menu' | 'create' | 'join'>('menu')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // 建立遊戲表單
  const [maxPlayers, setMaxPlayers] = useState(2)
  
  // 加入遊戲表單
  const [gameIdInput, setGameIdInput] = useState('')
  const [playerName, setPlayerName] = useState('')
  const [companyType, setCompanyType] = useState<CompanyType>('startup')

  const handleCreateGame = async () => {
    setLoading(true)
    setError(null)
    try {
      console.log('Creating game with max_players:', maxPlayers)
      const response = await createGame({ max_players: maxPlayers })
      console.log('Create game response:', response)
      if (response && response.game_id) {
        // 建立後直接進入加入流程
        setGameIdInput(response.game_id)
        setMode('join')
      } else {
        setError('建立遊戲失敗：伺服器回應格式錯誤')
      }
    } catch (err) {
      console.error('Create game error:', err)
      setError(err instanceof Error ? err.message : '建立遊戲失敗，請確認後端服務是否正常運作')
    } finally {
      setLoading(false)
    }
  }

  const handleJoinGame = async () => {
    if (!gameIdInput.trim() || !playerName.trim()) {
      setError('請填寫所有欄位')
      return
    }

    setLoading(true)
    setError(null)
    try {
      const playerId = `player_${Date.now()}`
      await joinGame(gameIdInput, {
        player_id: playerId,
        player_name: playerName,
        company_type: companyType,
      })
      // 儲存玩家 ID 到 localStorage
      localStorage.setItem('playerId', playerId)
      localStorage.setItem('playerName', playerName)
      // 導航到遊戲頁面
      navigate(`/game/${gameIdInput}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : '加入遊戲失敗')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="lobby">
      <div className="lobby-card">
        <h2>🎲 遊戲大廳</h2>
        
        {error && <div className="error-message">{error}</div>}

        {mode === 'menu' && (
          <div className="menu-buttons">
            <button 
              className="btn btn-primary btn-large"
              onClick={() => setMode('create')}
            >
              🆕 建立新遊戲
            </button>
            <button 
              className="btn btn-secondary btn-large"
              onClick={() => setMode('join')}
            >
              🚪 加入遊戲
            </button>
          </div>
        )}

        {mode === 'create' && (
          <div className="form-section">
            <h3>建立新遊戲</h3>
            <div className="form-group">
              <label>玩家人數</label>
              <select 
                value={maxPlayers} 
                onChange={(e) => setMaxPlayers(Number(e.target.value))}
              >
                <option value={2}>2 人</option>
                <option value={3}>3 人</option>
                <option value={4}>4 人</option>
              </select>
            </div>
            <div className="button-group">
              <button 
                className="btn btn-primary"
                onClick={handleCreateGame}
                disabled={loading}
              >
                {loading ? '建立中...' : '建立遊戲'}
              </button>
              <button 
                className="btn btn-secondary"
                onClick={() => setMode('menu')}
                disabled={loading}
              >
                返回
              </button>
            </div>
          </div>
        )}

        {mode === 'join' && (
          <div className="form-section">
            <h3>加入遊戲</h3>
            {gameIdInput && (
              <div className="game-id-display">
                <span className="game-id-label">🎮 遊戲 ID：</span>
                <span className="game-id-value">{gameIdInput}</span>
                <button 
                  className="btn-copy"
                  onClick={() => navigator.clipboard.writeText(gameIdInput)}
                  title="複製 ID"
                >
                  📋
                </button>
              </div>
            )}
            {!gameIdInput && (
              <div className="form-group">
                <label>遊戲 ID</label>
                <input
                  type="text"
                  value={gameIdInput}
                  onChange={(e) => setGameIdInput(e.target.value)}
                  placeholder="輸入遊戲 ID"
                />
              </div>
            )}
            <div className="form-group">
              <label>玩家名稱</label>
              <input
                type="text"
                value={playerName}
                onChange={(e) => setPlayerName(e.target.value)}
                placeholder="輸入你的名稱"
              />
            </div>
            <div className="form-group">
              <label>選擇公司類型</label>
              <div className="company-types">
                {(Object.keys(COMPANY_TYPES) as CompanyType[]).map((type) => {
                  const info = COMPANY_TYPES[type]
                  return (
                    <div
                      key={type}
                      className={`company-card ${companyType === type ? 'selected' : ''}`}
                      onClick={() => setCompanyType(type)}
                    >
                      <div className="company-card-header">
                        <span className="company-icon">{info.icon}</span>
                        <h4>{info.name}</h4>
                      </div>
                      <p className="company-desc">{info.description}</p>
                      <div className="company-stats">
                        <div className="stat-item">
                          <span className="stat-label">💰 資本</span>
                          <span className="stat-value">{info.initialStats.capital} 萬</span>
                        </div>
                        <div className="stat-item">
                          <span className="stat-label">👥 員工</span>
                          <span className="stat-value">{info.initialStats.employees} 人</span>
                        </div>
                        <div className="stat-item">
                          <span className="stat-label">🔒 資安</span>
                          <span className="stat-value">Lv.{info.initialStats.securityLevel}</span>
                        </div>
                        <div className="stat-item">
                          <span className="stat-label">☁️ 雲端</span>
                          <span className="stat-value">{info.initialStats.cloudAdoption}%</span>
                        </div>
                      </div>
                    </div>
                  )
                })}
              </div>
            </div>
            <div className="button-group">
              <button 
                className="btn btn-primary"
                onClick={handleJoinGame}
                disabled={loading}
              >
                {loading ? '加入中...' : '加入遊戲'}
              </button>
              <button 
                className="btn btn-secondary"
                onClick={() => setMode('menu')}
                disabled={loading}
              >
                返回
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default Lobby
