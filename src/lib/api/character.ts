import { supabase } from '../supabaseClient';

export interface Character {
    id: string;
    name: string;
    avatar_url: string | null;
    stats: any;
    inventory: any;
    is_npc: boolean;
    created_at: string;
    race: string;
    max_hp: number;
    current_hp: number;
    money: number;
    initiative: number;
    age?: string;
    height?: string;
    weight?: string;
    max_spells?: number;
    spells?: any;
    abilities?: string;
    experience?: number;
    type: string; // 'PLAYER', 'NPC', 'MONSTER', 'GM_HIDDEN'
    sub_race?: string | null;
    armor_class: number;
    speed: number;
    game_id: string;
    user_id?: string | null;
}

export async function fetchCharacters(gameId: string): Promise<Character[]> {
    const { data, error } = await supabase
        .from('game_characters')
        .select('*, characters(*)')
        .eq('game_id', gameId);

    if (error) throw error;

    return (data || []).map((gc: any) => {
        const charObj = Array.isArray(gc.characters) ? gc.characters[0] : gc.characters;
        if (!charObj || charObj.type === 'MONSTER' || charObj.type === 'GM_HIDDEN') return null;
        return {
            ...charObj,
            game_id: gc.game_id,
            user_id: gc.user_id
        };
    }).filter((c: any) => c !== null) as Character[];
}

export async function fetchCharacter(gameId: string, characterId: string): Promise<Character> {
    const { data: charData, error: charError } = await supabase
        .from('characters')
        .select('*')
        .eq('id', characterId)
        .single();

    if (charError) throw charError;

    const { data: gc } = await supabase
        .from('game_characters')
        .select('game_id, user_id')
        .eq('game_id', gameId)
        .eq('character_id', characterId)
        .single();

    return {
        ...charData,
        game_id: gc?.game_id || gameId,
        user_id: gc?.user_id || null
    };
}

export async function createCharacter(gameId: string, payload: any): Promise<Character> {
    const { data: charId, error: charErr } = await supabase.rpc('create_character_rpc', {
        p_game_id: gameId,
        p_user_id: payload.user_id || null,
        p_name: payload.name,
        p_race: payload.race || '',
        p_max_hp: payload.max_hp || 10,
        p_is_npc: payload.is_npc || false,
        p_avatar_url: payload.avatar_url || null,
        p_stats: payload.stats || {},
        p_inventory: payload.inventory || [],
        p_money: payload.money || 0,
        p_initiative: payload.initiative || 0,
        p_type: payload.type || 'PLAYER',
        p_sub_race: payload.sub_race || null,
        p_armor_class: payload.armor_class || 10,
        p_speed: payload.speed || 30
    });

    if (charErr) throw charErr;

    // Fetch the newly created character
    const { data: char, error: fetchErr } = await supabase
        .from('characters')
        .select('*')
        .eq('id', charId)
        .single();

    if (fetchErr) throw fetchErr;
    return char;
}

export async function updateCharacter(gameId: string, characterId: string, payload: any): Promise<Character> {
    const { data, error } = await supabase
        .from('characters')
        .update(payload)
        .eq('id', characterId)
        .select()
        .single();

    if (error) throw error;
    return data;
}

export async function deleteCharacter(gameId: string, characterId: string): Promise<void> {
    const { error } = await supabase
        .from('characters')
        .delete()
        .eq('id', characterId);

    if (error) throw error;
}

export async function assignCharacter(gameId: string, characterId: string, userId: string | null): Promise<void> {
    const { error } = await supabase
        .from('game_characters')
        .update({ user_id: userId || null })
        .eq('game_id', gameId)
        .eq('character_id', characterId);

    if (error) throw error;
}
