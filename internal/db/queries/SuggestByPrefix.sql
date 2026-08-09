-- name: SuggestByPrefix :one
-- Fetches the most recently started command matching the given
-- prefix. Used to power shell autosuggestion.
select cmd
from history
where instr(cmd, sqlc.arg(prefix)) = 1
order by started_at desc
limit 1;
