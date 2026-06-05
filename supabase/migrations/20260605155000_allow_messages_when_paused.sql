-- Allow players to insert messages when the game is 'ongoing' OR 'paused'
DROP POLICY IF EXISTS "GMs can insert messages; players can insert only when game is active" ON public.messages;

CREATE POLICY "Members can insert messages when game is ongoing or paused"
  ON public.messages
  FOR INSERT
  WITH CHECK (
    public.is_game_member(game_id, auth.uid())
    AND (
      public.is_game_gm(game_id, auth.uid())
      OR EXISTS (
        SELECT 1 FROM public.games
        WHERE id = game_id AND state IN ('ongoing', 'paused')
      )
    )
  );
