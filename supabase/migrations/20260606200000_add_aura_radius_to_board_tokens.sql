-- Add aura_radius column to board_tokens for area-of-effect rendering
ALTER TABLE public.board_tokens ADD COLUMN IF NOT EXISTS aura_radius DOUBLE PRECISION;
