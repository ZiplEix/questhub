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

export interface BoardToken {
    id: string;
    board_id: string;
    game_id: string;
    character_id: string;
    x: number;
    y: number;
    created_at: string;
    character?: {
        name: string;
        avatar_url: string | null;
        is_npc: boolean;
        type: string;
        user_id?: string | null;
    };
}

export async function fetchBoardTokens(boardId: string): Promise<BoardToken[]> {
    const { data, error } = await supabase
        .from('board_tokens')
        .select('*, character:characters(*)')
        .eq('board_id', boardId);

    if (error) throw error;

    if (data && data.length > 0) {
        const gameId = data[0].game_id;
        const { data: gcData } = await supabase
            .from('game_characters')
            .select('character_id, user_id')
            .eq('game_id', gameId);

        const assignmentMap = new Map((gcData || []).map(gc => [gc.character_id, gc.user_id]));

        return data.map((t: any) => {
            const charObj = Array.isArray(t.character) ? t.character[0] : t.character;
            return {
                ...t,
                character: charObj ? {
                    ...charObj,
                    user_id: assignmentMap.get(t.character_id) || null
                } : undefined
            };
        });
    }

    return [];
}

export async function addBoardToken(boardId: string, gameId: string, characterId: string, x: number, y: number): Promise<BoardToken> {
    // Check if token already exists for this board and character
    const { data: existing, error: checkError } = await supabase
        .from('board_tokens')
        .select('id')
        .eq('board_id', boardId)
        .eq('character_id', characterId)
        .maybeSingle();

    if (checkError) throw checkError;
    if (existing) {
        throw new Error("Ce personnage est déjà sur le plateau.");
    }

    const { data, error } = await supabase
        .from('board_tokens')
        .insert({
            board_id: boardId,
            game_id: gameId,
            character_id: characterId,
            x,
            y
        })
        .select('*, character:characters(*)')
        .single();

    if (error) throw error;

    const { data: gc } = await supabase
        .from('game_characters')
        .select('user_id')
        .eq('game_id', gameId)
        .eq('character_id', characterId)
        .single();

    const charObj = Array.isArray(data.character) ? data.character[0] : data.character;
    return {
        ...data,
        character: charObj ? {
            ...charObj,
            user_id: gc?.user_id || null
        } : undefined
    };
}

export async function updateBoardTokenPosition(tokenId: string, x: number, y: number): Promise<void> {
    const { error } = await supabase
        .from('board_tokens')
        .update({ x, y })
        .eq('id', tokenId);

    if (error) throw error;
}

export async function deleteBoardToken(tokenId: string): Promise<void> {
    const { error } = await supabase
        .from('board_tokens')
        .delete()
        .eq('id', tokenId);

    if (error) throw error;
}

