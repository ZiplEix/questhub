-- Migration pour activer le Realtime sur la table characters
ALTER PUBLICATION supabase_realtime ADD TABLE public.characters;
ALTER TABLE public.characters REPLICA IDENTITY FULL;

-- Ajouter la colonne faction à la table board_tokens
ALTER TABLE public.board_tokens ADD COLUMN IF NOT EXISTS faction VARCHAR(50) NOT NULL DEFAULT 'enemy';

-- Mettre à jour la faction par défaut pour les joueurs
UPDATE public.board_tokens bt
SET faction = 'ally'
FROM public.characters c
WHERE bt.character_id = c.id AND c.type = 'PLAYER';
