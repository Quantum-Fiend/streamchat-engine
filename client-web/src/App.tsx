import { useEffect, useRef, useState } from 'react'
import { useChatStore } from './store'
import './App.css'

function App() {
  const {
    messages,
    isConnected,
    room,
    connect,
    sendMessage,
    setRoom
  } = useChatStore()

  const [inputValue, setInputValue] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    connect(room)
  }, []) // Initial connection

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSend = () => {
    if (!inputValue.trim()) return
    sendMessage(inputValue)
    setInputValue('')
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') handleSend()
  }

  return (
    <div className="app-container">
      <header className="chat-header">
        <div className="header-info">
          <h1>ClusterTalk</h1>
          <span className={`status-dot ${isConnected ? 'online' : 'offline'}`}></span>
          <span className="status-text">{isConnected ? `LIVE @ ${room.toUpperCase()}` : 'OFFLINE'}</span>
        </div>
        <div className="room-selector">
          <select
            value={room}
            onChange={(e) => setRoom(e.target.value)}
            className="room-select"
          >
            <option value="general"># General</option>
            <option value="tech"># Tech</option>
            <option value="random"># Random</option>
          </select>
        </div>
      </header>

      <main className="chat-area">
        <div className="message-list">
          {messages.length === 0 && (
            <div className="empty-state">
              <p>Joined channel #{room}.</p>
              <p>No transmissions yet.</p>
            </div>
          )}
          {messages.map((msg, idx) => (
            <div key={idx} className={`message-bubble ${msg.sender.includes('User-') ? 'other' : 'me'}`}>
              <div className="message-sender">{msg.sender}</div>
              <div className="message-content">{msg.payload}</div>
              <div className="message-time">{new Date(msg.timestamp * 1000).toLocaleTimeString()}</div>
            </div>
          ))}
          <div ref={messagesEndRef} />
        </div>
      </main>

      <footer className="chat-input-area">
        <div className="input-wrapper">
          <input
            type="text"
            className="chat-input"
            placeholder={`Message #${room}...`}
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyDown={handleKeyDown}
            disabled={!isConnected}
          />
          <button className="send-btn" onClick={handleSend} disabled={!isConnected}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <line x1="22" y1="2" x2="11" y2="13"></line>
              <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
            </svg>
          </button>
        </div>
      </footer>
    </div>
  )
}

export default App
