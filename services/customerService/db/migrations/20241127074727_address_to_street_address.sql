-- +goose Up
-- +goose StatementBegin
ALTER TABLE Address RENAME COLUMN Address TO Street_Address;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Address RENAME COLUMN Street_Address TO Address;
-- +goose StatementEnd