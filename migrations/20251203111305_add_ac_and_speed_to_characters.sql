-- +goose Up
-- +goose StatementBegin
ALTER TABLE characters ADD COLUMN armor_class INTEGER DEFAULT 10;
ALTER TABLE characters ADD COLUMN speed INTEGER DEFAULT 30;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE characters DROP COLUMN armor_class;
ALTER TABLE characters DROP COLUMN speed;
-- +goose StatementEnd
