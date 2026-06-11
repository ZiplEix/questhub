-- Backfill size column in public.media_library from storage.objects metadata
UPDATE public.media_library ml
SET size = COALESCE((o.metadata->>'size')::INTEGER, 0)
FROM storage.objects o
WHERE o.bucket_id = 'images'
  AND ml.url LIKE '%' || o.name;
