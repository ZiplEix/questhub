-- Fix the shadowing issue in "Members can view characters in their game" SELECT policy on public.characters
DROP POLICY IF EXISTS "Members can view characters in their game" ON public.characters;

CREATE POLICY "Members can view characters in their game" ON public.characters 
FOR SELECT 
USING (
    EXISTS (
        SELECT 1 FROM public.game_characters gc 
        WHERE gc.character_id = public.characters.id 
        AND public.is_game_member(gc.game_id, auth.uid())
    )
);
