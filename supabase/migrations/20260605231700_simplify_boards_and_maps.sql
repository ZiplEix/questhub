-- Supprimer la table maps devenue obsolète
DROP TABLE IF EXISTS public.maps CASCADE;

-- Ajouter l'image_url directement à la table game_boards
ALTER TABLE public.game_boards 
ADD COLUMN IF NOT EXISTS image_url TEXT;
