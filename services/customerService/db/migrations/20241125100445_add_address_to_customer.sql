-- +goose Up
-- +goose StatementBegin
CREATE TABLE Customer (
    ID serial PRIMARY KEY,
    Name varchar(255),
    Email varchar(255) UNIQUE,
    PASSWORD varchar(255),
    PhoneNumber varchar(15),
    AddressID int REFERENCES Address(ID),
    CONSTRAINT email_not_null CHECK (Email IS NOT NULL),
    CONSTRAINT password_not_null CHECK (Password IS NOT NULL)
);

CREATE TABLE Address (
    ID serial PRIMARY KEY,
    Address text,
    Zip_Code int REFERENCES ZipCode(Zip_Code),
    CONSTRAINT address_not_null CHECK (Address IS NOT NULL)
);

CREATE TABLE ZipCode (
    Zip_Code INT PRIMARY KEY,
    City varchar(255),
    CONSTRAINT zipcode_not_null CHECK (ZipCode IS NOT NULL)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE Customer;

-- +goose StatementEnd
