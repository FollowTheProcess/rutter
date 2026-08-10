-- name: ListCandidatesInSession :many
-- Fetches N unique history entries, scoped to the current session.
select *
from history
where
    id in (
        select max(h.id) from history h
        where h.session = ?
        group by h.cmd
    )
order by started_at desc
limit ?;
