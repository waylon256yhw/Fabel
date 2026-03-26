-- name: GetConversation :one
SELECT id, character_id, preset_id, created_at
FROM conversations
WHERE id = ?;

-- name: CreateConversation :execlastid
INSERT INTO conversations (character_id, preset_id)
VALUES (?, ?);

-- name: GetFirstConversationID :one
SELECT id FROM conversations ORDER BY id LIMIT 1;
