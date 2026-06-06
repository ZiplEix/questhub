-- Add ping_color to game_players
ALTER TABLE public.game_players ADD COLUMN IF NOT EXISTS ping_color TEXT NOT NULL DEFAULT '#E07A5F';

-- Recreate game_player_users view
DROP VIEW IF EXISTS public.game_player_users;
CREATE OR REPLACE VIEW public.game_player_users AS
SELECT 
    gp.game_id, 
    gp.user_id, 
    COALESCE(u.raw_user_meta_data->>'name', u.email) as name, 
    u.email, 
    gp.joined_at,
    gp.ping_color
FROM public.game_players gp
JOIN auth.users u ON gp.user_id = u.id
WHERE public.is_game_member(gp.game_id, auth.uid());

GRANT SELECT ON public.game_player_users TO authenticated;
