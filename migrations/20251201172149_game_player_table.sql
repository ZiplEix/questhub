-- +goose Up
-- +goose StatementBegin
-- TABLE DE LIAISON JOUEURS <-> PARTIES
CREATE TABLE game_players (
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Un joueur ne peut rejoindre qu'une seule fois la même partie
    PRIMARY KEY (game_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE game_players;
-- +goose StatementEnd
