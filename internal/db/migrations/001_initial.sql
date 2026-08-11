-- +goose Up
CREATE TABLE history (
    id INTEGER PRIMARY KEY,
    -- The actual command the user has run.
    cmd TEXT NOT NULL,
    -- The absolute path the user was in when running the command.
    cwd TEXT NOT NULL,
    -- A unique ID (UUIDv4) identifying the session, an environment variable
    -- $RUTTER_SESSION_ID is set by the shell hook on session start.
    session TEXT NOT NULL,
    -- Unix timestamp (nanoseconds) at which the command was executed.
    started_at DATETIME NOT NULL,
    -- Duration in nanoseconds of the command's execution, a command that
    -- has been started but not yet finished has a duration of -1.
    duration INTEGER NOT NULL DEFAULT -1,
    -- The exit code of the command. A command that has not yet finished
    -- has an exit code of -1.
    exit INTEGER NOT NULL DEFAULT -1,

    UNIQUE (started_at, session, cmd) -- Idempotency key
);

-- Fast history candidate list, newest first
CREATE INDEX IF NOT EXISTS history_cmd_started_at
ON history (cmd, started_at DESC);

-- Filter by session
CREATE INDEX IF NOT EXISTS history_session_started_at
ON history (session, started_at DESC);

-- Filter by directory
CREATE INDEX IF NOT EXISTS history_cwd_started_at
ON history (cwd, started_at DESC);

-- +goose Down
DROP TABLE history;
