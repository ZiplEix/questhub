-- Migration to add delete_user_account RPC function
CREATE OR REPLACE FUNCTION public.delete_user_account(
    p_delete_monsters BOOLEAN,
    p_delete_games BOOLEAN,
    p_delete_characters BOOLEAN
)
RETURNS VOID SECURITY DEFINER AS $$
DECLARE
    v_user_id UUID;
BEGIN
    v_user_id := auth.uid();
    IF v_user_id IS NULL THEN
        RAISE EXCEPTION 'Non autorisé';
    END IF;

    -- 1. Delete Monsters if requested (type = 'MONSTER' in games mastered by user)
    IF p_delete_monsters THEN
        DELETE FROM public.characters 
        WHERE id IN (
            SELECT gc.character_id 
            FROM public.game_characters gc
            JOIN public.games g ON gc.game_id = g.id
            WHERE g.gm_id = v_user_id
        ) AND type = 'MONSTER';
    END IF;

    -- 2. Delete Characters if requested
    -- Note: This deletes character sheets created in games they master OR assigned to them as a player.
    IF p_delete_characters THEN
        DELETE FROM public.characters
        WHERE id IN (
            SELECT gc.character_id
            FROM public.game_characters gc
            LEFT JOIN public.games g ON gc.game_id = g.id
            WHERE gc.user_id = v_user_id OR g.gm_id = v_user_id
        ) AND type <> 'MONSTER';
    END IF;

    -- 3. Delete Games if requested
    IF p_delete_games THEN
        DELETE FROM public.games WHERE gm_id = v_user_id;
    END IF;

    -- 4. Delete templates created by the user (normally ON DELETE CASCADE, but do it explicitly to be sure)
    DELETE FROM public.templates WHERE created_by = v_user_id;

    -- 5. Delete the auth user
    DELETE FROM auth.users WHERE id = v_user_id;
END;
$$ LANGUAGE plpgsql;
