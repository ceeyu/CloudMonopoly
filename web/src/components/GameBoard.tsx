import { PlayerState, CellType, CELL_TYPES } from '../api/types'
import './GameBoard.css'

interface GameBoardProps {
  boardSize: number
  players: PlayerState[]
  currentPlayerId: string
}

// 玩家顏色
const PLAYER_COLORS = ['#ff6b6b', '#4ecdc4', '#45b7d1', '#96ceb4']

function GameBoard({ boardSize, players, currentPlayerId }: GameBoardProps) {
  // 計算棋盤佈局 - 使用環形佈局
  const cellsPerSide = Math.ceil(boardSize / 4)
  
  // 產生格子位置 (環形棋盤) - 使用百分比實現響應式
  const getCellPosition = (index: number): React.CSSProperties => {
    const side = Math.floor(index / cellsPerSide)
    const posInSide = index % cellsPerSide
    const percent = posInSide / cellsPerSide
    
    // 使用百分比定位，讓格子位置響應式
    const edgeOffset = '3%' // 邊緣偏移
    const trackLength = 80 // 軌道長度百分比

    switch (side) {
      case 0: // 上邊 (左到右)
        return { 
          top: edgeOffset, 
          left: `calc(${percent * trackLength + 10}%)`,
          transform: 'translate(-50%, 0)'
        }
      case 1: // 右邊 (上到下)
        return { 
          top: `calc(${percent * trackLength + 10}%)`, 
          right: edgeOffset,
          transform: 'translate(0, -50%)'
        }
      case 2: // 下邊 (右到左)
        return { 
          bottom: edgeOffset, 
          right: `calc(${percent * trackLength + 10}%)`,
          transform: 'translate(50%, 0)'
        }
      case 3: // 左邊 (下到上)
        return { 
          bottom: `calc(${percent * trackLength + 10}%)`, 
          left: edgeOffset,
          transform: 'translate(0, 50%)'
        }
      default:
        return { top: '0%', left: '0%' }
    }
  }

  // 取得格子類型 (模擬後端邏輯)
  const getCellType = (index: number): CellType => {
    if (index === 0) return 'start'
    if (index % 10 === 0) return 'bonus'
    if (index % 7 === 0) return 'opportunity'
    if (index % 5 === 0) return 'fate'
    if (index % 11 === 0) return 'challenge'
    return 'normal'
  }

  // 取得格子上的玩家
  const getPlayersOnCell = (position: number) => {
    return players.filter(p => p.position === position)
  }

  // 取得玩家在列表中的索引 (用於顏色)
  const getPlayerIndex = (playerId: string) => {
    return players.findIndex(p => p.player_id === playerId)
  }

  return (
    <div className="game-board">
      <div className="board-container">
        {/* 棋盤格子 */}
        {Array.from({ length: boardSize }, (_, i) => {
          const cellType = getCellType(i)
          const cellInfo = CELL_TYPES[cellType]
          const playersOnCell = getPlayersOnCell(i)
          
          return (
            <div
              key={i}
              className={`board-cell cell-${cellType}`}
              style={getCellPosition(i)}
              title={`${cellInfo.name} - ${cellInfo.description}`}
            >
              <span className="cell-icon">{cellInfo.icon}</span>
              <span className="cell-number">{i + 1}</span>
              {playersOnCell.length > 0 && (
                <div className="cell-players">
                  {playersOnCell.map((player) => {
                    const playerIndex = getPlayerIndex(player.player_id)
                    const isCurrentPlayer = player.player_id === currentPlayerId
                    return (
                      <div
                        key={player.player_id}
                        className={`player-token ${isCurrentPlayer ? 'current' : ''}`}
                        style={{ 
                          backgroundColor: PLAYER_COLORS[playerIndex % PLAYER_COLORS.length],
                          zIndex: isCurrentPlayer ? 10 : 1
                        }}
                        title={`${player.player_name}${isCurrentPlayer ? ' (你)' : ''}`}
                      >
                        {player.player_name.charAt(0).toUpperCase()}
                      </div>
                    )
                  })}
                </div>
              )}
            </div>
          )
        })}

        {/* 中央資訊面板 */}
        <div className="board-center">
          <h3>☁️ AWS Learning Game</h3>
          <p className="board-subtitle">學習 AWS 架構決策</p>
          
          <div className="legend">
            <h4>格子類型</h4>
            {Object.entries(CELL_TYPES).map(([type, info]) => (
              <div key={type} className="legend-item">
                <span 
                  className={`legend-color cell-${type}`}
                  style={{ backgroundColor: info.color }}
                >
                  {info.icon}
                </span>
                <span className="legend-text">{info.name}</span>
              </div>
            ))}
          </div>

          {players.length > 0 && (
            <div className="player-legend">
              <h4>玩家</h4>
              {players.map((player, index) => (
                <div key={player.player_id} className="legend-item">
                  <span 
                    className="player-color"
                    style={{ backgroundColor: PLAYER_COLORS[index % PLAYER_COLORS.length] }}
                  >
                    {player.player_name.charAt(0).toUpperCase()}
                  </span>
                  <span className="legend-text">
                    {player.player_name}
                    {player.player_id === currentPlayerId && ' (你)'}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default GameBoard
