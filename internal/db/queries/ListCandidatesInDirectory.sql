-- name: ListCandidatesInDirectory :many
-- Fetches N unique history entries scoped to the current directory only.
select *
from history
where
    id in (
        select max(h.id) from history h
        where h.cwd = ?
        group by h.cmd
    )
order by started_at desc
limit ?;
