-- name: CreateUser :one
INSERT INTO users (username, email, hash_password)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT
    *
FROM
    users
WHERE
    username = $1
LIMIT 1;

-- name: GetUserByEmailOrUsername :one
SELECT
    *
FROM
    users
WHERE
    username = sqlc.arg (identifier)
    OR email = sqlc.arg (identifier)
LIMIT 1;

