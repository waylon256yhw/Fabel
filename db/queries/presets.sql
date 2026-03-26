-- name: ListPresets :many
SELECT id, name, model, temperature, max_tokens, system_prompt
FROM presets
WHERE user_id = ?
ORDER BY id;

-- name: GetPreset :one
SELECT id, name, model, temperature, max_tokens, system_prompt
FROM presets
WHERE id = ? AND user_id = ?;

-- name: BackfillPresets :exec
UPDATE presets SET user_id = ? WHERE user_id IS NULL;
