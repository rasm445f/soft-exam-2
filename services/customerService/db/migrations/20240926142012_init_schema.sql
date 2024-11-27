-- +goose Up
-- +goose StatementBegin

-- Create ZipCode table first
CREATE TABLE IF NOT EXISTS ZipCode (
    Zip_Code INT PRIMARY KEY,
    City varchar(255),
    CONSTRAINT zipcode_not_null CHECK (Zip_Code IS NOT NULL)
);

-- Create Address table second
CREATE TABLE IF NOT EXISTS Address (
    ID serial PRIMARY KEY,
    Address text,
    Zip_Code int REFERENCES ZipCode(Zip_Code),
    CONSTRAINT address_not_null CHECK (Address IS NOT NULL)
);

-- Create Customer table last
CREATE TABLE IF NOT EXISTS Customer (
    ID serial PRIMARY KEY,
    Name varchar(255),
    Email varchar(255) UNIQUE,
    PASSWORD varchar(255),
    PhoneNumber varchar(15),
    AddressID int REFERENCES Address(ID),
    CONSTRAINT email_not_null CHECK (Email IS NOT NULL),
    CONSTRAINT password_not_null CHECK (Password IS NOT NULL)
);

-- +goose StatementEnd
