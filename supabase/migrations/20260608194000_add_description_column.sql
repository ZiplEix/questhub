-- Migration: Add description column to characters and update create_character_rpc

-- 1. Add description column to characters table
ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS description TEXT DEFAULT '';

-- 2. Drop the previous signature of the function
DROP FUNCTION IF EXISTS public.create_character_rpc(
    p_game_id UUID,
    p_user_id UUID,
    p_name TEXT,
    p_race TEXT,
    p_max_hp INT,
    p_is_npc BOOLEAN,
    p_avatar_url TEXT,
    p_stats JSONB,
    p_inventory JSONB,
    p_money INT,
    p_initiative INT,
    p_type TEXT,
    p_sub_race TEXT,
    p_armor_class INT,
    p_speed INT,
    p_strength INT,
    p_strength_mod INT,
    p_dexterity INT,
    p_dexterity_mod INT,
    p_constitution INT,
    p_constitution_mod INT,
    p_intelligence INT,
    p_intelligence_mod INT,
    p_wisdom INT,
    p_wisdom_mod INT,
    p_charisma INT,
    p_charisma_mod INT,
    p_age TEXT,
    p_height TEXT,
    p_weight TEXT,
    p_max_spells INT,
    p_spells JSONB,
    p_abilities TEXT,
    p_experience INT
);

-- 3. Recreate the function with description field
CREATE OR REPLACE FUNCTION public.create_character_rpc(
    p_game_id UUID,
    p_user_id UUID,
    p_name TEXT,
    p_race TEXT,
    p_max_hp INT,
    p_is_npc BOOLEAN,
    p_avatar_url TEXT,
    p_stats JSONB,
    p_inventory JSONB,
    p_money INT,
    p_initiative INT,
    p_type TEXT,
    p_sub_race TEXT,
    p_armor_class INT,
    p_speed INT,
    p_strength INT DEFAULT 10,
    p_strength_mod INT DEFAULT 0,
    p_dexterity INT DEFAULT 10,
    p_dexterity_mod INT DEFAULT 0,
    p_constitution INT DEFAULT 10,
    p_constitution_mod INT DEFAULT 0,
    p_intelligence INT DEFAULT 10,
    p_intelligence_mod INT DEFAULT 0,
    p_wisdom INT DEFAULT 10,
    p_wisdom_mod INT DEFAULT 0,
    p_charisma INT DEFAULT 10,
    p_charisma_mod INT DEFAULT 0,
    p_age TEXT DEFAULT '',
    p_height TEXT DEFAULT '',
    p_weight TEXT DEFAULT '',
    p_max_spells INT DEFAULT 0,
    p_spells JSONB DEFAULT '{}'::jsonb,
    p_abilities TEXT DEFAULT '',
    p_experience INT DEFAULT 0,
    p_description TEXT DEFAULT ''
)
RETURNS UUID SECURITY DEFINER AS $$
DECLARE
    v_char_id UUID;
BEGIN
    INSERT INTO public.characters (
        name, race, max_hp, current_hp, is_npc, avatar_url, stats, inventory, money,
        initiative, type, sub_race, armor_class, speed,
        strength, strength_mod, dexterity, dexterity_mod,
        constitution, constitution_mod, intelligence, intelligence_mod,
        wisdom, wisdom_mod, charisma, charisma_mod,
        age, height, weight, max_spells, spells, abilities, experience,
        description
    )
    VALUES (
        p_name, COALESCE(p_race, ''), p_max_hp, p_max_hp, p_is_npc, NULLIF(p_avatar_url, ''), p_stats, p_inventory, p_money,
        p_initiative, p_type, NULLIF(p_sub_race, ''), p_armor_class, p_speed,
        p_strength, p_strength_mod, p_dexterity, p_dexterity_mod,
        p_constitution, p_constitution_mod, p_intelligence, p_intelligence_mod,
        p_wisdom, p_wisdom_mod, p_charisma, p_charisma_mod,
        COALESCE(p_age, ''), COALESCE(p_height, ''), COALESCE(p_weight, ''), p_max_spells, p_spells, COALESCE(p_abilities, ''), p_experience,
        COALESCE(p_description, '')
    )
    RETURNING id INTO v_char_id;

    INSERT INTO public.game_characters (game_id, character_id, user_id)
    VALUES (p_game_id, v_char_id, p_user_id);

    RETURN v_char_id;
END;
$$ LANGUAGE plpgsql;

GRANT EXECUTE ON FUNCTION public.create_character_rpc TO authenticated, anon;
