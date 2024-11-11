-- +goose Up
-- +goose StatementBegin
ALTER TABLE Customer
    ADD COLUMN PASSWORD varchar(255) NOT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
