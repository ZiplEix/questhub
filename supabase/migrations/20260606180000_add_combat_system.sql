-- Migration pour ajouter le système de combat et d'initiative
ALTER TABLE public.game_boards ADD COLUMN IF NOT EXISTS combat_active BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE public.game_boards ADD COLUMN IF NOT EXISTS combat_round INT NOT NULL DEFAULT 1;
ALTER TABLE public.game_boards ADD COLUMN IF NOT EXISTS combat_active_token_id UUID;
ALTER TABLE public.game_boards ADD COLUMN IF NOT EXISTS hide_monster_stats BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE public.board_tokens ADD COLUMN IF NOT EXISTS in_combat BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE public.board_tokens ADD COLUMN IF NOT EXISTS initiative INT NOT NULL DEFAULT 0;
ALTER TABLE public.board_tokens ADD COLUMN IF NOT EXISTS is_hidden BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE public.characters ADD COLUMN IF NOT EXISTS conditions JSONB NOT NULL DEFAULT '[]'::jsonb;

-- Contrainte de clé étrangère pour le token actif
ALTER TABLE public.game_boards 
DROP CONSTRAINT IF EXISTS fk_game_boards_combat_active_token,
ADD CONSTRAINT fk_game_boards_combat_active_token 
FOREIGN KEY (combat_active_token_id) REFERENCES public.board_tokens(id) ON DELETE SET NULL;
