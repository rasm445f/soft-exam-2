-- +goose Up
-- +goose StatementBegin
CREATE TABLE todo (
    id serial PRIMARY KEY,
    title varchar(255) NOT NULL,
    text text NOT NULL,
    isCompleted boolean NOT NULL DEFAULT FALSE,
    category varchar(255),
    deadline date
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE todo
-- +goose StatementEnd