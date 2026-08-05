-- name: CreateTask :one
INSERT INTO tasks (workspace_id, created_by, assignee_id, title, description, status, priority, due_date)
    VALUES ($1, $2, $3, $4, $5, COALESCE($6::task_status, 'TODO'::task_status), COALESCE($7::task_priority, 'MEDIUM'::task_priority), $8)
RETURNING
    *;

-- name: GetTaskByID :one
SELECT
    *
FROM
    tasks
WHERE
    id = $1
LIMIT 1;

-- name: ListTasksByWorkspace :many
SELECT
    *
FROM
    tasks
WHERE
    workspace_id = $1
ORDER BY
    created_at DESC;

-- name: ListTaskByAssignee :many
SELECT
    *
FROM
    tasks
WHERE
    assignee_id = $1
    AND workspace_id = $2
ORDER BY
    created_at DESC;

-- name: UpdateTask :one
UPDATE
    tasks
SET
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    assignee_id = COALESCE($4, assignee_id),
    status = COALESCE($5::task_status, status),
    priority = COALESCE($6::task_priority, priority),
    due_date = COALESCE($7, due_date),
    updated_at = NOW()
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;
