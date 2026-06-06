-- Add scale system to game_boards for proper token sizing
ALTER TABLE public.game_boards
  ADD COLUMN pixels_per_cell INT NOT NULL DEFAULT 70,
  ADD COLUMN grid_offset_x DOUBLE PRECISION NOT NULL DEFAULT 0,
  ADD COLUMN grid_offset_y DOUBLE PRECISION NOT NULL DEFAULT 0;
