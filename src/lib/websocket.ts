import { writable } from 'svelte/store';
import { supabase } from './supabaseClient';
import { fetchBoards } from './api/board';

export const websocketStore = writable<{
    connected: boolean;
    messages: any[];
}>({
    connected: false,
    messages: []
});

export const boardsStore = writable<any[]>([]);
export const activeBoardStore = writable<any | null>(null);

let channel: any = null;

function syncActiveBoard(boards: any[]) {
    const activeBoard = boards.find(b => b.is_active);
    activeBoardStore.set(activeBoard || null);
}

export async function fetchInitialBoards(gameId: string) {
    try {
        const boards = await fetchBoards(gameId);
        boardsStore.set(boards);
        syncActiveBoard(boards);
    } catch (e) {
        console.error("Failed to fetch initial boards:", e);
    }
}

export async function fetchHistory(gameId: string) {
    try {
        const { data, error } = await supabase
            .from('messages')
            .select('*')
            .eq('game_id', gameId)
            .order('created_at', { ascending: true });

        if (error) throw error;

        websocketStore.update(store => ({
            ...store,
            messages: data || []
        }));
    } catch (e) {
        console.error("Failed to fetch chat history:", e);
    }
}

export async function connectWebSocket(gameId: string) {
    if (channel) {
        if (channel.topic === `realtime:game:${gameId}`) return;
        closeWebSocket();
    }

    websocketStore.update(s => ({ ...s, connected: false }));

    // Wait for the session to be loaded so the channel is authenticated
    await supabase.auth.getSession();

    // Load initial campaign data
    await fetchInitialBoards(gameId);

    // Subscribe to messages channel for the specific game
    channel = supabase.channel(`game:${gameId}`)
        .on(
            'postgres_changes',
            {
                event: 'INSERT',
                schema: 'public',
                table: 'messages',
                filter: `game_id=eq.${gameId}`
            },
            (payload) => {
                websocketStore.update(s => {
                    // Prevent duplicate rendering
                    if (s.messages.some(m => m.id === payload.new.id)) return s;
                    return {
                        ...s,
                        messages: [...s.messages, payload.new]
                    };
                });
            }
        )
        .on(
            'postgres_changes',
            {
                event: 'DELETE',
                schema: 'public',
                table: 'messages',
                filter: `game_id=eq.${gameId}`
            },
            (payload) => {
                websocketStore.update(s => ({
                    ...s,
                    messages: s.messages.filter(m => m.id !== payload.old.id)
                }));
            }
        )
        .on(
            'postgres_changes',
            {
                event: '*',
                schema: 'public',
                table: 'game_boards',
                filter: `game_id=eq.${gameId}`
            },
            (payload) => {
                boardsStore.update(boards => {
                    let updatedBoards = [...boards];
                    if (payload.eventType === 'INSERT') {
                        updatedBoards.push(payload.new);
                    } else if (payload.eventType === 'UPDATE') {
                        updatedBoards = updatedBoards.map(b => b.id === payload.new.id ? payload.new : b);
                    } else if (payload.eventType === 'DELETE') {
                        updatedBoards = updatedBoards.filter(b => b.id !== payload.old.id);
                    }
                    syncActiveBoard(updatedBoards);
                    return updatedBoards;
                });
            }
        )
        .on(
            'broadcast',
            { event: 'ping' },
            (payload) => {
                console.log('Received broadcast ping:', payload);
                if (payload && payload.payload) {
                    pingCallbacks.forEach(cb => cb(payload.payload));
                }
            }
        )
        .on(
            'broadcast',
            { event: 'token_drag' },
            (payload) => {
                if (payload && payload.payload) {
                    tokenDragCallbacks.forEach(cb => cb(payload.payload));
                }
            }
        )
        .subscribe((status) => {
            if (status === 'SUBSCRIBED') {
                console.log('Supabase Realtime connected');
                websocketStore.update(s => ({ ...s, connected: true }));
            } else if (status === 'CLOSED' || status === 'CHANNEL_ERROR') {
                console.log('Supabase Realtime disconnected');
                websocketStore.update(s => ({ ...s, connected: false }));
            }
        });
}

export function closeWebSocket() {
    if (channel) {
        supabase.removeChannel(channel);
        channel = null;
    }
    websocketStore.update(s => ({ ...s, connected: false, messages: [] }));
    boardsStore.set([]);
    activeBoardStore.set(null);
}

let pingCallbacks: ((ping: { x: number; y: number; id: number; color?: string }) => void)[] = [];

export function onPingReceived(callback: (ping: { x: number; y: number; id: number; color?: string }) => void) {
    pingCallbacks.push(callback);
    return () => {
        pingCallbacks = pingCallbacks.filter(cb => cb !== callback);
    };
}

export function sendPing(x: number, y: number, color: string) {
    if (channel) {
        channel.send({
            type: 'broadcast',
            event: 'ping',
            payload: { x, y, id: Date.now(), color }
        });
    }
}

let tokenDragCallbacks: ((drag: { tokenId: string; x: number; y: number }) => void)[] = [];

export function onTokenDragged(callback: (drag: { tokenId: string; x: number; y: number }) => void) {
    tokenDragCallbacks.push(callback);
    return () => {
        tokenDragCallbacks = tokenDragCallbacks.filter(cb => cb !== callback);
    };
}

export function sendTokenDrag(tokenId: string, x: number, y: number) {
    if (channel) {
        channel.send({
            type: 'broadcast',
            event: 'token_drag',
            payload: { tokenId, x, y }
        });
    }
}



