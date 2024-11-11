-- +goose Up
-- +goose StatementBegin
CREATE TABLE Customer (
    ID serial PRIMARY KEY,
    Name varchar(255) NOT NULL,
    Email varchar(255) UNIQUE NOT NULL,
    PhoneNumber varchar(15),
    Address text
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE Customer;

-- DROP TABLE Feedback;
--
-- DROP TABLE Complaint;
-- +goose StatementEnd
