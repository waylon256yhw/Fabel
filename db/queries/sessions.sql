-- name: GetSession :one
SELECT data, expiry FROM sessions WHERE token = ?;

-- name: UpsertSession :exec
INSERT INTO sessions (token, data, expiry) VALUES (?, ?, ?)
ON CONFLICT(token) DO UPDATE SET data = excluded.data, expiry = excluded.expiry;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token = ?;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expiry < ?;

-- name: GetAllSessions :many
SELECT token, data FROM sessions WHERE expiry >= ?;
