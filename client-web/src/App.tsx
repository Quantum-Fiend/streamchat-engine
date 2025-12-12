import { useState, useEffect, useRef } from 'react'
import './App.css'

interface WSMessage {
  type: string
  payload: string
  sender: string
  room_id: string
  timestamp: number
}

function App() {
  const [messages, setMessages] = useState<WSMessage[]>([])
  const [inputValue, setInputValue] = useState('')
  const [isConnected, setIsConnected] = useState(false)
  const [room, setRoom] = useState('general')

  const ws = useRef<WebSocket | null>(null)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    // Connect to WebSocket with Room Query Param
    const socket = new WebSocket(`ws://localhost:8080/ws?room=${room}`)
    ws.current = socket

    socket.onopen = () => {
      console.log(`Connected to ClusterTalk Room: ${room}`)
      setIsConnected(true)
    }

    socket.onclose = () => {
      console.log('Disconnected from WebSocket')
      setIsConnected(false)
    }

    socket.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data)
        setMessages((prev) => [...prev, msg])
      } catch (e) {
        console.error("Failed to parse message", event.data)
      }
    }

    return () => {
      socket.close()
    }
  }, [room]) // Reconnect if room changes

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const sendMessage = () => {
    if (!inputValue.trim() || !ws.current) return

    const msg = {
      type: "message",
      payload: inputValue,
      room_id: room,
      // Sender ID is handled by server for security usually, but we send it for now if needed, 
      // though server overrides it in current implementation.
      sender: "me"
    }

    ws.current.send(JSON.stringify(msg))
    setInputValue('')
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') sendMessage()
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
            onChange={(e) => {
              setMessages([]); // Clear chat on room switch
              setRoom(e.target.value);
            }}
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
          <button className="send-btn" onClick={sendMessage} disabled={!isConnected}>
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
