-- +goose Up
-- +goose StatementBegin
CREATE TABLE Customer (
    ID serial PRIMARY KEY,
    Name varchar(255),
    Email varchar(255) UNIQUE,
    PASSWORD varchar(255),
    PhoneNumber varchar(15),
    Address text,
    CONSTRAINT email_not_null CHECK (Email IS NOT NULL), -- Ensures non-null for essential fields
    CONSTRAINT password_not_null CHECK (Password IS NOT NULL) -- Validates non-null password
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE Customer;

-- +goose StatementEnd
