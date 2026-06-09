import { supabase } from '../supabaseClient';
import { createCharacter } from './character';

export interface MarketplaceTemplate {
    id: string;
    created_by: string;
    author_name: string;
    name: string;
    description: string | null;
    type: 'PERSONNAGE' | 'PNJ' | 'MONSTRE' | 'BUNDLE';
    data: any; // Contains the character stats or { items: [...] } for bundle
    is_public: boolean;
    uses: number;
    created_at: string;
}

export async function fetchMarketplaceTemplates(filters?: {
    type?: string;
    search?: string;
    onlyUser?: boolean;
}): Promise<MarketplaceTemplate[]> {
    const { data: { user } } = await supabase.auth.getUser();
    let query = supabase.from('marketplace_templates').select('*');
    
    // Filter by type if specified
    if (filters?.type && filters.type !== 'TOUT') {
        if (filters.type === 'BUNDLE_ONLY') {
            query = query.eq('type', 'BUNDLE');
        } else {
            query = query.eq('type', filters.type);
        }
    }
    
    // Filter by ownership if specified
    if (filters?.onlyUser && user) {
        query = query.eq('created_by', user.id);
    }
    
    const { data, error } = await query
        .order('uses', { ascending: false })
        .order('created_at', { ascending: false });
    if (error) throw error;
    
    let list = data || [];
    
    // Search query matching name, description, or author
    if (filters?.search) {
        const searchLower = filters.search.toLowerCase();
        list = list.filter(t => 
            t.name.toLowerCase().includes(searchLower) ||
            (t.description && t.description.toLowerCase().includes(searchLower)) ||
            (t.author_name && t.author_name.toLowerCase().includes(searchLower))
        );
    }
    
    return list as MarketplaceTemplate[];
}

export async function createTemplate(payload: {
    name: string;
    description: string | null;
    type: 'PERSONNAGE' | 'PNJ' | 'MONSTRE' | 'BUNDLE';
    data: any;
    is_public?: boolean;
}): Promise<MarketplaceTemplate> {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Unauthorized');
    
    const { data, error } = await supabase
        .from('templates')
        .insert({
            created_by: user.id,
            name: payload.name,
            description: payload.description,
            type: payload.type,
            data: payload.data,
            is_public: payload.is_public !== undefined ? payload.is_public : true
        })
        .select()
        .single();
        
    if (error) throw error;
    return data;
}

export async function updateTemplate(
    templateId: string,
    payload: {
        name?: string;
        description?: string | null;
        is_public?: boolean;
        data?: any;
    }
): Promise<void> {
    const { error } = await supabase
        .from('templates')
        .update(payload)
        .eq('id', templateId);
        
    if (error) throw error;
}

export async function deleteTemplate(templateId: string): Promise<void> {
    const { error } = await supabase
        .from('templates')
        .delete()
        .eq('id', templateId);
        
    if (error) throw error;
}

async function cloneCharacterToGame(charData: any, type: 'PERSONNAGE' | 'PNJ' | 'MONSTRE', gameId: string) {
    // Map from template type to game character database types
    let charType = 'PLAYER';
    let isNpc = false;
    
    if (type === 'PNJ') {
        charType = 'NPC';
        isNpc = true;
    } else if (type === 'MONSTRE') {
        charType = 'MONSTER';
        isNpc = true;
    }
    
    const payload = {
        name: charData.name,
        race: charData.race || '',
        sub_race: charData.sub_race || null,
        max_hp: charData.max_hp || 10,
        is_npc: isNpc,
        avatar_url: charData.avatar_url || null,
        stats: charData.stats || {},
        inventory: charData.inventory || [],
        money: charData.money || 0,
        initiative: charData.initiative || 0,
        type: charType,
        armor_class: charData.armor_class || 10,
        speed: charData.speed || 30,
        user_id: null, // initially unassigned
        description: charData.description || ''
    };
    
    return await createCharacter(gameId, payload);
}

export async function importTemplateToGame(
    template: any,
    gameId: string
): Promise<void> {
    if (template.type === 'BUNDLE' && !template.is_virtual) {
        // Bundle imports all bundled items recursively
        const items = template.data?.items || [];
        for (const item of items) {
            // item shape: { name, type, data }
            await cloneCharacterToGame(item.data, item.type, gameId);
        }
    } else {
        // Single character/monster template import
        await cloneCharacterToGame(template.data, template.type as any, gameId);
    }
    
    // Increment the uses counter atomically via RPC
    const targetId = template.is_virtual ? template.parent_bundle_id : template.id;
    const { error } = await supabase.rpc('increment_template_uses', { p_template_id: targetId });
    if (error) console.error('Failed to increment template uses:', error);
}
