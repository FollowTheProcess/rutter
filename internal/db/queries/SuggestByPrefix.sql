-- name: SuggestByPrefix :one
-- Fetches the most recently started command matching the given
-- prefix. Used to power shell autosuggestion.
SELECT cmd
FROM history
WHERE
    cmd >= sqlc.arg(prefix)
    AND cmd < sqlc.arg(prefix) || char(1114111)
ORDER BY started_at DESC
LIMIT 1;
