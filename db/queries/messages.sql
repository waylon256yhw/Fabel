-- name: GetMessagesByConversation :many
SELECT id, conversation_id, role, content, parent_id, active_child_id, created_at
FROM messages
WHERE conversation_id = ?
ORDER BY id;

-- name: AppendMessage :exec
INSERT INTO messages (conversation_id, role, content, parent_id)
VALUES (?, ?, ?, ?);
