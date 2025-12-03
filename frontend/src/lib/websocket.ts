import { writable } from 'svelte/store';
import { PUBLIC_BASE_WS_URL, PUBLIC_BASE_API_URL } from '$env/static/public';

export const websocketStore = writable<{
    connected: boolean;
    messages: any[];
}>({
    connected: false,
    messages: []
});

let socket: WebSocket | null = null;
let reconnectTimer: ReturnType<typeof setTimeout> | null = null;

export function sendMessage(message: any) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
    } else {
        console.error("WebSocket is not connected");
    }
}

export async function fetchHistory(gameId: string, token: string) {
    try {
        const response = await fetch(`${PUBLIC_BASE_API_URL}/table/${gameId}/chat`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        if (response.ok) {
            const history = await response.json();
            if (Array.isArray(history)) {
                websocketStore.update(store => ({
                    ...store,
                    messages: [...history, ...store.messages] // Prepend history
                }));
            }
        }
    } catch (e) {
        console.error("Failed to fetch chat history:", e);
    }
}

export function connectWebSocket(token: string) {
    if (socket?.readyState === WebSocket.OPEN) return;

    const url = `${PUBLIC_BASE_WS_URL}/ws?token=${token}`;
    socket = new WebSocket(url);

    socket.onopen = () => {
        console.log('WebSocket connected');
        websocketStore.update(s => ({ ...s, connected: true }));
        if (reconnectTimer) {
            clearTimeout(reconnectTimer);
            reconnectTimer = null;
        }
    };

    socket.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            websocketStore.update(s => ({
                ...s,
                messages: [...s.messages, data]
            }));
        } catch (e) {
            console.error('Failed to parse WebSocket message:', e);
        }
    };

    socket.onclose = () => {
        console.log('WebSocket disconnected');
        websocketStore.update(s => ({ ...s, connected: false }));
        socket = null;
        // Reconnect after 3 seconds
        reconnectTimer = setTimeout(() => connectWebSocket(token), 3000);
    };

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        socket?.close();
    };
}

export function closeWebSocket() {
    if (socket) {
        socket.close();
        socket = null;
    }
    if (reconnectTimer) {
        clearTimeout(reconnectTimer);
        reconnectTimer = null;
    }
}
