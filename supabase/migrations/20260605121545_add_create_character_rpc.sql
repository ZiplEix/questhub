-- RPC: Create character and link to game atomically (bypasses initial RLS check on unlinked characters)
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
    p_speed INT
)
RETURNS UUID SECURITY DEFINER AS $$
DECLARE
    v_char_id UUID;
BEGIN
    -- 1. Insert into characters
    INSERT INTO public.characters (
        name, race, max_hp, current_hp, is_npc, avatar_url, stats, inventory, money,
        initiative, type, sub_race, armor_class, speed
    )
    VALUES (
        p_name, COALESCE(p_race, ''), p_max_hp, p_max_hp, p_is_npc, NULLIF(p_avatar_url, ''), p_stats, p_inventory, p_money,
        p_initiative, p_type, NULLIF(p_sub_race, ''), p_armor_class, p_speed
    )
    RETURNING id INTO v_char_id;

    -- 2. Insert into game_characters
    INSERT INTO public.game_characters (game_id, character_id, user_id)
    VALUES (p_game_id, v_char_id, p_user_id);

    RETURN v_char_id;
END;
$$ LANGUAGE plpgsql;

GRANT EXECUTE ON FUNCTION public.create_character_rpc TO authenticated, anon;
