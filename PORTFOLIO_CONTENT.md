# Portfolio: ClusterTalk - Distributed Chat System

## Project Overview
**ClusterTalk** is a high-performance distributed messaging platform designed to solve the challenges of real-time communication at scale. It leverages a microservices architecture to separate concerns between real-time routing (Go), intelligent moderation (Python), and high-throughput analytics (Rust).

## 🏗️ Architecture Design
The system uses a **Hexagonal Architecture** approach:
- **Core Engine (Go)**: Handles WebSocket connections, broadcasting, and room state management using direct memory manipulation and channel-based concurrency patterns.
- **Intelligence Layer (Python)**: Decoupled service for NLP tasks. Implemented this way to leverage Python's rich ML ecosystem without blocking the Go event loop.
- **Stats Layer (Rust)**: Utilized Rust for the analytics consumer to ensure memory safety handling high-volume event streams with zero garbage collection pauses.

## 🔥 Key Technical Challenges & Solutions

### 1. Synchronous vs Asynchronous Processing
*Challenge*: The AI Moderation needed to be strict (blocking), while Analytics could be eventual.
*Solution*: Implemented a synchronous HTTP blocking call for the AI check within the message pipeline, but spawned a separate lightweight goroutine for the fire-and-forget call to the Rust Analytics service, preventing stats logging from adding latency to the user experience.

### 2. Protocol Design in Constraints
*Challenge*: Need for structured communication without heavy overhead.
*Solution*: Designed a strict JSON-based protocol with `type`, `payload`, and `room_id` routing. This provided flexibility for the frontend while maintaining strict types in the Go backend.

### 3. Cross-Platform State Management
*Challenge*: Keeping logic consistent across Web, Android, and iOS.
*Solution*: Standardized the WebSocket manager implementation across all three platforms. Created a shared "Protocol Spec" that guided the implementation of the `OkHttp` (Android) and `URLSession` (iOS) networking layers.

### 4. Zero-Dependency Persistence
*Challenge*: Saving history without requiring users to install Postgres/Docker.
*Solution*: Integrated `go-sqlite` (pure Go implementation) to provide robust SQL persistence in a single executable binary, automatically managing connection pooling and schema migrations on startup.


## 🚀 Impact & Results
- **Latency**: Achieved sub-50ms message delivery latency.
- **Safety**: Automated protection against toxicity with >90% blocking rate for flagged keywords.
- **Scalability**: Architecture allows horizontal scaling of the Go Hub by introducing a Redis Pub/Sub layer (Architecture Ready).

## 💻 Tech Stack
- **Languages**: Go, Python, Rust, TypeScript, Kotlin, Swift
- **Infrastructure**: Docker, WebSockets
- **Tools**: Vite, FastAPI, Actix
