import Foundation

class WebSocketManager: ObservableObject {
    private var webSocketTask: URLSessionWebSocketTask?
    
    @Published var isConnected = false
    @Published var messages: [String] = []
    
    // Localhost on iOS Simulator
    private let url = URL(string: "ws://localhost:8080/ws?room=general")!
    
    func connect() {
        let session = URLSession(configuration: .default)
        webSocketTask = session.webSocketTask(with: url)
        webSocketTask?.resume()
        
        receiveMessage()
        DispatchQueue.main.async { self.isConnected = true }
    }
    
    private func receiveMessage() {
        webSocketTask?.receive { [weak self] result in
            switch result {
            case .failure(let error):
                print("Error in receiving message: \(error)")
                DispatchQueue.main.async { self?.isConnected = false }
            case .success(let message):
                switch message {
                case .string(let text):
                    DispatchQueue.main.async {
                        self?.messages.append(text)
                        // Recursively listen for next message
                        self?.receiveMessage()
                    }
                case .data(_):
                    break
                @unknown default:
                    break
                }
            }
        }
    }
    
    func sendMessage(_ text: String, room: String) {
        let json: [String: Any] = [
            "type": "message",
            "payload": text,
            "room_id": room,
            "sender": "iOSUser"
        ]
        
        if let data = try? JSONSerialization.data(withJSONObject: json),
           let jsonString = String(data: data, encoding: .utf8) {
            let message = URLSessionWebSocketTask.Message.string(jsonString)
            webSocketTask?.send(message) { error in
                if let error = error {
                    print("WebSocket sending error: \(error)")
                }
            }
        }
    }
    
    func disconnect() {
        webSocketTask?.cancel(with: .goingAway, reason: nil)
        isConnected = false
    }
}
