import { supabase } from '../supabaseClient';

export interface Invitation {
    id: string;
    game_id: string;
    user_id: string;
    user_name: string;
    created_at: string;
}

export async function fetchInvitations(gameId: string): Promise<Invitation[]> {
    const { data, error } = await supabase
        .from('game_invitation_users')
        .select('*')
        .eq('game_id', gameId);

    if (error) throw error;
    return data || [];
}

export async function acceptInvitation(gameId: string, userId: string): Promise<void> {
    // 1. Add to players
    const { error: insertErr } = await supabase
        .from('game_players')
        .insert({ game_id: gameId, user_id: userId });

    if (insertErr) throw insertErr;

    // 2. Delete invitation
    const { error: deleteErr } = await supabase
        .from('game_invitations')
        .delete()
        .eq('game_id', gameId)
        .eq('user_id', userId);

    if (deleteErr) throw deleteErr;
}

export async function declineInvitation(gameId: string, userId: string): Promise<void> {
    const { error } = await supabase
        .from('game_invitations')
        .delete()
        .eq('game_id', gameId)
        .eq('user_id', userId);

    if (error) throw error;
}
