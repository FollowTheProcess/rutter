-- name: ListCandidates :many
-- Fetches N history entries scoped globally, i.e. not bound by session
-- or current directory.
select *
from history
where
    id in (
        select max(id) from history
        group by cmd
    )
order by started_at desc
limit ?;
