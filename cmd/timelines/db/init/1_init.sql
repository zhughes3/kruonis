CREATE TABLE IF NOT EXISTS groups (
                                      id SERIAL,
                                      title TEXT NOT NULL,
                                      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                      PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS timelines (
                                         id SERIAL,
                                         group_id INTEGER NOT NULL,
                                         title TEXT NOT NULL,
                                         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                         updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                         PRIMARY KEY(id),
                                         FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS events (
                                      id SERIAL,
                                      timeline_id INTEGER NOT NULL,
                                      title TEXT NOT NULL,
                                      timestamp TIMESTAMP NOT NULL,
                                      description TEXT,
                                      content TEXT,
                                      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                      PRIMARY KEY(id),
                                      FOREIGN KEY (timeline_id) REFERENCES timelines(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS tags (
                                    id SERIAL,
                                    tag TEXT NOT NULL,
                                    timeline_id INTEGER NOT NULL,
                                    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                    PRIMARY KEY(id, timeline_id),
                                    FOREIGN KEY (timeline_id) REFERENCES timelines (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL,
                                     email text NOT NULL,
                                     hash text NOT NULL,
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                     PRIMARY KEY (id)
);