-- name: ListCandidatesInSession :many
-- Fetches N history entries, scoped to the current session.
select *
from history
where session = ?
order by started_at desc
limit ?;
