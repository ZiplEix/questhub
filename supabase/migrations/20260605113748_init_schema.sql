-- ----------------------------------------------------
-- TABLES
-- ----------------------------------------------------

CREATE TABLE IF NOT EXISTS public.games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    gm_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    invite_code TEXT NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    image_url TEXT,
    state VARCHAR(50) DEFAULT 'ongoing',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_games_gm_id ON public.games(gm_id);

CREATE TABLE IF NOT EXISTS public.templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_by UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    data JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_public BOOLEAN DEFAULT FALSE,
    uses INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_templates_type ON public.templates(type);
CREATE INDEX IF NOT EXISTS idx_templates_public ON public.templates(is_public);

CREATE TABLE IF NOT EXISTS public.game_players (
    game_id UUID NOT NULL REFERENCES public.games(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (game_id, user_id)
);

CREATE TABLE IF NOT EXISTS public.characters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    avatar_url TEXT,
    stats JSONB NOT NULL DEFAULT '{}'::jsonb, 
    inventory JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_npc BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    race TEXT NOT NULL DEFAULT '',
    max_hp INT NOT NULL DEFAULT 10,
    current_hp INT NOT NULL DEFAULT 10,
    money INTEGER DEFAULT 0,
    initiative INTEGER DEFAULT 0,
    age TEXT DEFAULT '',
    height TEXT DEFAULT '',
    weight TEXT DEFAULT '',
    max_spells INTEGER DEFAULT 0,
    spells JSONB DEFAULT '{}'::jsonb,
    abilities TEXT DEFAULT '',
    experience INTEGER DEFAULT 0,
    type VARCHAR(50) DEFAULT 'PLAYER',
    sub_race VARCHAR(255),
    armor_class INTEGER DEFAULT 10,
    speed INTEGER DEFAULT 30
);

CREATE TABLE IF NOT EXISTS public.game_characters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES public.games(id) ON DELETE CASCADE,
    character_id UUID NOT NULL REFERENCES public.characters(id) ON DELETE CASCADE,
    user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(game_id, character_id)
);

CREATE INDEX IF NOT EXISTS idx_game_characters_game_id ON public.game_characters(game_id);
CREATE INDEX IF NOT EXISTS idx_game_characters_user_id ON public.game_characters(user_id);

CREATE TABLE IF NOT EXISTS public.game_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES public.games(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(game_id, user_id)
);

CREATE TABLE IF NOT EXISTS public.notes (
    game_id UUID NOT NULL REFERENCES public.games(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    content TEXT DEFAULT '',
    PRIMARY KEY (game_id, user_id)
);

CREATE TABLE IF NOT EXISTS public.messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES public.games(id) ON DELETE CASCADE,
    sender_id TEXT NOT NULL,
    sender_name TEXT NOT NULL,
    content TEXT NOT NULL,
    type TEXT NOT NULL, -- 'CHAT_GLOBAL', 'CHAT_PRIVATE', 'EVENT'
    target_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_messages_game_id ON public.messages(game_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON public.messages(created_at);

-- ----------------------------------------------------
-- HELPER FUNCTIONS & RPC PROCEDURES
-- ----------------------------------------------------

CREATE OR REPLACE FUNCTION public.is_game_member(p_game_id UUID, p_user_id UUID)
RETURNS BOOLEAN SECURITY DEFINER AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM public.games WHERE id = p_game_id AND gm_id = p_user_id
    ) OR EXISTS (
        SELECT 1 FROM public.game_players WHERE game_id = p_game_id AND user_id = p_user_id
    );
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION public.is_game_gm(p_game_id UUID, p_user_id UUID)
RETURNS BOOLEAN SECURITY DEFINER AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM public.games WHERE id = p_game_id AND gm_id = p_user_id
    );
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION public.join_game_by_invite_code(p_invite_code TEXT)
RETURNS UUID SECURITY DEFINER AS $$
DECLARE
    v_game_id UUID;
BEGIN
    SELECT id INTO v_game_id FROM public.games WHERE UPPER(invite_code) = UPPER(p_invite_code);
    IF v_game_id IS NULL THEN
        RAISE EXCEPTION 'Code d''invitation invalide.';
    END IF;
    
    INSERT INTO public.game_invitations (game_id, user_id)
    VALUES (v_game_id, auth.uid())
    ON CONFLICT (game_id, user_id) DO NOTHING;
    
    RETURN v_game_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION public.roll_dice(p_game_id UUID, p_sides INT, p_is_secret BOOLEAN DEFAULT FALSE)
RETURNS INT SECURITY DEFINER AS $$
DECLARE
    v_result INT;
    v_sender_name TEXT;
    v_sender_id TEXT;
    v_character_name TEXT;
    v_is_gm BOOLEAN;
    v_msg_type TEXT := 'EVENT';
    v_target_id UUID := NULL;
BEGIN
    v_result := floor(random() * p_sides)::int + 1;
    v_sender_id := auth.uid()::text;
    
    SELECT (gm_id = auth.uid()) INTO v_is_gm FROM public.games WHERE id = p_game_id;
    
    IF v_is_gm THEN
        v_sender_name := 'GM';
    ELSE
        SELECT c.name INTO v_character_name 
        FROM public.characters c
        JOIN public.game_characters gc ON c.id = gc.character_id
        WHERE gc.game_id = p_game_id AND gc.user_id = auth.uid()
        LIMIT 1;
        
        v_sender_name := COALESCE(v_character_name, (
            SELECT COALESCE(raw_user_meta_data->>'name', email) 
            FROM auth.users WHERE id = auth.uid()
        ));
    END IF;
    
    IF p_is_secret AND v_is_gm THEN
        v_msg_type := 'CHAT_PRIVATE';
        v_target_id := auth.uid();
    END IF;
    
    INSERT INTO public.messages (game_id, sender_id, sender_name, content, type, target_id)
    VALUES (
        p_game_id, 
        v_sender_id, 
        v_sender_name, 
        '🎲 d' || p_sides || ' : ' || v_result, 
        v_msg_type, 
        v_target_id
    );
    
    RETURN v_result;
END;
$$ LANGUAGE plpgsql;

-- ----------------------------------------------------
-- VIEWS
-- ----------------------------------------------------

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

CREATE OR REPLACE VIEW public.game_player_users AS
SELECT 
    gp.game_id, 
    gp.user_id, 
    COALESCE(u.raw_user_meta_data->>'name', u.email) as name, 
    u.email, 
    gp.joined_at
FROM public.game_players gp
JOIN auth.users u ON gp.user_id = u.id;

CREATE OR REPLACE VIEW public.game_invitation_users AS
SELECT 
    gi.id, 
    gi.game_id, 
    gi.user_id, 
    COALESCE(u.raw_user_meta_data->>'name', u.email) as user_name, 
    gi.created_at
FROM public.game_invitations gi
JOIN auth.users u ON gi.user_id = u.id;

-- ----------------------------------------------------
-- SECURITY POLICIES (RLS)
-- ----------------------------------------------------

ALTER TABLE public.games ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.templates ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.game_players ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.characters ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.game_characters ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.game_invitations ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.notes ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.messages ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can create games" ON public.games FOR INSERT WITH CHECK (auth.uid() = gm_id);
CREATE POLICY "Members can view the game" ON public.games FOR SELECT USING (auth.uid() = gm_id OR EXISTS (SELECT 1 FROM public.game_players WHERE game_id = id AND user_id = auth.uid()));
CREATE POLICY "GMs can edit or delete their games" ON public.games FOR ALL USING (auth.uid() = gm_id);

CREATE POLICY "GMs can manage game players" ON public.game_players FOR ALL USING (public.is_game_gm(game_id, auth.uid()));
CREATE POLICY "Members can view other game players" ON public.game_players FOR SELECT USING (public.is_game_member(game_id, auth.uid()));

CREATE POLICY "Users can manage their own templates" ON public.templates FOR ALL USING (created_by = auth.uid()) WITH CHECK (created_by = auth.uid());
CREATE POLICY "Users can view public templates or their own" ON public.templates FOR SELECT USING (is_public = TRUE OR created_by = auth.uid());

CREATE POLICY "Members can view characters in their game" ON public.characters FOR SELECT USING (EXISTS (SELECT 1 FROM public.game_characters gc WHERE gc.character_id = id AND public.is_game_member(gc.game_id, auth.uid())));
CREATE POLICY "Any authenticated user can insert characters" ON public.characters FOR INSERT WITH CHECK (auth.role() = 'authenticated');
CREATE POLICY "GMs or assigned players can update characters when game is active" ON public.characters FOR UPDATE USING (EXISTS (SELECT 1 FROM public.game_characters gc JOIN public.games g ON gc.game_id = g.id WHERE gc.character_id = characters.id AND (g.gm_id = auth.uid() OR (gc.user_id = auth.uid() AND g.state = 'ongoing'))));
CREATE POLICY "Only GMs can delete characters" ON public.characters FOR DELETE USING (EXISTS (SELECT 1 FROM public.game_characters gc JOIN public.games g ON gc.game_id = g.id WHERE gc.character_id = characters.id AND g.gm_id = auth.uid()));

CREATE POLICY "GMs can manage game links" ON public.game_characters FOR ALL USING (public.is_game_gm(game_id, auth.uid()));
CREATE POLICY "Members can view game links" ON public.game_characters FOR SELECT USING (public.is_game_member(game_id, auth.uid()));

CREATE POLICY "GMs can view and manage invitations" ON public.game_invitations FOR ALL USING (public.is_game_gm(game_id, auth.uid()));
CREATE POLICY "Users can view or create their own invitations" ON public.game_invitations FOR ALL USING (user_id = auth.uid());

CREATE POLICY "Users can view their own notes or GMs can view all" ON public.notes FOR SELECT USING (user_id = auth.uid() OR public.is_game_gm(game_id, auth.uid()));
CREATE POLICY "GMs can modify notes, players can modify theirs only when game is active" ON public.notes FOR ALL USING (public.is_game_gm(game_id, auth.uid()) OR (user_id = auth.uid() AND EXISTS (SELECT 1 FROM public.games WHERE id = game_id AND state = 'ongoing')));

CREATE POLICY "Members can read global or targeted chat messages" ON public.messages FOR SELECT USING (public.is_game_member(game_id, auth.uid()) AND (target_id IS NULL OR target_id = auth.uid() OR sender_id = auth.uid()::text));
CREATE POLICY "GMs can insert messages; players can insert only when game is active" ON public.messages FOR INSERT WITH CHECK (public.is_game_member(game_id, auth.uid()) AND (public.is_game_gm(game_id, auth.uid()) OR EXISTS (SELECT 1 FROM public.games WHERE id = game_id AND state = 'ongoing')));
CREATE POLICY "GMs can delete messages" ON public.messages FOR DELETE USING (public.is_game_gm(game_id, auth.uid()));
