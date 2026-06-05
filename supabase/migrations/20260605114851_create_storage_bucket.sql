-- Create storage bucket 'images' for campaign uploads
INSERT INTO storage.buckets (id, name, public)
VALUES ('images', 'images', true)
ON CONFLICT (id) DO NOTHING;

-- Set up RLS policies for storage objects in the 'images' bucket
-- Give public read access to all files in the 'images' bucket
CREATE POLICY "Public Read Access"
ON storage.objects FOR SELECT
USING (bucket_id = 'images');

-- Give authenticated users permission to upload to the 'images' bucket
CREATE POLICY "Authenticated Upload Access"
ON storage.objects FOR INSERT
TO authenticated
WITH CHECK (bucket_id = 'images');

-- Give owners permission to update/delete their own objects in the 'images' bucket
CREATE POLICY "Owner Update/Delete Access"
ON storage.objects FOR ALL
TO authenticated
USING (bucket_id = 'images' AND auth.uid()::text = owner::text);
