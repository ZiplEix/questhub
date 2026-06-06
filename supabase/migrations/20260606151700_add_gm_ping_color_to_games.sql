-- Add gm_ping_color to games
ALTER TABLE public.games ADD COLUMN IF NOT EXISTS gm_ping_color TEXT NOT NULL DEFAULT '#E07A5F';
