-- Recreate views without security_invoker = true (since auth.users is protected)
-- Instead, we filter by auth.uid() directly in the view definitions and grant SELECT to authenticated role.

-- 1. user_games
DROP VIEW IF EXISTS public.user_games;
CREATE OR REPLACE VIEW public.user_games AS
SELECT DISTINCT 
    g.id, 
    g.name, 
    g.gm_id, 
    COALESCE(u.raw_user_meta_data->>'name', u.email) as gm_name, 
    g.invite_code, 
    g.is_active, 
    g.image_url, 
    g.state, 
    g.created_at
FROM public.games g
JOIN auth.users u ON g.gm_id = u.id
LEFT JOIN public.game_players gp ON g.id = gp.game_id
WHERE g.gm_id = auth.uid() OR gp.user_id = auth.uid();

GRANT SELECT ON public.user_games TO authenticated;


-- 2. game_player_users
DROP VIEW IF EXISTS public.game_player_users;
CREATE OR REPLACE VIEW public.game_player_users AS
SELECT 
    gp.game_id, 
    gp.user_id, 
    COALESCE(u.raw_user_meta_data->>'name', u.email) as name, 
    u.email, 
    gp.joined_at
FROM public.game_players gp
JOIN auth.users u ON gp.user_id = u.id
WHERE public.is_game_member(gp.game_id, auth.uid());

GRANT SELECT ON public.game_player_users TO authenticated;


-- 3. game_invitation_users
DROP VIEW IF EXISTS public.game_invitation_users;
CREATE OR REPLACE VIEW public.game_invitation_users AS
SELECT 
    gi.id, 
    gi.game_id, 
    gi.user_id, 
    COALESCE(u.raw_user_meta_data->>'name', u.email) as user_name, 
    gi.created_at
FROM public.game_invitations gi
JOIN auth.users u ON gi.user_id = u.id
WHERE public.is_game_gm(gi.game_id, auth.uid()) OR gi.user_id = auth.uid();

GRANT SELECT ON public.game_invitation_users TO authenticated;
