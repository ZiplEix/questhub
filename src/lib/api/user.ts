import { supabase } from '../supabaseClient';

export interface UserProfile {
    id: string;
    email: string | undefined;
    name: string;
}

export interface UserStats {
    games_played: number;
    games_mastered: number;
    characters_created: number;
}

export interface UserCampaign {
    game_id: string;
    game_name: string;
    game_image_url: string;
    character_name: string;
    character_avatar_url: string;
    joined_at: string;
}

export async function fetchUserProfile(): Promise<UserProfile> {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Unauthorized');

    return {
        id: user.id,
        email: user.email,
        name: user.user_metadata?.name || user.email?.split('@')[0] || 'Joueur'
    };
}

export async function fetchUserStats(): Promise<UserStats> {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Unauthorized');

    const { count: gamesPlayed, error: err1 } = await supabase
        .from('game_players')
        .select('*', { count: 'exact', head: true })
        .eq('user_id', user.id);
    if (err1) throw err1;

    const { count: gamesMastered, error: err2 } = await supabase
        .from('games')
        .select('*', { count: 'exact', head: true })
        .eq('gm_id', user.id);
    if (err2) throw err2;

    const { count: charactersCreated, error: err3 } = await supabase
        .from('game_characters')
        .select('*', { count: 'exact', head: true })
        .eq('user_id', user.id);
    if (err3) throw err3;

    return {
        games_played: gamesPlayed || 0,
        games_mastered: gamesMastered || 0,
        characters_created: charactersCreated || 0
    };
}

export async function fetchUserCampaigns(): Promise<UserCampaign[]> {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Unauthorized');

    const { data, error } = await supabase
        .from('game_characters')
        .select('assigned_at, games(id, name, image_url), characters(name, avatar_url)')
        .eq('user_id', user.id)
        .order('assigned_at', { ascending: false });

    if (error) throw error;

    return (data || []).map((gc: any) => {
        const gameObj = Array.isArray(gc.games) ? gc.games[0] : gc.games;
        const charObj = Array.isArray(gc.characters) ? gc.characters[0] : gc.characters;
        return {
            game_id: gameObj?.id || '',
            game_name: gameObj?.name || '',
            game_image_url: gameObj?.image_url || '',
            character_name: charObj?.name || '',
            character_avatar_url: charObj?.avatar_url || '',
            joined_at: gc.assigned_at
        };
    });
}
