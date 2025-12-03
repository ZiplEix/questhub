-- +goose Up
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    sender_id TEXT NOT NULL, -- Can be a UserID or "GM" or "System"
    sender_name TEXT NOT NULL,
    content TEXT NOT NULL,
    type TEXT NOT NULL, -- "CHAT_GLOBAL", "CHAT_PRIVATE", "EVENT"
    target_id TEXT, -- UserID for private messages, NULL for global
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_messages_game_id ON messages(game_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);

-- +goose Down
DROP TABLE messages;
