export interface Character {
    id: string;
    game_id: string;
    user_id?: string;
    name: string;
    race: string;
    max_hp: number;
    current_hp: number;
    avatar_url?: string;
    stats: Record<string, string>;
    inventory: InventoryItem[];
    is_npc: boolean;
    money: number;
    created_at: string;
    player_name?: string;

    // New fields
    initiative: number;
    age: string;
    height: string;
    weight: string;
    max_spells: number;
    spells: Record<string, { name: string; description: string }[]>;
    abilities: string;
    experience: number;
}

export interface InventoryItem {
    name: string;
    quantity: string;
    image_url?: string;
    icon_name?: string;
}
