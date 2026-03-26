-- name: GetSetting :one
SELECT value FROM server_settings WHERE key = ?;

-- name: UpsertSetting :exec
INSERT INTO server_settings (key, value) VALUES (?, ?)
ON CONFLICT(key) DO UPDATE SET value = excluded.value;

-- name: ListSettings :many
SELECT key, value FROM server_settings ORDER BY key;
