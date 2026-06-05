import { supabase } from '../supabaseClient';

export interface GameBoard {
    id: string;
    game_id: string;
    name: string;
    image_url: string | null;
    is_active: boolean;
    created_at: string;
}

export async function fetchBoards(gameId: string): Promise<GameBoard[]> {
    const { data, error } = await supabase
        .from('game_boards')
        .select('*')
        .eq('game_id', gameId)
        .order('created_at', { ascending: true });

    if (error) throw error;
    return data || [];
}

export async function addBoard(gameId: string, name: string, imageUrl: string): Promise<GameBoard> {
    const { data, error } = await supabase
        .from('game_boards')
        .insert({
            game_id: gameId,
            name,
            image_url: imageUrl,
            is_active: false
        })
        .select()
        .single();

    if (error) throw error;
    return data;
}

export async function deleteBoard(boardId: string): Promise<void> {
    const { error } = await supabase
        .from('game_boards')
        .delete()
        .eq('id', boardId);

    if (error) throw error;
}

export async function activateBoard(gameId: string, boardId: string): Promise<void> {
    // 1. Désactiver tous les plateaux de la partie
    const { error: deactivateError } = await supabase
        .from('game_boards')
        .update({ is_active: false })
        .eq('game_id', gameId);

    if (deactivateError) throw deactivateError;

    // 2. Activer le plateau sélectionné
    const { error: activateError } = await supabase
        .from('game_boards')
        .update({ is_active: true })
        .eq('id', boardId);

    if (activateError) throw activateError;
}
