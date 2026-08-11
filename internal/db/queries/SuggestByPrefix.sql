-- name: SuggestByPrefix :one
-- Fetches the most recently started command matching the given
-- prefix. Used to power shell autosuggestion.
SELECT cmd
FROM history
WHERE instr(cmd, sqlc.arg(prefix)) = 1
ORDER BY started_at DESC
LIMIT 1;
