-- +goose Up
-- +goose StatementBegin
ALTER TABLE characters ADD COLUMN money INTEGER DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE characters DROP COLUMN money;
-- +goose StatementEnd
