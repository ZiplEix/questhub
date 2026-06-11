-- Create user_roles table
CREATE TABLE IF NOT EXISTS public.user_roles (
    user_id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
    role VARCHAR(50) CHECK (role IN ('admin', 'moderator')),
    is_banned BOOLEAN DEFAULT FALSE NOT NULL
);

-- Enable RLS on user_roles
ALTER TABLE public.user_roles ENABLE ROW LEVEL SECURITY;

-- Allow select access to authenticated users to view roles (e.g. for badges)
CREATE POLICY "Allow select access to authenticated"
ON public.user_roles FOR SELECT
TO authenticated
USING (true);

-- Functions to check user permissions (SECURITY DEFINER to bypass RLS dynamically)
CREATE OR REPLACE FUNCTION public.is_admin(p_user_id UUID)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM public.user_roles
        WHERE user_id = p_user_id AND role = 'admin' AND is_banned = FALSE
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE OR REPLACE FUNCTION public.is_moderator(p_user_id UUID)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM public.user_roles
        WHERE user_id = p_user_id AND role IN ('admin', 'moderator') AND is_banned = FALSE
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- RPC: Get overall platform statistics
CREATE OR REPLACE FUNCTION public.admin_get_stats()
RETURNS JSONB AS $$
DECLARE
    v_total_users INT;
    v_total_games INT;
    v_total_templates INT;
    v_pending_tickets INT;
    v_banned_users INT;
BEGIN
    -- Auth check
    IF NOT public.is_moderator(auth.uid()) THEN
        RAISE EXCEPTION 'Unauthorized';
    END IF;

    SELECT COUNT(*)::INT INTO v_total_users FROM auth.users;
    SELECT COUNT(*)::INT INTO v_total_games FROM public.games;
    SELECT COUNT(*)::INT INTO v_total_templates FROM public.templates;
    SELECT COUNT(*)::INT INTO v_pending_tickets FROM public.support_tickets WHERE status = 'PENDING';
    SELECT COUNT(*)::INT INTO v_banned_users FROM public.user_roles WHERE is_banned = TRUE;

    RETURN jsonb_build_object(
        'total_users', v_total_users,
        'total_games', v_total_games,
        'total_templates', v_total_templates,
        'pending_tickets', v_pending_tickets,
        'banned_users', v_banned_users
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- RPC: Get list of all users
CREATE OR REPLACE FUNCTION public.admin_get_users()
RETURNS TABLE (
    id UUID,
    email VARCHAR(255),
    name TEXT,
    created_at TIMESTAMPTZ,
    role VARCHAR(50),
    is_banned BOOLEAN
) AS $$
BEGIN
    -- Auth check
    IF NOT public.is_moderator(auth.uid()) THEN
        RAISE EXCEPTION 'Unauthorized';
    END IF;

    RETURN QUERY
    SELECT 
        u.id,
        u.email::VARCHAR(255),
        COALESCE(u.raw_user_meta_data->>'name', u.email) as name,
        u.created_at,
        ur.role::VARCHAR(50),
        COALESCE(ur.is_banned, FALSE) as is_banned
    FROM auth.users u
    LEFT JOIN public.user_roles ur ON u.id = ur.user_id
    ORDER BY u.created_at DESC;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- RPC: Set a user's role
CREATE OR REPLACE FUNCTION public.admin_set_user_role(p_target_user_id UUID, p_role VARCHAR)
RETURNS VOID AS $$
BEGIN
    -- Auth check
    IF NOT public.is_admin(auth.uid()) THEN
        RAISE EXCEPTION 'Unauthorized';
    END IF;

    -- Target cannot be the caller
    IF auth.uid() = p_target_user_id THEN
        RAISE EXCEPTION 'Vous ne pouvez pas modifier votre propre rôle.';
    END IF;

    -- Validate role
    IF p_role IS NOT NULL AND p_role NOT IN ('admin', 'moderator') THEN
        RAISE EXCEPTION 'Rôle invalide';
    END IF;

    INSERT INTO public.user_roles (user_id, role)
    VALUES (p_target_user_id, p_role)
    ON CONFLICT (user_id) DO UPDATE
    SET role = EXCLUDED.role;
    
    -- Clean up row if empty
    DELETE FROM public.user_roles
    WHERE user_id = p_target_user_id AND role IS NULL AND is_banned = FALSE;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- RPC: Set a user's banned status
CREATE OR REPLACE FUNCTION public.admin_set_user_banned(p_target_user_id UUID, p_is_banned BOOLEAN)
RETURNS VOID AS $$
BEGIN
    -- Auth check
    IF NOT public.is_admin(auth.uid()) THEN
        RAISE EXCEPTION 'Unauthorized';
    END IF;

    -- Target cannot be the caller
    IF auth.uid() = p_target_user_id THEN
        RAISE EXCEPTION 'Vous ne pouvez pas vous bannir vous-même.';
    END IF;

    INSERT INTO public.user_roles (user_id, is_banned)
    VALUES (p_target_user_id, p_is_banned)
    ON CONFLICT (user_id) DO UPDATE
    SET is_banned = EXCLUDED.is_banned;
    
    -- Clean up row if empty
    DELETE FROM public.user_roles
    WHERE user_id = p_target_user_id AND role IS NULL AND is_banned = FALSE;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- RPC: Get all games/tables
CREATE OR REPLACE FUNCTION public.admin_get_games()
RETURNS TABLE (
    id UUID,
    name TEXT,
    gm_name TEXT,
    gm_email VARCHAR(255),
    state VARCHAR(50),
    is_active BOOLEAN,
    created_at TIMESTAMPTZ,
    player_count INT
) AS $$
BEGIN
    -- Auth check
    IF NOT public.is_moderator(auth.uid()) THEN
        RAISE EXCEPTION 'Unauthorized';
    END IF;

    RETURN QUERY
    SELECT 
        g.id,
        g.name,
        COALESCE(u.raw_user_meta_data->>'name', u.email) as gm_name,
        u.email::VARCHAR(255) as gm_email,
        g.state::VARCHAR(50),
        g.is_active,
        g.created_at,
        (SELECT COUNT(*)::INT FROM public.game_players gp WHERE gp.game_id = g.id) as player_count
    FROM public.games g
    JOIN auth.users u ON g.gm_id = u.id
    ORDER BY g.created_at DESC;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- RPC: Archive/Restore a game
CREATE OR REPLACE FUNCTION public.admin_archive_game(p_game_id UUID, p_is_active BOOLEAN)
RETURNS VOID AS $$
BEGIN
    -- Auth check
    IF NOT public.is_moderator(auth.uid()) THEN
        RAISE EXCEPTION 'Unauthorized';
    END IF;

    UPDATE public.games
    SET is_active = p_is_active
    WHERE id = p_game_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- RPC: Delete a game
CREATE OR REPLACE FUNCTION public.admin_delete_game(p_game_id UUID)
RETURNS VOID AS $$
BEGIN
    -- Auth check
    IF NOT public.is_moderator(auth.uid()) THEN
        RAISE EXCEPTION 'Unauthorized';
    END IF;

    DELETE FROM public.games
    WHERE id = p_game_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- RPC: Get all templates (marketplace)
CREATE OR REPLACE FUNCTION public.admin_get_templates()
RETURNS TABLE (
    id UUID,
    name TEXT,
    type VARCHAR(50),
    creator_name TEXT,
    is_public BOOLEAN,
    uses INT,
    created_at TIMESTAMPTZ
) AS $$
BEGIN
    -- Auth check
    IF NOT public.is_moderator(auth.uid()) THEN
        RAISE EXCEPTION 'Unauthorized';
    END IF;

    RETURN QUERY
    SELECT 
        t.id,
        t.name,
        t.type::VARCHAR(50),
        COALESCE(u.raw_user_meta_data->>'name', u.email) as creator_name,
        t.is_public,
        t.uses,
        t.created_at
    FROM public.templates t
    JOIN auth.users u ON t.created_by = u.id
    ORDER BY t.created_at DESC;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- RPC: Toggle template public status
CREATE OR REPLACE FUNCTION public.admin_toggle_template_public(p_template_id UUID, p_is_public BOOLEAN)
RETURNS VOID AS $$
BEGIN
    -- Auth check
    IF NOT public.is_moderator(auth.uid()) THEN
        RAISE EXCEPTION 'Unauthorized';
    END IF;

    UPDATE public.templates
    SET is_public = p_is_public
    WHERE id = p_template_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- RPC: Delete a template
CREATE OR REPLACE FUNCTION public.admin_delete_template(p_template_id UUID)
RETURNS VOID AS $$
BEGIN
    -- Auth check
    IF NOT public.is_moderator(auth.uid()) THEN
        RAISE EXCEPTION 'Unauthorized';
    END IF;

    DELETE FROM public.templates
    WHERE id = p_template_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- RPC: Get all support tickets
CREATE OR REPLACE FUNCTION public.admin_get_tickets()
RETURNS TABLE (
    id UUID,
    created_at TIMESTAMPTZ,
    email TEXT,
    category TEXT,
    subject TEXT,
    message TEXT,
    attachment_url TEXT,
    status TEXT,
    user_name TEXT
) AS $$
BEGIN
    -- Auth check
    IF NOT public.is_moderator(auth.uid()) THEN
        RAISE EXCEPTION 'Unauthorized';
    END IF;

    RETURN QUERY
    SELECT 
        st.id,
        st.created_at,
        st.email,
        st.category,
        st.subject,
        st.message,
        st.attachment_url,
        st.status,
        COALESCE(u.raw_user_meta_data->>'name', u.email) as user_name
    FROM public.support_tickets st
    LEFT JOIN auth.users u ON st.user_id = u.id
    ORDER BY 
        CASE st.status 
            WHEN 'PENDING' THEN 1
            WHEN 'IN_PROGRESS' THEN 2
            WHEN 'RESOLVED' THEN 3
            WHEN 'CLOSED' THEN 4
            ELSE 5
        END,
        st.created_at DESC;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- RPC: Update support ticket status
CREATE OR REPLACE FUNCTION public.admin_update_ticket_status(p_ticket_id UUID, p_status TEXT)
RETURNS VOID AS $$
BEGIN
    -- Auth check
    IF NOT public.is_moderator(auth.uid()) THEN
        RAISE EXCEPTION 'Unauthorized';
    END IF;

    -- Validate status
    IF p_status NOT IN ('PENDING', 'IN_PROGRESS', 'RESOLVED', 'CLOSED') THEN
        RAISE EXCEPTION 'Status invalide';
    END IF;

    UPDATE public.support_tickets
    SET status = p_status
    WHERE id = p_ticket_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;
