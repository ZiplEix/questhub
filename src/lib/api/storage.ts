import { supabase } from '../supabaseClient';

export async function uploadImage(file: File): Promise<string> {
    const fileExt = file.name.split('.').pop();
    const fileName = `${Math.random().toString(36).substring(2, 15)}.${fileExt}`;
    const filePath = `uploads/${fileName}`;

    const { error } = await supabase.storage
        .from('images')
        .upload(filePath, file);

    if (error) throw error;

    // Get public URL
    const { data: { publicUrl } } = supabase.storage
        .from('images')
        .getPublicUrl(filePath);

    return publicUrl;
}

export async function uploadStoryAsset(gameId: string, file: File): Promise<string> {
    // Clean filename to be URL-safe (keep alphanumeric, dots, dashes, underscores)
    const cleanedName = file.name
        .replace(/\s+/g, '_')
        .replace(/[^a-zA-Z0-9._-]/g, '');
    
    const filePath = `story/${gameId}/${Date.now()}_${cleanedName}`;

    const { error } = await supabase.storage
        .from('images')
        .upload(filePath, file);

    if (error) throw error;

    const { data: { publicUrl } } = supabase.storage
        .from('images')
        .getPublicUrl(filePath);

    return publicUrl;
}

export async function listStoryAssets(gameId: string): Promise<Array<{ name: string; url: string; created_at: string; size: number }>> {
    const folderPath = `story/${gameId}`;
    const { data, error } = await supabase.storage
        .from('images')
        .list(folderPath, {
            sortBy: { column: 'name', order: 'desc' }
        });

    if (error) throw error;
    if (!data) return [];

    return data.map(item => {
        const filePath = `${folderPath}/${item.name}`;
        const { data: { publicUrl } } = supabase.storage
            .from('images')
            .getPublicUrl(filePath);

        return {
            name: item.name,
            url: publicUrl,
            created_at: item.created_at || '',
            size: item.metadata?.size || 0
        };
    });
}

export async function deleteStoryAsset(gameId: string, fileName: string): Promise<void> {
    const filePath = `story/${gameId}/${fileName}`;
    const { error } = await supabase.storage
        .from('images')
        .remove([filePath]);

    if (error) throw error;
}
