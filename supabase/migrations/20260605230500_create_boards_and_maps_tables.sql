-- Création de la table des plateaux
CREATE TABLE IF NOT EXISTS public.game_boards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES public.games(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index pour game_boards
CREATE INDEX IF NOT EXISTS idx_game_boards_game_id ON public.game_boards(game_id);

-- Création de la table des cartes (reliées à un plateau et à la partie)
CREATE TABLE IF NOT EXISTS public.maps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id UUID NOT NULL REFERENCES public.game_boards(id) ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES public.games(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    image_url TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index pour maps
CREATE INDEX IF NOT EXISTS idx_maps_board_id ON public.maps(board_id);
CREATE INDEX IF NOT EXISTS idx_maps_game_id ON public.maps(game_id);

-- Activer RLS sur les deux tables
ALTER TABLE public.game_boards ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.maps ENABLE ROW LEVEL SECURITY;

-- Politiques de sécurité pour game_boards
CREATE POLICY "Members can view game boards" 
ON public.game_boards FOR SELECT 
USING (public.is_game_member(game_id, auth.uid()));

CREATE POLICY "GMs can manage game boards" 
ON public.game_boards FOR ALL 
USING (public.is_game_gm(game_id, auth.uid()));

-- Politiques de sécurité pour maps (simplifiées grâce à game_id)
CREATE POLICY "Members can view maps" 
ON public.maps FOR SELECT 
USING (public.is_game_member(game_id, auth.uid()));

CREATE POLICY "GMs can manage maps" 
ON public.maps FOR ALL 
USING (public.is_game_gm(game_id, auth.uid()));

-- Activer Supabase Realtime pour les deux tables
ALTER PUBLICATION supabase_realtime ADD TABLE public.game_boards;
ALTER PUBLICATION supabase_realtime ADD TABLE public.maps;
