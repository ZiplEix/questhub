-- +goose Up
-- +goose StatementBegin
CREATE TABLE game_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Un joueur ne peut avoir qu'une seule invitation en attente pour une même partie
    UNIQUE(game_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE game_invitations;
-- +goose StatementEnd
