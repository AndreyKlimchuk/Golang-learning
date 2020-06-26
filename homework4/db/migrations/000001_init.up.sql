BEGIN;

CREATE TABLE IF NOT EXISTS projects (
    id serial PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL
);

CREATE TABLE IF NOT EXISTS columns (
    id serial PRIMARY KEY,
    project_id integer REFERENCES projects(id),
    name text NOT NULL,
    rank text NOT NULL
);

CREATE UNIQUE INDEX ON columns (project_id, name);

CREATE TABLE IF NOT EXISTS tasks (
    id serial PRIMARY KEY,
    project_id integer REFERENCES projects(id),
    column_id integer REFERENCES columns(id),
    name text NOT NULL,
    description text NOT NULL,
    rank text NOT NULL
);

CREATE INDEX ON tasks (project_id);

CREATE TABLE IF NOT EXISTS comments (
    id serial PRIMARY KEY,
    task_id integer REFERENCES tasks(id),
    text text NOT NULL
);

CREATE INDEX ON comments (task_id);

COMMIT;