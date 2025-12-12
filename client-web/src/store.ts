import { create } from 'zustand'

interface WSMessage {
    type: string
    payload: string
    sender: string
    room_id: string
    timestamp: number
}

interface ChatState {
    messages: WSMessage[]
    isConnected: boolean
    room: string
    socket: WebSocket | null

    // Actions
    connect: (room: string) => void
    disconnect: () => void
    sendMessage: (text: string) => void
    setRoom: (room: string) => void
    addMessage: (msg: WSMessage) => void
}

export const useChatStore = create<ChatState>((set, get) => ({
    messages: [],
    isConnected: false,
    room: 'general',
    socket: null,

    connect: (room: string) => {
        const { socket } = get()
        if (socket) {
            socket.close()
        }

        const ws = new WebSocket(`ws://localhost:8080/ws?room=${room}`)

        ws.onopen = () => set({ isConnected: true })
        ws.onclose = () => set({ isConnected: false, socket: null })

        ws.onmessage = (event) => {
            try {
                const msg: WSMessage = JSON.parse(event.data)
                get().addMessage(msg)
            } catch (e) {
                console.error("Failed to parse message", event.data)
            }
        }

        set({ socket: ws, room })
    },

    disconnect: () => {
        const { socket } = get()
        if (socket) socket.close()
        set({ socket: null, isConnected: false })
    },

    sendMessage: (text: string) => {
        const { socket, room } = get()
        if (!socket || !text.trim()) return

        const msg = {
            type: "message",
            payload: text,
            room_id: room,
            sender: "me"
        }
        socket.send(JSON.stringify(msg))
    },

    setRoom: (room: string) => {
        // Create new connection on room switch
        get().connect(room)
        // Clear messages
        set({ messages: [] })
    },

    addMessage: (msg) => set((state) => ({ messages: [...state.messages, msg] })),
}))
