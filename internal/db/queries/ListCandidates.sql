-- name: ListCandidates :many
-- Fetches N unique history entries scoped globally, i.e. not bound by session
-- or current directory.
select *
from history
where
    id in (
        select max(h.id) from history h
        group by h.cmd
    )
order by started_at desc
limit ?;
