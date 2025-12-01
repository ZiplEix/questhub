-- +goose Up
-- +goose StatementBegin
CREATE TABLE characters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id TEXT REFERENCES "user"(id) ON DELETE SET NULL, -- NULL = PNJ ou Perso non assigné
    
    name TEXT NOT NULL,
    avatar_url TEXT,

    stats JSONB NOT NULL DEFAULT '{}'::jsonb, 

    inventory JSONB NOT NULL DEFAULT '[]'::jsonb,
    
    is_npc BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_characters_game_id ON characters(game_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE characters;
DROP INDEX idx_characters_game_id;
-- +goose StatementEnd
