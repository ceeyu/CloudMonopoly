import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Lobby from './pages/Lobby'
import Game from './pages/Game'
import './App.css'

function App() {
  return (
    <BrowserRouter>
      <div className="app">
        <header className="app-header">
          <h1>🎮 AWS Learning Game</h1>
          <p>AWS SAA 證照學習遊戲</p>
        </header>
        <main className="app-main">
          <Routes>
            <Route path="/" element={<Lobby />} />
            <Route path="/game/:gameId" element={<Game />} />
          </Routes>
        </main>
        <footer className="app-footer">
          <p>© 2024 AWS Learning Game - 透過遊戲學習 AWS 架構</p>
        </footer>
      </div>
    </BrowserRouter>
  )
}

export default App
