-- name: UpsertResource :exec
INSERT INTO resources (id, user_id, kind, title, content, tags, meta, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
ON CONFLICT(id) DO UPDATE SET
    title = excluded.title,
    content = excluded.content,
    tags = excluded.tags,
    meta = excluded.meta,
    updated_at = excluded.updated_at;

-- name: DeleteResource :exec
DELETE FROM resources WHERE id = ?;

-- name: DeleteResourcesByKindAndUser :exec
DELETE FROM resources WHERE kind = ? AND user_id = ?;

-- name: UpsertRelation :exec
INSERT OR IGNORE INTO resource_relations (src_id, dst_id, relation)
VALUES (?, ?, ?);

-- name: DeleteRelation :exec
DELETE FROM resource_relations WHERE src_id = ? AND dst_id = ? AND relation = ?;

-- name: DeleteRelationsBySrc :exec
DELETE FROM resource_relations WHERE src_id = ?;

-- name: DeleteRelationsByDst :exec
DELETE FROM resource_relations WHERE dst_id = ?;

-- name: GetRelationsBySrc :many
SELECT src_id, dst_id, relation FROM resource_relations WHERE src_id = ?;

-- name: GetRelationsByDst :many
SELECT src_id, dst_id, relation FROM resource_relations WHERE dst_id = ?;
