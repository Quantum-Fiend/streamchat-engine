# ClusterTalk - Distributed Chat System (Advanced Edition)

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/backend-Go-00ADD8.svg)
![React](https://img.shields.io/badge/frontend-React-61DAFB.svg)
![Python](https://img.shields.io/badge/AI-Python-3776AB.svg)
![Rust](https://img.shields.io/badge/Analytics-Rust-000000.svg)

**ClusterTalk** is a production-grade, distributed real-time messaging platform engineered with a microservices architecture. It demonstrates mastery of **distributed systems**, **cross-platform engineering**, and **applied AI integration**.

Designed for high-throughput and low-latency, it seamlessly orchestrates a **Go** websocket engine, **Python** AI moderation service, **Rust** analytics consumer, and **React/Mobile** clients.

---

## 🏗️ System Architecture

```mermaid
graph TD
    subgraph Clients
        Web[React Web Client]
        Android[Android App (Kotlin)]
        iOS[iOS App (Swift)]
    end

    subgraph "Core Cluster"
        LB[Load Balancer / Gateway]
        Hub[Go WebSocket Hub]
    end

    subgraph "Microservices"
        AI[Python AI Service]
        Analytics[Rust Analytics Engine]
    end

    Web <-->|WebSocket| Hub
    Android <-->|WebSocket| Hub
    iOS <-->|WebSocket| Hub
    
    Hub -->|HTTP Sync| AI
    Hub -.->|HTTP Async| Analytics
```

---

## ✨ Key Features

- **🚀 High-Performance Core**: Go-based event loop using lightweight goroutines handling thousands of concurrent connections.
- **🛡️ AI-Powered Moderation**: Real-time toxicity detection using a Python Microservice. Automatically censors spam and offensive content before broadcast.
- **📊 Rust Analytics**: Zero-overhead event ingestion pipeline implemented in Rust for analyzing chat throughput.
- **🌐 Multi-Room Clustering**: logical partitioning of chat streams (`#general`, `#tech`, `#random`).
- **📱 Cross-Platform Architecture**: Unified protocol supporting Web, Android, and iOS clients.
- **🎨 Cyberpunk UI**: Modern, glassmorphism-based interface with responsive design.

---

## 🛠️ Technology Stack

| Component | Technology | Role |
|-----------|------------|------|
| **Backend** | **Go (Golang)** | WebSocket Hub, Concurrency, Routing |
| **Frontend** | **React + Vite** | SPA Client, State Management |
| **AI Service** | **Python (FastAPI)** | NLP, Toxicity Detection, Filtering |
| **Analytics** | **Rust (Actix)** | High-speed Event Logging |
| **Mobile** | **Kotlin / Swift** | Native Client Scaffolding |
| **Protocol** | **JSON/WebSockets** | Real-time bidirectional communication |

---

## 🏁 Getting Started

### Prerequisites
- Go 1.19+
- Node.js 16+
- Python 3.9+
- Rust (Cargo)

### 1. Start the Backend Engine
The heart of the system.
```bash
cd backend-go
go mod tidy
go run cmd/server/main.go
# 🟢 Server listening on :8080
```

### 2. Activate AI Intelligence
The brain of the system.
```bash
cd ai-python
# Create virtual env
python -m venv venv
./venv/Scripts/activate
pip install -r requirements.txt
# Start Service
uvicorn main:app --reload --port 8000
# 🟢 AI Service online at :8000
```

### 3. Launch Web Interface
The face of the system.
```bash
cd client-web
npm install
npm run dev
# 🟢 Client ready at http://localhost:5173
```

### 4. (Optional) Run Analytics
The memory of the system.
```bash
cd analytics-rust
cargo run
# 🟢 Analytics ingestion at :9000
```

---

## 🧪 Verification & Demo

1.  **Open Client**: Navigate to `http://localhost:5173`.
2.  **Join Room**: You are automatically placed in `#general`.
3.  **Test Chat**: Send a message. It appears instantly.
4.  **Test AI**: Send the message `"This is spam"`. It will be intercepted and replaced with `[CENSORED BY AI]`.
5.  **Multi-Room**: Open a second tab, switch to `#tech`. Messages are isolated between rooms.

---

## 📂 Directory Structure

```
ClusterTalk/
├── backend-go/         # Go WebSocket Server
├── ai-python/          # Python AI Microservice
├── analytics-rust/     # Rust Data Consumer
├── client-web/         # React Frontend
├── client-android/     # Native Android Project
├── client-ios/         # Native iOS Project
└── infra/              # Docker & Kubernetes Configs
```

---

## ⚖️ License

Distributed under the MIT License. See `LICENSE` for more information.
