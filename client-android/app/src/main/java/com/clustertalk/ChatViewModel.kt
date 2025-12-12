package com.clustertalk.ui

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.clustertalk.network.WebSocketManager

class ChatViewModel : ViewModel(), WebSocketManager.MessageListener {
    private val webSocketManager = WebSocketManager(this)
    
    private val _messages = MutableLiveData<List<String>>(emptyList())
    val messages: LiveData<List<String>> = _messages
    
    private val _isConnected = MutableLiveData(false)
    val isConnected: LiveData<Boolean> = _isConnected

    fun connect() {
        webSocketManager.connect()
    }
    
    fun sendMessage(text: String) {
        webSocketManager.sendMessage(text, "general")
    }

    override fun onConnect() {
        _isConnected.postValue(true)
    }

    override fun onMessage(text: String) {
        val currentList = _messages.value.orEmpty().toMutableList()
        currentList.add(text)
        _messages.postValue(currentList)
    }

    override fun onDisconnect() {
        _isConnected.postValue(false)
    }

    override fun onCleared() {
        super.onCleared()
        webSocketManager.close()
    }
}
