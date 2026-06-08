-- Migration: Extract the 6 classic D&D ability scores from JSONB stats into dedicated columns

-- 1. Add dedicated columns for each ability score (value + modifier)
ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS strength INT NOT NULL DEFAULT 10;
ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS strength_mod INT NOT NULL DEFAULT 0;
ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS dexterity INT NOT NULL DEFAULT 10;
ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS dexterity_mod INT NOT NULL DEFAULT 0;
ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS constitution INT NOT NULL DEFAULT 10;
ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS constitution_mod INT NOT NULL DEFAULT 0;
ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS intelligence INT NOT NULL DEFAULT 10;
ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS intelligence_mod INT NOT NULL DEFAULT 0;
ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS wisdom INT NOT NULL DEFAULT 10;
ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS wisdom_mod INT NOT NULL DEFAULT 0;
ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS charisma INT NOT NULL DEFAULT 10;
ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS charisma_mod INT NOT NULL DEFAULT 0;

-- 2. Migrate existing data from JSONB stats into the new columns
UPDATE public.characters SET
    strength       = COALESCE((stats->'Force'->>'value')::int, 10),
    strength_mod   = COALESCE((stats->'Force'->>'modifier')::int, 0),
    dexterity      = COALESCE((stats->'Dextérité'->>'value')::int, 10),
    dexterity_mod  = COALESCE((stats->'Dextérité'->>'modifier')::int, 0),
    constitution   = COALESCE((stats->'Constitution'->>'value')::int, 10),
    constitution_mod = COALESCE((stats->'Constitution'->>'modifier')::int, 0),
    intelligence   = COALESCE((stats->'Intelligence'->>'value')::int, 10),
    intelligence_mod = COALESCE((stats->'Intelligence'->>'modifier')::int, 0),
    wisdom         = COALESCE((stats->'Sagesse'->>'value')::int, 10),
    wisdom_mod     = COALESCE((stats->'Sagesse'->>'modifier')::int, 0),
    charisma       = COALESCE((stats->'Charisme'->>'value')::int, 10),
    charisma_mod   = COALESCE((stats->'Charisme'->>'modifier')::int, 0)
WHERE stats IS NOT NULL AND stats != '{}'::jsonb;

-- 3. Remove the 6 core stats from the JSONB field (keep only custom stats)
UPDATE public.characters SET
    stats = stats - 'Force' - 'Dextérité' - 'Constitution' - 'Intelligence' - 'Sagesse' - 'Charisme'
WHERE stats IS NOT NULL AND stats != '{}'::jsonb;

-- 4. Update the create_character_rpc to include the new columns
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
    p_speed INT
);

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
    p_charisma_mod INT DEFAULT 0
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
        wisdom, wisdom_mod, charisma, charisma_mod
    )
    VALUES (
        p_name, COALESCE(p_race, ''), p_max_hp, p_max_hp, p_is_npc, NULLIF(p_avatar_url, ''), p_stats, p_inventory, p_money,
        p_initiative, p_type, NULLIF(p_sub_race, ''), p_armor_class, p_speed,
        p_strength, p_strength_mod, p_dexterity, p_dexterity_mod,
        p_constitution, p_constitution_mod, p_intelligence, p_intelligence_mod,
        p_wisdom, p_wisdom_mod, p_charisma, p_charisma_mod
    )
    RETURNING id INTO v_char_id;

    INSERT INTO public.game_characters (game_id, character_id, user_id)
    VALUES (p_game_id, v_char_id, p_user_id);

    RETURN v_char_id;
END;
$$ LANGUAGE plpgsql;

GRANT EXECUTE ON FUNCTION public.create_character_rpc TO authenticated, anon;
