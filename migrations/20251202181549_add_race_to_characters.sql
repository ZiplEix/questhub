-- +goose Up
-- +goose StatementBegin
ALTER TABLE characters ADD COLUMN race TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE characters DROP COLUMN race;
-- +goose StatementEnd
