-- +goose Up
-- +goose StatementBegin
CREATE TABLE character_notes (
    character_id UUID PRIMARY KEY REFERENCES characters(id) ON DELETE CASCADE,
    content TEXT DEFAULT ''
);

-- Migrate existing notes if any (optional, but good practice)
-- INSERT INTO character_notes (character_id, content)
-- SELECT id, notes FROM characters WHERE notes IS NOT NULL AND notes != '';

-- ALTER TABLE characters DROP COLUMN notes;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE characters ADD COLUMN notes TEXT DEFAULT '';

-- Restore notes (best effort)
UPDATE characters c
SET notes = cn.content
FROM character_notes cn
WHERE c.id = cn.character_id;

DROP TABLE character_notes;
-- +goose StatementEnd
