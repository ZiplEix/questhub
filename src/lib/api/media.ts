import { supabase } from '../supabaseClient';
import { validateImage, checkStorageLimit } from './storage';

export interface MediaAsset {
    id: string;
    created_at: string;
    user_id: string;
    name: string;
    url: string;
    size: number;
    mime_type: string;
}

export async function fetchMediaLibrary(): Promise<MediaAsset[]> {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Unauthorized');

    const { data, error } = await supabase
        .from('media_library')
        .select('*')
        .order('created_at', { ascending: false });

    if (error) throw error;
    return data || [];
}

export async function uploadToMediaLibrary(file: File): Promise<MediaAsset> {
    const validation = validateImage(file, 10);
    if (!validation.valid) {
        throw new Error(validation.error);
    }

    await checkStorageLimit(file.size);

    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Unauthorized');

    const cleanedName = file.name
        .replace(/\s+/g, '_')
        .replace(/[^a-zA-Z0-9._-]/g, '');
    const filePath = `media/${user.id}/${Date.now()}_${cleanedName}`;

    // 1. Upload to storage
    const { error: uploadError } = await supabase.storage
        .from('images')
        .upload(filePath, file);

    if (uploadError) throw uploadError;

    // Get public URL
    const { data: { publicUrl } } = supabase.storage
        .from('images')
        .getPublicUrl(filePath);

    // 2. Insert record in database
    const { data, error: dbError } = await supabase
        .from('media_library')
        .insert({
            name: file.name,
            url: publicUrl,
            size: file.size,
            mime_type: file.type
        })
        .select()
        .single();

    if (dbError) {
        // Cleanup storage file on failure
        await supabase.storage.from('images').remove([filePath]);
        throw dbError;
    }

    return data;
}

export async function deleteFromMediaLibrary(id: string, url: string): Promise<void> {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Unauthorized');

    // Extract storage path from public URL
    const marker = '/storage/v1/object/public/images/';
    const index = url.indexOf(marker);
    if (index === -1) throw new Error('Invalid media URL');
    const filePath = url.substring(index + marker.length);

    // 1. Delete record from database
    const { error: dbError } = await supabase
        .from('media_library')
        .delete()
        .eq('id', id);

    if (dbError) throw dbError;

    // 2. Delete file from storage
    const { error: storageError } = await supabase.storage
        .from('images')
        .remove([filePath]);

    if (storageError) {
        console.error("Failed to remove storage file", storageError);
    }
}
