-- Migration to create story tables

-- 1. Create story_folders table
CREATE TABLE IF NOT EXISTS public.story_folders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES public.games(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    is_visible BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Create story_pages table
CREATE TABLE IF NOT EXISTS public.story_pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES public.games(id) ON DELETE CASCADE,
    folder_id UUID REFERENCES public.story_folders(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT DEFAULT '',
    is_visible BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_story_folders_game_id ON public.story_folders(game_id);
CREATE INDEX IF NOT EXISTS idx_story_pages_game_id ON public.story_pages(game_id);
CREATE INDEX IF NOT EXISTS idx_story_pages_folder_id ON public.story_pages(folder_id);

-- 4. Enable Row Level Security (RLS)
ALTER TABLE public.story_folders ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.story_pages ENABLE ROW LEVEL SECURITY;

-- 5. RLS Policies for story_folders
CREATE POLICY "GMs can manage story folders" ON public.story_folders
    FOR ALL USING (public.is_game_gm(game_id, auth.uid()));

CREATE POLICY "Members can view visible story folders" ON public.story_folders
    FOR SELECT USING (public.is_game_member(game_id, auth.uid()) AND is_visible = TRUE);

-- 6. RLS Policies for story_pages
CREATE POLICY "GMs can manage story pages" ON public.story_pages
    FOR ALL USING (public.is_game_gm(game_id, auth.uid()));

CREATE POLICY "Members can view visible story pages" ON public.story_pages
    FOR SELECT USING (
        public.is_game_member(game_id, auth.uid()) 
        AND is_visible = TRUE 
        AND (
            folder_id IS NULL 
            OR folder_id IN (SELECT id FROM public.story_folders WHERE is_visible = TRUE)
        )
    );

-- 7. Enable Supabase Realtime
ALTER PUBLICATION supabase_realtime ADD TABLE public.story_folders;
ALTER PUBLICATION supabase_realtime ADD TABLE public.story_pages;

ALTER TABLE public.story_folders REPLICA IDENTITY FULL;
ALTER TABLE public.story_pages REPLICA IDENTITY FULL;
