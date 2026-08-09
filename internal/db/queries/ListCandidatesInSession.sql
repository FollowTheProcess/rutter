-- name: ListCandidatesInSession :many
-- Fetches N history entries, scoped to the current session.
select
    cmd,
    cwd,
    started_at,
    duration,
    exit
from history
where session = ?
order by started_at desc
limit ?;
