-- name: ListCharacters :many
SELECT id, name, description, personality, scenario, first_mes, mes_example
FROM characters
WHERE user_id = ?
ORDER BY id;

-- name: GetCharacter :one
SELECT id, name, description, personality, scenario, first_mes, mes_example
FROM characters
WHERE id = ? AND user_id = ?;

-- name: BackfillCharacters :exec
UPDATE characters SET user_id = ? WHERE user_id IS NULL;
