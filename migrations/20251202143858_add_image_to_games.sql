-- +goose Up
-- +goose StatementBegin
ALTER TABLE games ADD COLUMN image_url TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE games DROP COLUMN image_url;
-- +goose StatementEnd
