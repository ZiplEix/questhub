-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    gm_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE, -- Le GM est un user
    invite_code TEXT NOT NULL UNIQUE, -- Code court ex: "A7X-99"
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_games_gm_id ON games(gm_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE games;
DROP INDEX idx_games_gm_id;
-- +goose StatementEnd
