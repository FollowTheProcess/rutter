-- +goose Up
create table history (
    id integer primary key,
    -- The actual command the user has run.
    cmd text not null,
    -- The absolute path the user was in when running the command.
    cwd text not null,
    -- A unique ID (UUIDv4) identifying the session, an environment variable
    -- $RUTTER_SESSION_ID is set by the shell hook on session start.
    session text not null,
    -- Unix timestamp (nanoseconds) at which the command was executed.
    started_at datetime not null,
    -- Duration in nanoseconds of the command's execution, a command that
    -- has been started but not yet finished has a duration of -1.
    duration integer not null default -1,
    -- The exit code of the command. A command that has not yet finished
    -- has an exit code of -1.
    exit integer not null default -1,

    unique (started_at, session, cmd) -- Idempotency key
);

-- Fast history candidate list, newest first
create index if not exists history_cmd_started_at
on history (cmd, started_at desc);

-- Filter by session
create index if not exists history_session_started_at
on history (session, started_at desc);

-- Filter by directory
create index if not exists history_cwd_started_at
on history (cwd, started_at desc);

-- +goose Down
drop table history;
