-- name: FetchAllTodos :many
SELECT
    id,
    title,
    text,
    iscompleted,
    category,
    deadline
FROM
    todo;

-- name: GetTodoById :one
SELECT
    id,
    title,
    text,
    iscompleted,
    category,
    deadline
FROM
    todo
WHERE
    id = $1;

-- name: DeleteTodoById :exec
DELETE FROM todo
WHERE id = $1;

-- name: CreateTodo :one
INSERT INTO todo (title, text, iscompleted, category, deadline)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    id;

