-- Create media_library table
CREATE TABLE IF NOT EXISTS public.media_library (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc'::text, now()) NOT NULL,
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE DEFAULT auth.uid() NOT NULL,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    size INTEGER NOT NULL,
    mime_type TEXT NOT NULL
);

-- Enable RLS
ALTER TABLE public.media_library ENABLE ROW LEVEL SECURITY;

-- Select policy: users can only see their own files
CREATE POLICY "Users can view own media"
ON public.media_library FOR SELECT
TO authenticated
USING (auth.uid() = user_id);

-- Insert policy: users can insert their own files
CREATE POLICY "Users can insert own media"
ON public.media_library FOR INSERT
TO authenticated
WITH CHECK (auth.uid() = user_id);

-- Delete policy: users can delete their own files
CREATE POLICY "Users can delete own media"
ON public.media_library FOR DELETE
TO authenticated
USING (auth.uid() = user_id);
