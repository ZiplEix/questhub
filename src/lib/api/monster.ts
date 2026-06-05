import { supabase } from '../supabaseClient';
import type { Character } from './character';

export async function fetchMonsters(gameId: string): Promise<Character[]> {
    const { data, error } = await supabase
        .from('game_characters')
        .select('*, characters(*)')
        .eq('game_id', gameId);

    if (error) throw error;

    return (data || []).map((gc: any) => {
        const charObj = Array.isArray(gc.characters) ? gc.characters[0] : gc.characters;
        if (!charObj || charObj.type !== 'MONSTER') return null;
        return {
            ...charObj,
            game_id: gc.game_id,
            user_id: gc.user_id
        };
    }).filter((c: any) => c !== null) as Character[];
}
