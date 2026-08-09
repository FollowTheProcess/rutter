-- name: ListCandidatesInDirectory :many
-- Fetches N history entries scoped to the current directory only.
select
    cmd,
    cwd,
    started_at,
    duration,
    exit
from history
where
    id in (
        select max(h.id) from history h
        where h.cwd = ?
        group by h.cmd
    )
order by started_at desc
limit ?;
