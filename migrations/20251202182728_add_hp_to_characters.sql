-- +goose Up
-- +goose StatementBegin
ALTER TABLE characters ADD COLUMN max_hp INT NOT NULL DEFAULT 10;
ALTER TABLE characters ADD COLUMN current_hp INT NOT NULL DEFAULT 10;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE characters DROP COLUMN max_hp;
ALTER TABLE characters DROP COLUMN current_hp;
-- +goose StatementEnd
