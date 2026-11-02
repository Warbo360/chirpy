-- name: GetChirps :many

SELECT
    *
FROM
    chirps
ORDER BY
    created_at;

-- name: GetChirpsByUser :many

SELECT
    *
FROM
    chirps
WHERE
    user_id = $1;
