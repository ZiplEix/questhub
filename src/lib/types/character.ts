export interface Character {
    id: string;
    game_id: string;
    user_id?: string | null;
    name: string;
    race: string;
    max_hp: number;
    current_hp: number;
    avatar_url?: string | null;
    stats: Record<string, { value: number; modifier: number }>;
    inventory: InventoryItem[];
    is_npc: boolean;
    money: number;
    created_at: string;
    player_name?: string;
    type: string;
    sub_race?: string | null;

    // New fields
    initiative: number;
    age?: string | null;
    height?: string | null;
    weight?: string | null;
    max_spells?: number;
    spells?: Record<string, { name: string; description: string; charges: string }[]>;
    abilities?: string | null;
    experience?: number;
    armor_class: number;
    speed: number;
    conditions: string[];
}

export interface InventoryItem {
    name: string;
    quantity: string;
    description?: string;
    image_url?: string;
    icon_name?: string;
}
