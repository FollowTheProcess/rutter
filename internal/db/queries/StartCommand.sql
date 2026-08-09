-- name: StartCommand :one
-- Inserts the command from the pre-exec hook, before duration
-- and exit are known. Returns the id so FinishCommand can
-- update it with exit and duration.
insert into history (cmd, cwd, session, started_at)
values (?, ?, ?, ?)
returning id;
