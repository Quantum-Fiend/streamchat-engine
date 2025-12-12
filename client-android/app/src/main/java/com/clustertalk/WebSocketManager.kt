package com.clustertalk.network

import okhttp3.*
import org.json.JSONObject
import java.util.concurrent.TimeUnit

class WebSocketManager(private val listener: MessageListener) {
    private val client = OkHttpClient.Builder()
        .readTimeout(0, TimeUnit.MILLISECONDS)
        .build()
    
    private var webSocket: WebSocket? = null
    
    // Switch between 10.0.2.2 (Emulator) or localhost if tunneling
    private val serverUrl = "ws://10.0.2.2:8080/ws" 

    interface MessageListener {
        fun onConnect()
        fun onMessage(text: String)
        fun onDisconnect()
    }

    fun connect(room: String = "general") {
        val request = Request.Builder()
            .url("$serverUrl?room=$room")
            .build()
        
        webSocket = client.newWebSocket(request, object : WebSocketListener() {
            override fun onOpen(webSocket: WebSocket, response: Response) {
                listener.onConnect()
            }

            override fun onMessage(webSocket: WebSocket, text: String) {
                listener.onMessage(text)
            }

            override fun onClosing(webSocket: WebSocket, code: Int, reason: String) {
                listener.onDisconnect()
            }
            
            override fun onFailure(webSocket: WebSocket, t: Throwable, response: Response?) {
                listener.onDisconnect()
            }
        })
    }

    fun sendMessage(text: String, room: String) {
        val json = JSONObject().apply {
            put("type", "message")
            put("payload", text)
            put("room_id", room)
            put("sender", "AndroidUser")
        }
        webSocket?.send(json.toString())
    }

    fun close() {
        webSocket?.close(1000, "User Exit")
    }
}
