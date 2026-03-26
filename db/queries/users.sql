-- name: CreateUser :exec
INSERT INTO users (id, username, display_name, password_hash, role)
VALUES (?, ?, ?, ?, ?);

-- name: GetUserByUsername :one
SELECT id, username, display_name, password_hash, role, created_at
FROM users
WHERE username = ?;

-- name: GetUserByID :one
SELECT id, username, display_name, password_hash, role, created_at
FROM users
WHERE id = ?;

-- name: ListUsers :many
SELECT id, username, display_name, role, created_at
FROM users
ORDER BY created_at;

-- name: UpdateUserRole :exec
UPDATE users SET role = ? WHERE id = ?;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;
