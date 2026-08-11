-- name: ListCandidatesInSession :many
-- Fetches N unique history entries, scoped to the current session.
SELECT *
FROM history
WHERE
    id IN (
        SELECT max(h.id) FROM history h
        WHERE h.session = ?
        GROUP BY h.cmd
    )
ORDER BY started_at DESC
LIMIT ?;
