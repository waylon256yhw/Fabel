-- name: ListCharacters :many
SELECT id, name, description, personality, scenario, first_mes, mes_example
FROM characters
ORDER BY id;

-- name: GetCharacter :one
SELECT id, name, description, personality, scenario, first_mes, mes_example
FROM characters
WHERE id = ?;
