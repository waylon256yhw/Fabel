-- +goose Up
CREATE TABLE IF NOT EXISTS characters (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    name         TEXT    NOT NULL,
    description  TEXT    NOT NULL DEFAULT '',
    personality  TEXT    NOT NULL DEFAULT '',
    scenario     TEXT    NOT NULL DEFAULT '',
    first_mes    TEXT    NOT NULL DEFAULT '',
    mes_example  TEXT    NOT NULL DEFAULT '',
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS presets (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    name          TEXT    NOT NULL,
    model         TEXT    NOT NULL DEFAULT 'gpt-4o-mini',
    temperature   REAL    NOT NULL DEFAULT 1.0,
    max_tokens    INTEGER NOT NULL DEFAULT 1000,
    system_prompt TEXT    NOT NULL DEFAULT '',
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS conversations (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL REFERENCES characters(id),
    preset_id    INTEGER NOT NULL REFERENCES presets(id),
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    conversation_id INTEGER NOT NULL REFERENCES conversations(id),
    role            TEXT    NOT NULL CHECK(role IN ('user', 'assistant', 'system')),
    content         TEXT    NOT NULL,
    parent_id       INTEGER REFERENCES messages(id),
    active_child_id INTEGER REFERENCES messages(id),
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS presets;
DROP TABLE IF EXISTS characters;
