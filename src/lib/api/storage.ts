import { supabase } from '../supabaseClient';

export const ALLOWED_IMAGE_TYPES = ['image/png', 'image/jpeg', 'image/jpg', 'image/webp'];

export function validateImage(file: File, maxMb: number = 2): { valid: boolean; error?: string } {
    if (!ALLOWED_IMAGE_TYPES.includes(file.type)) {
        if (file.type === 'image/gif') {
            return { valid: false, error: "Les fichiers GIF ne sont pas autorisés." };
        }
        return { valid: false, error: "Format d'image non supporté. Veuillez utiliser PNG, JPG, JPEG ou WEBP." };
    }

    const maxSize = maxMb * 1024 * 1024;
    if (file.size > maxSize) {
        return { valid: false, error: `La taille de l'image ne doit pas dépasser ${maxMb} Mo.` };
    }

    return { valid: true };
}

export async function getUserStorageUsage(): Promise<number> {
    const { data, error } = await supabase.rpc('get_user_storage_usage');
    if (error) {
        console.error('Error fetching storage usage:', error);
        return 0;
    }
    return data || 0;
}

export const STORAGE_LIMIT = 50 * 1024 * 1024; // 50 MB in bytes

export async function checkStorageLimit(newFileSize: number): Promise<void> {
    const currentUsage = await getUserStorageUsage();
    if (currentUsage + newFileSize > STORAGE_LIMIT) {
        throw new Error("Limite de stockage de 50 Mo dépassée. Veuillez utiliser des liens ou supprimer d'autres médias.");
    }
}

export async function uploadImage(file: File, maxMb: number = 2): Promise<string> {
    const validation = validateImage(file, maxMb);
    if (!validation.valid) {
        throw new Error(validation.error);
    }

    await checkStorageLimit(file.size);

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
    const validation = validateImage(file, 2); // 2MB limit for story assets
    if (!validation.valid) {
        throw new Error(validation.error);
    }

    await checkStorageLimit(file.size);

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
