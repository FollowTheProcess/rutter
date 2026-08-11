-- name: ListCandidatesInDirectory :many
-- Fetches N unique history entries scoped to the current directory only.
SELECT *
FROM history
WHERE
    id IN (
        SELECT max(h.id) FROM history h
        WHERE h.cwd = ?
        GROUP BY h.cmd
    )
ORDER BY started_at DESC
LIMIT ?;
