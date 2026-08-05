-- +goose Up
CREATE TYPE task_status AS enum (
    'TODO',
    'IN_PROGRESS',
    'DONE',
    'CANCELLED'
);

CREATE TYPE task_priority AS enum (
    'LOW',
    'MEDIUM',
    'HIGH',
    'URGENT'
);

CREATE TYPE workspace_role AS enum (
    'OWNER',
    'ADMIN',
    'MEMBER',
    'VIEWER'
);

CREATE TABLE IF NOT EXISTS workspaces (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    name varchar(100) NOT NULL,
    slug varchar(100) UNIQUE NOT NULL,
    owner_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS workspace_members (
    workspace_id uuid NOT NULL REFERENCES workspaces (id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    ROLE workspace_role NOT NULL DEFAULT 'MEMBER',
    joined_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (workspace_id, user_id)
);

CREATE TABLE IF NOT EXISTS tasks (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    workspace_id uuid NOT NULL REFERENCES workspaces (id) ON DELETE CASCADE,
    created_by uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    assignee_id uuid REFERENCES users (id) ON DELETE SET NULL,
    title varchar(255) NOT NULL,
    description text,
    status task_status NOT NULL DEFAULT 'TODO',
    priority task_priority NOT NULL DEFAULT 'MEDIUM',
    due_date timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tasks_workspace_id ON tasks (workspace_id);

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks (status);

-- +goose Down
DROP TABLE IF EXISTS tasks;

DROP TABLE IF EXISTS workspace_members;

DROP TABLE IF EXISTS workspaces;

DROP TYPE IF EXISTS task_priority;

DROP TYPE IF EXISTS task_status;

DROP TYPE IF EXISTS workspace_role;
