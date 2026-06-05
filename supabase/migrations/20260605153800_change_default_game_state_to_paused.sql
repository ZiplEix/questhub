-- Change default value of state column in games table to 'paused'
ALTER TABLE public.games ALTER COLUMN state SET DEFAULT 'paused';
