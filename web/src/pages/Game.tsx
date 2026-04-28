import { useState, useEffect, useCallback, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { getGameState, executeTurn, startGame, submitDecision, getRandomEvent } from '../api/client'
import { GameStateResponse, TurnResponse, PlayerState, EventResponse, DecisionResponse, CompanyType } from '../api/types'
import GameBoard from '../components/GameBoard'
import CompanyStatus from '../components/CompanyStatus'
import EventModal from '../components/EventModal'
import GameRulesModal from '../components/GameRulesModal'
import VictoryModal from '../components/VictoryModal'
import './Game.css'

function Game() {
  const { gameId } = useParams<{ gameId: string }>()
  const navigate = useNavigate()
  const [gameState, setGameState] = useState<GameStateResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [turnResult, setTurnResult] = useState<TurnResponse | null>(null)
  const [showEvent, setShowEvent] = useState(false)
  const [currentEvent, setCurrentEvent] = useState<EventResponse | null>(null)
  const [actionLoading, setActionLoading] = useState(false)
  
  // 遊戲規則彈窗狀態 - Requirements 4.1
  const [showRulesModal, setShowRulesModal] = useState(false)
  const hasShownRulesRef = useRef(false)
  
  // 勝利彈窗狀態 - Requirements 5.1
  const [showVictoryModal, setShowVictoryModal] = useState(false)
  const hasShownVictoryRef = useRef(false)

  const playerId = sessionStorage.getItem('playerId') || ''

  const gameStateRef = useRef<GameStateResponse | null>(null)

  const fetchGameState = useCallback(async () => {
    if (!gameId) return
    try {
      const state = await getGameState(gameId)
      
      // 如果當前玩家改變了，清除回合結果（除非正在顯示事件）
      if (gameStateRef.current && state.current_player_id !== gameStateRef.current.current_player_id && !showEvent) {
        setTurnResult(null)
      }
      
      gameStateRef.current = state
      setGameState(state)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : '載入遊戲失敗')
    } finally {
      setLoading(false)
    }
  }, [gameId, showEvent])

  useEffect(() => {
    fetchGameState()
    // 定期更新遊戲狀態
    const interval = setInterval(fetchGameState, 3000)
    return () => clearInterval(interval)
  }, [fetchGameState])

  // 遊戲開始時顯示規則彈窗 - Requirements 4.1
  useEffect(() => {
    if (gameState?.status === 'in_progress' && !hasShownRulesRef.current) {
      setShowRulesModal(true)
      hasShownRulesRef.current = true
    }
  }, [gameState?.status])

  // 遊戲結束時顯示勝利彈窗 - Requirements 5.1
  useEffect(() => {
    if (gameState?.status === 'finished' && gameState?.winner_id && !hasShownVictoryRef.current) {
      setShowVictoryModal(true)
      hasShownVictoryRef.current = true
    }
  }, [gameState?.status, gameState?.winner_id])

  const handleStartGame = async () => {
    if (!gameId) return
    setActionLoading(true)
    try {
      await startGame(gameId)
      await fetchGameState()
    } catch (err) {
      setError(err instanceof Error ? err.message : '開始遊戲失敗')
    } finally {
      setActionLoading(false)
    }
  }

  const handleRollDice = async () => {
    if (!gameId) return
    setActionLoading(true)
    try {
      const result = await executeTurn(gameId, {
        player_id: playerId,
        action_type: 'roll_dice',
      })
      setTurnResult(result)
      
      // 如果需要決策，根據格子類型取得對應事件
      if (result.decision_required && result.cell_type) {
        try {
          const event = await getRandomEvent(gameId, result.cell_type)
          setCurrentEvent(event)
          setShowEvent(true)
        } catch (eventErr) {
          console.error('Failed to get event:', eventErr)
          // 如果取得事件失敗，清除回合結果讓遊戲繼續
          setTimeout(() => {
            setTurnResult(null)
          }, 3000)
        }
      } else {
        // 如果不需要決策（普通格子或獎勵格），3秒後清除回合結果
        setTimeout(() => {
          setTurnResult(null)
        }, 3000)
      }
      
      await fetchGameState()
    } catch (err) {
      setError(err instanceof Error ? err.message : '執行回合失敗')
    } finally {
      setActionLoading(false)
    }
  }

  const handleDecision = async (eventId: string, choiceId: number): Promise<DecisionResponse> => {
    if (!gameId) throw new Error('遊戲 ID 不存在')
    
    const result = await submitDecision(gameId, {
      player_id: playerId,
      event_id: eventId,
      choice_id: choiceId,
    })
    
    // 決策完成後更新遊戲狀態
    await fetchGameState()
    
    return result
  }

  const handleCloseEvent = () => {
    setShowEvent(false)
    setCurrentEvent(null)
    setTurnResult(null)
  }

  // 關閉規則彈窗 - Requirements 4.5
  const handleStartFromRules = () => {
    setShowRulesModal(false)
  }

  // 關閉勝利彈窗並返回大廳 - Requirements 5.4
  const handleCloseVictory = () => {
    setShowVictoryModal(false)
    navigate('/')
  }

  // 取得當前玩家的公司類型
  const currentPlayerCompanyType = gameState?.players.find(
    p => p.player_id === playerId
  )?.company?.type as CompanyType | null

  // 取得贏家資訊
  const winner = gameState?.winner_id 
    ? gameState.players.find(p => p.player_id === gameState.winner_id) || null
    : null

  const isMyTurn = gameState?.current_player_id === playerId
  const isWaiting = gameState?.status === 'waiting'
  const isFinished = gameState?.status === 'finished'

  if (loading) {
    return (
      <div className="game-loading">
        <div className="spinner"></div>
        <p>載入遊戲中...</p>
      </div>
    )
  }

  if (error && !gameState) {
    return (
      <div className="game-error">
        <h2>❌ 錯誤</h2>
        <p>{error}</p>
        <button className="btn btn-primary" onClick={() => navigate('/')}>
          返回大廳
        </button>
      </div>
    )
  }

  return (
    <div className="game">
      {error && <div className="error-banner">{error}</div>}

      <div className="game-header">
        <div className="game-info">
          <span className="game-id">遊戲 ID: {gameId}</span>
          <span className={`game-status status-${gameState?.status}`}>
            {gameState?.status === 'waiting' && '等待玩家'}
            {gameState?.status === 'in_progress' && '遊戲進行中'}
            {gameState?.status === 'finished' && '遊戲結束'}
          </span>
        </div>
        <div className="turn-info">
          {gameState?.status === 'in_progress' && (
            <>
              <span>回合: {gameState.current_turn}</span>
              <span className={isMyTurn ? 'my-turn' : ''}>
                {isMyTurn ? '🎯 輪到你了！' : `等待 ${gameState.players.find(p => p.player_id === gameState.current_player_id)?.player_name || '其他玩家'}`}
              </span>
            </>
          )}
        </div>
      </div>

      <div className="game-content">
        <div className="game-left">
          <GameBoard 
            boardSize={gameState?.board_size || 30}
            players={gameState?.players || []}
            currentPlayerId={playerId}
            cells={gameState?.cells}
          />
          
          {turnResult && !showEvent && (
            <div className="turn-result">
              <h3>🎲 回合結果</h3>
              <p>骰子點數: <strong>{turnResult.dice_value}</strong></p>
              <p>移動: {turnResult.old_position} → {turnResult.new_position}</p>
              {turnResult.capital_change !== 0 && (
                <p className={turnResult.capital_change > 0 ? 'positive' : 'negative'}>
                  資本變化: {turnResult.capital_change > 0 ? '+' : ''}{turnResult.capital_change} 萬
                </p>
              )}
              {turnResult.employee_change !== 0 && (
                <p className={turnResult.employee_change > 0 ? 'positive' : 'negative'}>
                  員工變化: {turnResult.employee_change > 0 ? '+' : ''}{turnResult.employee_change} 人
                </p>
              )}
              {turnResult.circuit_completed && (
                <p className="bonus">🎉 完成一圈！獲得獎勵！</p>
              )}
            </div>
          )}
        </div>

        <div className="game-right">
          <div className="players-section">
            <h3>👥 玩家列表</h3>
            {gameState?.players && gameState.players.length > 0 ? (
              gameState.players.map((player: PlayerState) => (
                <CompanyStatus 
                  key={player.player_id}
                  player={player}
                  isCurrentPlayer={player.player_id === playerId}
                  isCurrentTurn={player.player_id === gameState.current_player_id}
                />
              ))
            ) : (
              <p className="no-players">尚無玩家加入</p>
            )}
          </div>

          <div className="actions-section">
            {isWaiting && (
              <div className="waiting-room">
                <h3>等待玩家加入</h3>
                <p>目前玩家: {gameState?.players.length} 人</p>
                {gameState && gameState.players.length >= 2 && (
                  <button 
                    className="btn btn-primary"
                    onClick={handleStartGame}
                    disabled={actionLoading}
                  >
                    {actionLoading ? '開始中...' : '開始遊戲'}
                  </button>
                )}
              </div>
            )}

            {gameState?.status === 'in_progress' && isMyTurn && !turnResult && (
              <button 
                className="btn btn-primary btn-large roll-btn"
                onClick={handleRollDice}
                disabled={actionLoading}
              >
                {actionLoading ? '擲骰中...' : '🎲 擲骰子'}
              </button>
            )}

            {isFinished && (
              <div className="game-finished">
                <h3>🏆 遊戲結束</h3>
                {winner ? (
                  <p>恭喜 {winner.player_name} 獲勝！</p>
                ) : (
                  <p>恭喜完成遊戲！</p>
                )}
                <button 
                  className="btn btn-secondary" 
                  onClick={() => setShowVictoryModal(true)}
                  style={{ marginBottom: '0.5rem' }}
                >
                  查看結果
                </button>
                <button className="btn btn-primary" onClick={() => navigate('/')}>
                  返回大廳
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {showEvent && currentEvent && (
        <EventModal
          event={currentEvent}
          onDecision={handleDecision}
          onClose={handleCloseEvent}
        />
      )}

      {/* 遊戲規則彈窗 - Requirements 4.1, 4.2, 4.3, 4.4, 4.5, 4.6 */}
      <GameRulesModal
        isOpen={showRulesModal}
        onStart={handleStartFromRules}
        currentPlayerCompanyType={currentPlayerCompanyType}
      />

      {/* 勝利彈窗 - Requirements 5.1, 5.2, 5.3, 5.4, 5.5 */}
      <VictoryModal
        isOpen={showVictoryModal}
        winner={winner}
        players={gameState?.players || []}
        winReason={gameState?.win_reason || 'condition_met'}
        onClose={handleCloseVictory}
      />
    </div>
  )
}

export default Game
