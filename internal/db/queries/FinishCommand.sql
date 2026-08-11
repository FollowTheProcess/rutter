-- name: FinishCommand :exec
-- Updates the command initially inserted by StartCommand
-- populating duration and exit. Run from the post-exec hook once
-- the command's duration and exit code are known.
UPDATE history
SET duration = ?, exit = ?
WHERE id = ?;
