-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING id, created_at, updated_at, email;

-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    users.email = $1;

-- name: GetUserFromRefreshToken :one

SELECT
    users.*
FROM
    refresh_tokens
    INNER JOIN users ON users.id = refresh_tokens.user_id
WHERE
    refresh_tokens.token = $1;

-- name: UpdateUserEmailAndPassword :one

UPDATE users
SET email = $1,
    hashed_password = $2
WHERE id = $3
RETURNING *;

-- name: UpdateUserIsRedToTrue :exec

UPDATE users
SET is_red = TRUE
WHERE id = $1;
