-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    gm_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE, -- Le GM est un user
    invite_code TEXT NOT NULL UNIQUE, -- Code court ex: "A7X-99"
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    image_url TEXT,
    state VARCHAR(50) DEFAULT 'ongoing'
);

CREATE INDEX IF NOT EXISTS idx_games_gm_id ON games(gm_id);

-- TEMPLATES: Blueprints for shareable content
CREATE TABLE IF NOT EXISTS templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_by TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL, -- 'CHARACTER', 'NPC', 'MONSTER', 'CREATURE'
    data JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_public BOOLEAN DEFAULT FALSE,
    uses INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_templates_type ON templates(type);
CREATE INDEX IF NOT EXISTS idx_templates_public ON templates(is_public);

-- TABLE DE LIAISON JOUEURS <-> PARTIES
CREATE TABLE IF NOT EXISTS game_players (
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Un joueur ne peut rejoindre qu'une seule fois la même partie
    PRIMARY KEY (game_id, user_id)
);

CREATE TABLE IF NOT EXISTS characters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    -- Removed game_id and user_id to decouple

    name TEXT NOT NULL,
    avatar_url TEXT,

    stats JSONB NOT NULL DEFAULT '{}'::jsonb, 
    inventory JSONB NOT NULL DEFAULT '[]'::jsonb,
    
    is_npc BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    race TEXT NOT NULL DEFAULT '',
    max_hp INT NOT NULL DEFAULT 10,
    current_hp INT NOT NULL DEFAULT 10,
    money INTEGER DEFAULT 0,

    initiative INTEGER DEFAULT 0,
    age TEXT DEFAULT '',
    height TEXT DEFAULT '',
    weight TEXT DEFAULT '',
    max_spells INTEGER DEFAULT 0,
    spells JSONB DEFAULT '{}'::jsonb,
    abilities TEXT DEFAULT '',
    experience INTEGER DEFAULT 0,

    type VARCHAR(50) DEFAULT 'PLAYER',
    sub_race VARCHAR(255),

    armor_class INTEGER DEFAULT 10,
    speed INTEGER DEFAULT 30
);

-- New Table: Links a character instance to a game and potentially a user
CREATE TABLE IF NOT EXISTS game_characters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    character_id UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    user_id TEXT REFERENCES "user"(id) ON DELETE SET NULL, -- NULL if unassigned or pure NPC managed by GM?
    
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(game_id, character_id) -- A character instance can only be in a game once
);

CREATE INDEX IF NOT EXISTS idx_game_characters_game_id ON game_characters(game_id);
CREATE INDEX IF NOT EXISTS idx_game_characters_user_id ON game_characters(user_id);

CREATE TABLE IF NOT EXISTS game_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Un joueur ne peut avoir qu'une seule invitation en attente pour une même partie
    UNIQUE(game_id, user_id)
);

-- Renamed from character_notes to notes, linked to Game + User
CREATE TABLE IF NOT EXISTS notes (
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    content TEXT DEFAULT '',
    
    PRIMARY KEY (game_id, user_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    sender_id TEXT NOT NULL, -- Can be a UserID or "GM" or "System"
    sender_name TEXT NOT NULL,
    content TEXT NOT NULL,
    type TEXT NOT NULL, -- "CHAT_GLOBAL", "CHAT_PRIVATE", "EVENT"
    target_id TEXT, -- UserID for private messages, NULL for global
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_messages_game_id ON messages(game_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS game_invitations;
DROP TABLE IF EXISTS game_characters;
DROP TABLE IF EXISTS characters;
DROP TABLE IF EXISTS game_players;
DROP TABLE IF EXISTS templates;
DROP TABLE IF EXISTS games;
-- +goose StatementEnd
