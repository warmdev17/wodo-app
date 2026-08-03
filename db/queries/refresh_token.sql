-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: GetRefreshToken :one
SELECT
    *
FROM
    refresh_tokens
WHERE
    token = $1
LIMIT 1;

-- name: GetRefreshTokenByToken :one
SELECT
    id,
    user_id,
    token,
    is_revoked,
    expires_at,
    created_at
FROM
    refresh_tokens
WHERE
    token = $1;

-- name: RevokeToken :exec
UPDATE
    refresh_tokens
SET
    is_revoked = TRUE
WHERE
    token = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = $1;

-- name: DeleteUserRefreshToken :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

