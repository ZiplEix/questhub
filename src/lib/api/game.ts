import { supabase } from '../supabaseClient';

export interface Game {
    id: string;
    name: string;
    gm_id: string;
    gm_name?: string;
    invite_code: string;
    is_active: boolean;
    image_url: string | null;
    state: string;
    created_at: string;
    is_gm?: boolean;
    current_character_id?: string | null;
}

export async function fetchGames(): Promise<Game[]> {
    const { data, error } = await supabase
        .from('user_games')
        .select('*')
        .order('created_at', { ascending: false });

    if (error) throw error;
    return data || [];
}

export async function fetchGame(gameId: string): Promise<Game> {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Unauthorized');

    const { data: game, error: gameError } = await supabase
        .from('games')
        .select('*')
        .eq('id', gameId)
        .single();

    if (gameError) throw gameError;

    const isGM = game.gm_id === user.id;
    let currentCharacterId: string | null = null;

    if (isGM) {
        // Check if GM character exists (name "GM Notes", type "GM_HIDDEN")
        const { data: gmChars } = await supabase
            .from('game_characters')
            .select('*, characters!inner(*)')
            .eq('game_id', gameId)
            .eq('user_id', user.id)
            .eq('characters.type', 'GM_HIDDEN');

        if (gmChars && gmChars.length > 0) {
            currentCharacterId = gmChars[0].character_id;
        } else {
            // Create a new character for GM via RPC
            const { data: charId, error: charErr } = await supabase.rpc('create_character_rpc', {
                p_game_id: gameId,
                p_user_id: user.id,
                p_name: 'GM Notes',
                p_race: 'Game Master',
                p_max_hp: 1,
                p_is_npc: false,
                p_avatar_url: null,
                p_stats: {},
                p_inventory: [],
                p_money: 0,
                p_initiative: 0,
                p_type: 'GM_HIDDEN',
                p_sub_race: null,
                p_armor_class: 10,
                p_speed: 30
            });

            if (charErr) throw charErr;
            currentCharacterId = charId;
        }
    } else {
        // Regular player character
        const { data: playerChars } = await supabase
            .from('game_characters')
            .select('character_id')
            .eq('game_id', gameId)
            .eq('user_id', user.id);

        if (playerChars && playerChars.length > 0) {
            currentCharacterId = playerChars[0].character_id;
        }
    }

    return {
        ...game,
        is_gm: isGM,
        current_character_id: currentCharacterId
    };
}

export async function createGame(name: string, imageUrl?: string): Promise<Game> {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Unauthorized');

    const inviteCode = Math.random().toString(36).substring(2, 7).toUpperCase();
    const { data, error } = await supabase
        .from('games')
        .insert({
            name,
            image_url: imageUrl || null,
            gm_id: user.id,
            invite_code: inviteCode,
            state: 'paused'
        })
        .select()
        .single();

    if (error) throw error;
    return data;
}

export async function updateGame(gameId: string, payload: { name?: string; image_url?: string; is_active?: boolean; state?: string }): Promise<Game> {
    const { data, error } = await supabase
        .from('games')
        .update(payload)
        .eq('id', gameId)
        .select()
        .single();

    if (error) throw error;
    return data;
}

export async function deleteGame(gameId: string): Promise<void> {
    const { error } = await supabase
        .from('games')
        .delete()
        .eq('id', gameId);

    if (error) throw error;
}

export async function joinGame(inviteCode: string): Promise<{ id: string }> {
    const { data, error } = await supabase.rpc('join_game_by_invite_code', {
        p_invite_code: inviteCode
    });

    if (error) throw error;
    return { id: data };
}

export async function regenerateInviteCode(gameId: string): Promise<string> {
    const inviteCode = Math.random().toString(36).substring(2, 7).toUpperCase();
    const { data, error } = await supabase
        .from('games')
        .update({ invite_code: inviteCode })
        .eq('id', gameId)
        .select('invite_code')
        .single();

    if (error) throw error;
    return data.invite_code;
}
