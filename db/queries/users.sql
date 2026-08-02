-- name: ListUsers :many
SELECT id, name, email, password, created_at
FROM users
ORDER BY id;