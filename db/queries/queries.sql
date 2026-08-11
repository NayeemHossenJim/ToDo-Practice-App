-- ============================================
-- USERS
-- ============================================

-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    password
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET
    name = $2,
    email = $3,
    password = $4
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;


-- ============================================
-- TODOS
-- ============================================

-- name: CreateTodo :one
INSERT INTO todos (
    title,
    description,
    user_id,
    due_date
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetTodoByID :one
SELECT *
FROM todos
WHERE id = $1
LIMIT 1;

-- name: GetTodoByIDAndUserID :one
SELECT *
FROM todos
WHERE id = $1
  AND user_id = $2
LIMIT 1;

-- name: ListTodosByUser :many
SELECT *
FROM todos
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListCompletedTodosByUser :many
SELECT *
FROM todos
WHERE user_id = $1
  AND completed = true
ORDER BY created_at DESC;

-- name: ListPendingTodosByUser :many
SELECT *
FROM todos
WHERE user_id = $1
  AND completed = false
ORDER BY created_at DESC;

-- name: ListTodosByDueDate :many
SELECT *
FROM todos
WHERE user_id = $1
ORDER BY due_date ASC NULLS LAST;

-- name: UpdateTodo :one
UPDATE todos
SET
    title = $2,
    description = $3,
    completed = $4,
    due_date = $5
WHERE id = $1
RETURNING *;

-- name: UpdateTodoByUser :one
UPDATE todos
SET
    title = $3,
    description = $4,
    completed = $5,
    due_date = $6
WHERE id = $1
  AND user_id = $2
RETURNING *;

-- name: MarkTodoCompleted :one
UPDATE todos
SET completed = true
WHERE id = $1
  AND user_id = $2
RETURNING *;

-- name: MarkTodoPending :one
UPDATE todos
SET completed = false
WHERE id = $1
  AND user_id = $2
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;

-- name: DeleteTodoByUser :exec
DELETE FROM todos
WHERE id = $1
  AND user_id = $2;

-- name: DeleteAllTodosByUser :exec
DELETE FROM todos
WHERE user_id = $1;

-- name: CountTodosByUser :one
SELECT COUNT(*)
FROM todos
WHERE user_id = $1;

-- name: CountPendingTodosByUser :one
SELECT COUNT(*)
FROM todos
WHERE user_id = $1
  AND completed = false;

-- name: CountCompletedTodosByUser :one
SELECT COUNT(*)
FROM todos
WHERE user_id = $1
  AND completed = true;