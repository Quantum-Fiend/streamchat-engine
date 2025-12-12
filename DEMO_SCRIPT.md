# Demo Video Recording Script

**Goal**: Show off the system in 2 minutes or less.

## Setup
1. Open **3 windows** on your screen:
   - Window 1: `http://localhost:5173` (Chrome/Edge) - **Client A**
   - Window 2: `http://localhost:5173` (Incognito Mode) - **Client B**
   - Window 3: Terminal showing the server logs (Split screen: top=Go Backend, bottom=Python AI).

## The Script

### Scene 1: Introduction (0:00 - 0:15)
- **Action**: Start recording. Show the clean "ClusterTalk" UI.
- **Narration**: "This is ClusterTalk, a distributed chat system I built using Go, Python, and Rust. It features real-time messaging, AI moderation, and multi-room support."

### Scene 2: Real-time Speed (0:15 - 0:45)
- **Action**: Type "Hello from Client A" in the first window. Hit Enter.
- **Action**: Immediately move mouse to Window 2.
- **Narration**: "As you can see, message delivery is instant thanks to the Go-based WebSocket hub."
- **Action**: Reply from Client B: "Loud and clear!".

### Scene 3: Multi-Room Feature (0:45 - 1:15)
- **Action**: On Client A, use the dropdown to switch to `#tech`.
- **Action**: Type "Posting in Tech channel" in Client A.
- **Observation**: Point out that Client B (still in `#general`) did *not* see it.
- **Narration**: "The backend handles logical partitioning using a room management system, ensuring traffic isolation."

### Scene 4: AI Moderation (1:15 - 1:45)
- **Action**: On Client B, type: "This is spam and bad content".
- **Action**: Hit Send.
- **Observation**: The message appears as `[CENSORED BY AI]: **********`.
- **Narration**: "The Python microservice analyzes every message in real-time. It successfully intercepted and censored the toxic content before broadcast."

### Scene 5: Conclusion (1:45 - 2:00)
- **Action**: Briefly flip to the VS Code Architecture Diagram or the Terminal logs showing the traffic.
- **Narration**: "The system is also logging analytics to a Rust service in the background. Thanks for watching."
- **Action**: Stop Recording.
