-- Migration pour créer la table des jetons (tokens) sur le plateau de jeu
CREATE TABLE IF NOT EXISTS public.board_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id UUID NOT NULL REFERENCES public.game_boards(id) ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES public.games(id) ON DELETE CASCADE,
    character_id UUID NOT NULL REFERENCES public.characters(id) ON DELETE CASCADE,
    x DOUBLE PRECISION NOT NULL,
    y DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index pour accélérer les jointures
CREATE INDEX IF NOT EXISTS idx_board_tokens_board_id ON public.board_tokens(board_id);
CREATE INDEX IF NOT EXISTS idx_board_tokens_game_id ON public.board_tokens(game_id);

-- Activer RLS
ALTER TABLE public.board_tokens ENABLE ROW LEVEL SECURITY;

-- Politiques de sécurité
CREATE POLICY "Members can view tokens" 
ON public.board_tokens FOR SELECT 
USING (public.is_game_member(game_id, auth.uid()));

CREATE POLICY "Members can update tokens" 
ON public.board_tokens FOR UPDATE
USING (public.is_game_member(game_id, auth.uid()));

CREATE POLICY "GMs can manage tokens" 
ON public.board_tokens FOR ALL 
USING (public.is_game_gm(game_id, auth.uid()));

-- Activer le Realtime pour la table board_tokens
ALTER PUBLICATION supabase_realtime ADD TABLE public.board_tokens;

-- Définir REPLICA IDENTITY à FULL pour que les payloads de suppression (DELETE) incluent tous les champs (requis pour le filtrage par board_id en Realtime)
ALTER TABLE public.board_tokens REPLICA IDENTITY FULL;
