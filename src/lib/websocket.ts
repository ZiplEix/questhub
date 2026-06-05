import { writable } from 'svelte/store';
import { supabase } from './supabaseClient';

export const websocketStore = writable<{
    connected: boolean;
    messages: any[];
}>({
    connected: false,
    messages: []
});

let channel: any = null;

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
}
