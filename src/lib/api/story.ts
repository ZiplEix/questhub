import { supabase } from '../supabaseClient';

export interface StoryFolder {
    id: string;
    game_id: string;
    name: string;
    is_visible: boolean;
    created_at: string;
}

export interface StoryPage {
    id: string;
    game_id: string;
    folder_id: string | null;
    title: string;
    content: string;
    is_visible: boolean;
    created_at: string;
    updated_at: string;
}

// STORY FOLDERS API
export async function fetchStoryFolders(gameId: string): Promise<StoryFolder[]> {
    const { data, error } = await supabase
        .from('story_folders')
        .select('*')
        .eq('game_id', gameId)
        .order('name', { ascending: true });

    if (error) throw error;
    return data || [];
}

export async function createStoryFolder(gameId: string, name: string): Promise<StoryFolder> {
    const { data, error } = await supabase
        .from('story_folders')
        .insert({
            game_id: gameId,
            name,
            is_visible: false // GM makes it visible manually
        })
        .select()
        .single();

    if (error) throw error;
    return data;
}

export async function updateStoryFolder(
    folderId: string,
    updates: Partial<{ name: string; is_visible: boolean }>
): Promise<StoryFolder> {
    const { data, error } = await supabase
        .from('story_folders')
        .update(updates)
        .eq('id', folderId)
        .select()
        .single();

    if (error) throw error;
    return data;
}

export async function deleteStoryFolder(folderId: string): Promise<void> {
    const { error } = await supabase
        .from('story_folders')
        .delete()
        .eq('id', folderId);

    if (error) throw error;
}

// STORY PAGES API
export async function fetchStoryPages(gameId: string): Promise<StoryPage[]> {
    const { data, error } = await supabase
        .from('story_pages')
        .select('*')
        .eq('game_id', gameId)
        .order('title', { ascending: true });

    if (error) throw error;
    return data || [];
}

export async function createStoryPage(
    gameId: string,
    folderId: string | null,
    title: string,
    content: string = ''
): Promise<StoryPage> {
    const { data, error } = await supabase
        .from('story_pages')
        .insert({
            game_id: gameId,
            folder_id: folderId,
            title,
            content,
            is_visible: false
        })
        .select()
        .single();

    if (error) throw error;
    return data;
}

export async function updateStoryPage(
    pageId: string,
    updates: Partial<{ folder_id: string | null; title: string; content: string; is_visible: boolean }>
): Promise<StoryPage> {
    const { data, error } = await supabase
        .from('story_pages')
        .update({
            ...updates,
            updated_at: new Date().toISOString()
        })
        .eq('id', pageId)
        .select()
        .single();

    if (error) throw error;
    return data;
}

export async function deleteStoryPage(pageId: string): Promise<void> {
    const { error } = await supabase
        .from('story_pages')
        .delete()
        .eq('id', pageId);

    if (error) throw error;
}
