-- Create marketplace_templates view to securely join auth.users for author names
CREATE OR REPLACE VIEW public.marketplace_templates AS
SELECT 
    t.id,
    t.created_by,
    COALESCE(u.raw_user_meta_data->>'name', u.email) as author_name,
    t.name,
    t.description,
    t.type,
    t.data,
    t.is_public,
    t.uses,
    t.created_at
FROM public.templates t
LEFT JOIN auth.users u ON t.created_by = u.id
WHERE t.is_public = TRUE OR t.created_by = auth.uid();

-- Grant permissions to authenticated users to select from this view
GRANT SELECT ON public.marketplace_templates TO authenticated;

-- Create function to atomically increment template uses
CREATE OR REPLACE FUNCTION public.increment_template_uses(p_template_id UUID)
RETURNS VOID SECURITY DEFINER AS $$
BEGIN
    UPDATE public.templates
    SET uses = COALESCE(uses, 0) + 1
    WHERE id = p_template_id;
END;
$$ LANGUAGE plpgsql;

-- Grant execution permissions
GRANT EXECUTE ON FUNCTION public.increment_template_uses TO authenticated;
