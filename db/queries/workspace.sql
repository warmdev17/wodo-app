-- name: CreateWorkspace :one
INSERT INTO workspaces (name, slug, owner_id)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: AddWorkspaceMember :one
INSERT INTO workspace_members (workspace_id, user_id, role)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: GetWorkspaceByID :one
SELECT
    *
FROM
    workspaces
WHERE
    id = $1
LIMIT 1;

-- name: GetWorkspaceBySlug :one
SELECT
    *
FROM
    workspaces
WHERE
    slug = $1
LIMIT 1;

-- name: GetWorkspacesByUserID :many
SELECT
    w.*,
    wm.role
FROM
    workspaces w
    JOIN workspace_members wm ON w.id = wm.workspace_id
WHERE
    wm.user_id = $1
ORDER BY
    w.created_at DESC;

-- name: GetWorkspaceMember :one
SELECT
    *
FROM
    workspace_members
WHERE
    workspace_id = $1
    AND user_id = $2
LIMIT 1;

