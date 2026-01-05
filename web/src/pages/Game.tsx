import { useState, useEffect, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { getGameState, executeTurn, startGame, submitDecision, getRandomEvent } from '../api/client'
import { GameStateResponse, TurnResponse, PlayerState, EventResponse, DecisionResponse } from '../api/types'
import GameBoard from '../components/GameBoard'
import CompanyStatus from '../components/CompanyStatus'
import EventModal from '../components/EventModal'
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

  const playerId = localStorage.getItem('playerId') || ''

  const fetchGameState = useCallback(async () => {
    if (!gameId) return
    try {
      const state = await getGameState(gameId)
      setGameState(state)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : '載入遊戲失敗')
    } finally {
      setLoading(false)
    }
  }, [gameId])

  useEffect(() => {
    fetchGameState()
    // 定期更新遊戲狀態
    const interval = setInterval(fetchGameState, 3000)
    return () => clearInterval(interval)
  }, [fetchGameState])

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
          // 如果取得事件失敗，仍然顯示回合結果
        }
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
            {gameState?.players.map((player: PlayerState) => (
              <CompanyStatus 
                key={player.player_id}
                player={player}
                isCurrentPlayer={player.player_id === playerId}
                isCurrentTurn={player.player_id === gameState.current_player_id}
              />
            ))}
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
                <p>恭喜完成遊戲！</p>
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
    </div>
  )
}

export default Game
