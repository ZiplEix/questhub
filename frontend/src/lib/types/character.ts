export interface Character {
    id: string;
    game_id: string;
    user_id?: string;
    name: string;
    race: string;
    max_hp: number;
    current_hp: number;
    avatar_url?: string;
    stats: Record<string, { value: number; modifier: number }>;
    inventory: InventoryItem[];
    is_npc: boolean;
    money: number;
    created_at: string;
    player_name?: string;
    type: "PLAYER" | "NPC" | "MONSTER";
    sub_race?: string;

    // New fields
    initiative: number;
    age: string;
    height: string;
    weight: string;
    max_spells: number;
    spells: Record<string, { name: string; description: string; charges: string }[]>;
    abilities: string;
    experience: number;
    armor_class: number;
    speed: number;
}

export interface InventoryItem {
    name: string;
    quantity: string;
    description?: string;
    image_url?: string;
    icon_name?: string;
}
