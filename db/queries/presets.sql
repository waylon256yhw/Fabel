-- name: ListPresets :many
SELECT id, name, model, temperature, max_tokens, system_prompt
FROM presets
ORDER BY id;
