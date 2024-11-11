-- +goose Up
-- +goose StatementBegin
INSERT INTO todo (id, title, text, iscompleted, category, deadline)
    VALUES (1, 'Task 1', 'This is the first task', FALSE, 'Work', '2023-12-31'),
    (2, 'Task 2', 'This is the second task', TRUE, 'Personal', '2023-11-30'),
    (3, 'Task 3', 'This is the third task', FALSE, 'Work', '2023-12-01');

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM todo
WHERE id IN (1, 2, 3);

-- +goose StatementEnd
