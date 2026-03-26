-- name: GetConversation :one
SELECT id, character_id, preset_id, created_at
FROM conversations
WHERE id = ? AND user_id = ?;

-- name: CreateConversation :execlastid
INSERT INTO conversations (character_id, preset_id, user_id)
VALUES (?, ?, ?);

-- name: GetFirstConversationID :one
SELECT id FROM conversations WHERE user_id = ? ORDER BY id LIMIT 1;

-- name: BackfillConversations :exec
UPDATE conversations SET user_id = ? WHERE user_id IS NULL;
