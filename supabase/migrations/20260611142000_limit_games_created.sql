-- Limit the number of active games a user can create to 2 in beta mode (exempting admins and moderators)
CREATE OR REPLACE FUNCTION public.check_game_limit()
RETURNS TRIGGER AS $$
DECLARE
    v_game_count INT;
BEGIN
    -- Admins and moderators are exempt from the beta limit
    IF public.is_moderator(NEW.gm_id) THEN
        RETURN NEW;
    END IF;

    -- Count games created by the same GM
    SELECT COUNT(*)::INT INTO v_game_count
    FROM public.games
    WHERE gm_id = NEW.gm_id;

    IF v_game_count >= 2 THEN
        RAISE EXCEPTION 'Vous avez atteint la limite de 2 parties créées en mode bêta.';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Drop trigger if exists
DROP TRIGGER IF EXISTS tr_check_game_limit ON public.games;

-- Create trigger
CREATE TRIGGER tr_check_game_limit
BEFORE INSERT ON public.games
FOR EACH ROW
EXECUTE FUNCTION public.check_game_limit();
