-- Create support tickets table
CREATE TABLE IF NOT EXISTS public.support_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc'::text, now()) NOT NULL,
    user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    email TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('CONTACT', 'BUG', 'RECLAMATION')),
    subject TEXT NOT NULL,
    message TEXT NOT NULL,
    attachment_url TEXT,
    status TEXT DEFAULT 'PENDING' NOT NULL CHECK (status IN ('PENDING', 'IN_PROGRESS', 'RESOLVED', 'CLOSED'))
);

-- Enable RLS on support_tickets
ALTER TABLE public.support_tickets ENABLE ROW LEVEL SECURITY;

-- Allow anyone (anonymous and authenticated) to submit a ticket
CREATE POLICY "Public Support Ticket Insert"
ON public.support_tickets FOR INSERT
WITH CHECK (true);

-- Allow authenticated users to view only their own submitted tickets
CREATE POLICY "Authenticated Support Ticket Select Own"
ON public.support_tickets FOR SELECT
TO authenticated
USING (auth.uid() = user_id);

-- Storage bucket permissions: Allow anyone to upload feedback attachments
-- under the 'feedback/' prefix in the 'images' bucket.
CREATE POLICY "Public Upload Feedback Access"
ON storage.objects FOR INSERT
TO public
WITH CHECK (bucket_id = 'images' AND (position('feedback/' in name) = 1));
