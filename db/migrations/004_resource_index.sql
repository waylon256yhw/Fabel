-- +goose Up

-- Resource index: projection of all domain entities for unified search.
CREATE TABLE resources (
    rowid      INTEGER PRIMARY KEY,
    id         TEXT NOT NULL UNIQUE,        -- shares ULID with domain entity
    user_id    TEXT NOT NULL REFERENCES users(id),
    kind       TEXT NOT NULL,               -- character / preset / conversation / lorebook / lorebook_entry / persona
    title      TEXT NOT NULL DEFAULT '',
    content    TEXT NOT NULL DEFAULT '',     -- full text for search
    tags       TEXT NOT NULL DEFAULT '',     -- space-separated
    meta       TEXT NOT NULL DEFAULT '{}',   -- JSON metadata
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);
CREATE INDEX idx_resources_user_kind ON resources(user_id, kind);

-- FTS5 content-sync table (trigram tokenizer for CJK + substring matching).
CREATE VIRTUAL TABLE resource_fts USING fts5(
    title, content, tags,
    content='resources',
    content_rowid='rowid',
    tokenize='trigram'
);

-- Triggers to keep FTS in sync with resources table.
-- +goose StatementBegin
CREATE TRIGGER resources_ai AFTER INSERT ON resources BEGIN
    INSERT INTO resource_fts(rowid, title, content, tags)
    VALUES (new.rowid, new.title, new.content, new.tags);
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER resources_ad AFTER DELETE ON resources BEGIN
    INSERT INTO resource_fts(resource_fts, rowid, title, content, tags)
    VALUES ('delete', old.rowid, old.title, old.content, old.tags);
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER resources_au AFTER UPDATE ON resources BEGIN
    INSERT INTO resource_fts(resource_fts, rowid, title, content, tags)
    VALUES ('delete', old.rowid, old.title, old.content, old.tags);
    INSERT INTO resource_fts(rowid, title, content, tags)
    VALUES (new.rowid, new.title, new.content, new.tags);
END;
-- +goose StatementEnd

-- Resource relations: directed edges between domain entities.
CREATE TABLE resource_relations (
    src_id   TEXT NOT NULL,
    dst_id   TEXT NOT NULL,
    relation TEXT NOT NULL,  -- contains / attached_to / uses
    PRIMARY KEY (src_id, dst_id, relation)
);
CREATE INDEX idx_resource_relations_dst ON resource_relations(dst_id, relation);

-- +goose Down
DROP TABLE IF EXISTS resource_relations;
DROP TRIGGER IF EXISTS resources_au;
DROP TRIGGER IF EXISTS resources_ad;
DROP TRIGGER IF EXISTS resources_ai;
DROP TABLE IF EXISTS resource_fts;
DROP TABLE IF EXISTS resources;
