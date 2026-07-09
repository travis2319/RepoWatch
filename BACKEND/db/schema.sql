CREATE TABLE IF NOT EXISTS repos (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    name                 TEXT NOT NULL,
    owner                TEXT NOT NULL,
    full_name            TEXT NOT NULL,
    url                  TEXT,
    visibility           TEXT,
    private              BOOLEAN DEFAULT 0,
    is_fork              BOOLEAN DEFAULT 0,
    forked_from          TEXT,
    forked_from_owner    TEXT,
    forks_count          INTEGER DEFAULT 0,
    forked_to_count      INTEGER DEFAULT 0,
    stargazers_count     INTEGER DEFAULT 0,
    collaborators_count  INTEGER DEFAULT 0,
    collaborators_list   TEXT,
    who_has_access       TEXT,
    language             TEXT,
    size_kb              INTEGER DEFAULT 0,
    created_at           TEXT,
    updated_at           TEXT,
    pushed_at            TEXT,
    default_branch       TEXT,
    archived             BOOLEAN DEFAULT 0,
    disabled             BOOLEAN DEFAULT 0,
    license              TEXT,
    description          TEXT,
    last_checked         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, owner)
);

CREATE TABLE IF NOT EXISTS collaborators (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    repo_id    INTEGER NOT NULL,
    username   TEXT NOT NULL,
    has_access BOOLEAN NOT NULL DEFAULT 0,
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(repo_id) REFERENCES repos(id)
);

-- CREATE TABLE IF NOT EXISTS repos (
--     id        INTEGER PRIMARY KEY AUTOINCREMENT,
--     name      TEXT NOT NULL,
--     owner     TEXT NOT NULL,
--     full_name TEXT NOT NULL,
--     last_checked TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- CREATE TABLE IF NOT EXISTS collaborators (
--     id         INTEGER PRIMARY KEY AUTOINCREMENT,
--     repo_id    INTEGER NOT NULL,
--     username   TEXT NOT NULL,
--     has_access BOOLEAN NOT NULL DEFAULT 0,
--     checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY(repo_id) REFERENCES repos(id)
-- );