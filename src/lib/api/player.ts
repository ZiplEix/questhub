import { supabase } from '../supabaseClient';

export interface Player {
    game_id: string;
    user_id: string;
    name: string;
    email: string;
    joined_at: string;
}

export async function fetchPlayers(gameId: string): Promise<Player[]> {
    const { data, error } = await supabase
        .from('game_player_users')
        .select('*')
        .eq('game_id', gameId);

    if (error) throw error;
    return data || [];
}

export async function kickPlayer(gameId: string, userId: string): Promise<void> {
    const { error } = await supabase
        .from('game_players')
        .delete()
        .eq('game_id', gameId)
        .eq('user_id', userId);

    if (error) throw error;
}
