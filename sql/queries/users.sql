-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE name = $1 LIMIT 1;

-- name: ResetDatabase :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT name,
case 
WHEN name = $1 THEN name || ' (current)'
ELSE name
END AS name
FROM users;

-- name: GetUserfromID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

