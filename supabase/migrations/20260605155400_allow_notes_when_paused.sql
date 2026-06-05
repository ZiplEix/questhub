-- Allow players to read/write their own notes when game is 'ongoing' OR 'paused'
DROP POLICY IF EXISTS "GMs can modify notes, players can modify theirs only when game is active" ON public.notes;

CREATE POLICY "GMs can modify notes, players can modify theirs when game is ongoing or paused"
  ON public.notes
  FOR ALL
  USING (
    public.is_game_gm(game_id, auth.uid())
    OR (
      user_id = auth.uid()
      AND EXISTS (
        SELECT 1 FROM public.games
        WHERE id = game_id AND state IN ('ongoing', 'paused')
      )
    )
  );
