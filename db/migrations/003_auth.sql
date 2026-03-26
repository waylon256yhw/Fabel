-- +goose Up

-- Users
CREATE TABLE users (
    id            TEXT PRIMARY KEY,                -- ULID
    username      TEXT NOT NULL UNIQUE COLLATE NOCASE,
    display_name  TEXT NOT NULL DEFAULT '',
    password_hash TEXT NOT NULL,
    role          TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user')),
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Sessions (alexedwards/scs format)
CREATE TABLE sessions (
    token  TEXT PRIMARY KEY,
    data   BLOB NOT NULL,
    expiry REAL NOT NULL
);
CREATE INDEX idx_sessions_expiry ON sessions(expiry);

-- Server settings (key-value)
CREATE TABLE server_settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
INSERT INTO server_settings (key, value) VALUES ('allow_registration', 'false');

-- Add user_id to existing tables (nullable for seed data backwards compat)
ALTER TABLE characters ADD COLUMN user_id TEXT REFERENCES users(id);
ALTER TABLE presets ADD COLUMN user_id TEXT REFERENCES users(id);
ALTER TABLE conversations ADD COLUMN user_id TEXT REFERENCES users(id);

CREATE INDEX idx_characters_user ON characters(user_id);
CREATE INDEX idx_presets_user ON presets(user_id);
CREATE INDEX idx_conversations_user ON conversations(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_conversations_user;
DROP INDEX IF EXISTS idx_presets_user;
DROP INDEX IF EXISTS idx_characters_user;

-- SQLite cannot DROP COLUMN, so we leave user_id columns in place on rollback.
-- Use make reset for a clean slate.

DROP TABLE IF EXISTS server_settings;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
