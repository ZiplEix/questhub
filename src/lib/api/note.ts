import { supabase } from '../supabaseClient';

export interface Note {
    content: string;
}

export async function fetchNotes(gameId: string, characterId: string): Promise<Note> {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Unauthorized');

    const { data: gc } = await supabase
        .from('game_characters')
        .select('user_id')
        .eq('game_id', gameId)
        .eq('character_id', characterId)
        .single();

    const userId = gc?.user_id || user.id;

    const { data, error } = await supabase
        .from('notes')
        .select('content')
        .eq('game_id', gameId)
        .eq('user_id', userId)
        .maybeSingle();

    if (error) throw error;
    return { content: data?.content || '' };
}

export async function updateNotes(gameId: string, characterId: string, content: string): Promise<void> {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Unauthorized');

    const { data: gc } = await supabase
        .from('game_characters')
        .select('user_id')
        .eq('game_id', gameId)
        .eq('character_id', characterId)
        .single();

    const userId = gc?.user_id || user.id;

    const { error } = await supabase
        .from('notes')
        .upsert({
            game_id: gameId,
            user_id: userId,
            content
        });

    if (error) throw error;
}
