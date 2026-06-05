import { supabase } from '../supabaseClient';

export interface RollResult {
    result: number;
    sides: number;
}

export async function rollDice(gameId: string, sides: number = 20, secret: boolean = false): Promise<RollResult> {
    const { data, error } = await supabase.rpc('roll_dice', {
        p_game_id: gameId,
        p_sides: sides,
        p_is_secret: secret
    });

    if (error) throw error;
    return { result: data, sides };
}
