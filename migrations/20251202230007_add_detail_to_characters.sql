-- +goose Up
-- +goose StatementBegin
ALTER TABLE characters
ADD COLUMN initiative INTEGER DEFAULT 0,
ADD COLUMN age TEXT DEFAULT '',
ADD COLUMN height TEXT DEFAULT '',
ADD COLUMN weight TEXT DEFAULT '',
ADD COLUMN max_spells INTEGER DEFAULT 0,
ADD COLUMN spells JSONB DEFAULT '{}'::jsonb,
ADD COLUMN abilities TEXT DEFAULT '',
ADD COLUMN experience INTEGER DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE characters
DROP COLUMN initiative,
DROP COLUMN age,
DROP COLUMN height,
DROP COLUMN weight,
DROP COLUMN max_spells,
DROP COLUMN spells,
DROP COLUMN abilities,
DROP COLUMN experience;
-- +goose StatementEnd
