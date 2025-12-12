import SwiftUI

struct ContentView: View {
    @StateObject private var wsManager = WebSocketManager()
    @State private var messageText = ""
    
    var body: some View {
        VStack {
            // Header
            HStack {
                Text("ClusterTalk iOS")
                    .font(.headline)
                Spacer()
                Circle()
                    .fill(wsManager.isConnected ? Color.green : Color.red)
                    .frame(width: 10, height: 10)
            }
            .padding()
            
            // Message List
            ScrollView {
                LazyVStack(alignment: .leading) {
                    ForEach(wsManager.messages, id: \.self) { msg in
                        // Basic JSON parsing to show payload only
                        // In real app, we decode properly
                        Text(msg)
                            .padding()
                            .background(Color.blue.opacity(0.1))
                            .cornerRadius(8)
                            .padding(.horizontal)
                    }
                }
            }
            
            // Input Area
            HStack {
                TextField("Message...", text: $messageText)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
                
                Button(action: {
                    wsManager.sendMessage(messageText, room: "general")
                    messageText = ""
                }) {
                    Image(systemName: "paperplane.fill")
                }
            }
            .padding()
        }
        .onAppear {
            wsManager.connect()
        }
    }
}
