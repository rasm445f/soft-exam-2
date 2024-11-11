-- +goose Up
-- +goose StatementBegin
CREATE TABLE todo (
    id SERIAL8 PRIMARY KEY,
    title varchar(255) NOT NULL,
    text varchar(255) NOT NULL,
    iscompleted bool DEFAULT FALSE NOT NULL,
    category varchar(255),
    deadline timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE todo;

-- +goose StatementEnd
