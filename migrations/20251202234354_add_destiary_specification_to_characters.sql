-- +goose Up
-- +goose StatementBegin
ALTER TABLE characters ADD COLUMN type VARCHAR(50) DEFAULT 'PLAYER';
ALTER TABLE characters ADD COLUMN sub_race VARCHAR(255);

-- Update existing records based on is_npc flag
UPDATE characters SET type = 'NPC' WHERE is_npc = true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE characters DROP COLUMN type;
ALTER TABLE characters DROP COLUMN sub_race;
-- +goose StatementEnd
