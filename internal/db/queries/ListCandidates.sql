-- name: ListCandidates :many
-- Fetches N unique history entries scoped globally, i.e. not bound by session
-- or current directory.
SELECT *
FROM history
WHERE
    id IN (
        SELECT max(h.id) FROM history h
        GROUP BY h.cmd
    )
ORDER BY started_at DESC
LIMIT ?;
