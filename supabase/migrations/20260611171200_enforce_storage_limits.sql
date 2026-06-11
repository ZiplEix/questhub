-- Create function to get current user storage usage
CREATE OR REPLACE FUNCTION public.get_user_storage_usage()
RETURNS BIGINT AS $$
BEGIN
    RETURN COALESCE((
        SELECT SUM((metadata->>'size')::BIGINT)
        FROM storage.objects
        WHERE owner::text = auth.uid()::text
    ), 0);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Create function to check user storage limit before inserts/updates in storage.objects
CREATE OR REPLACE FUNCTION public.check_user_storage_limit()
RETURNS TRIGGER AS $$
DECLARE
    v_total_size BIGINT;
    v_new_size BIGINT;
    v_limit BIGINT := 50 * 1024 * 1024; -- 50 MB in bytes
    v_owner_text TEXT;
BEGIN
    -- If owner is null, allow it (e.g. system operations, admin)
    IF NEW.owner IS NULL THEN
        RETURN NEW;
    END IF;

    v_owner_text := NEW.owner::text;

    -- Extract file size from metadata
    v_new_size := COALESCE((NEW.metadata->>'size')::BIGINT, 0);

    -- Calculate current total storage size for this owner
    SELECT COALESCE(SUM((metadata->>'size')::BIGINT), 0)
    INTO v_total_size
    FROM storage.objects
    WHERE owner::text = v_owner_text
      AND id <> NEW.id;

    IF (v_total_size + v_new_size) > v_limit THEN
        RAISE EXCEPTION 'Limite de stockage de 50 Mo dépassée. Veuillez utiliser des liens ou supprimer d''autres médias.';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Drop trigger if exists
DROP TRIGGER IF EXISTS tr_check_user_storage_limit ON storage.objects;

-- Create trigger on storage.objects
CREATE TRIGGER tr_check_user_storage_limit
BEFORE INSERT OR UPDATE ON storage.objects
FOR EACH ROW
EXECUTE FUNCTION public.check_user_storage_limit();
