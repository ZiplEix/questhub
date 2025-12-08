-- +goose Up
-- +goose StatementBegin
ALTER TABLE games ADD COLUMN state VARCHAR(50) DEFAULT 'ongoing';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE games DROP COLUMN state;
-- +goose StatementEnd
