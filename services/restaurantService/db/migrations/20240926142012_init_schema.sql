-- +goose Up
-- +goose StatementBegin

CREATE TABLE Address (
    zip_code INT PRIMARY KEY,
    city VARCHAR(255) NOT NULL
);

CREATE TABLE Restaurant (
    ID serial PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Address TEXT NOT NULL,
    Rating DECIMAL(2, 1),
    category VARCHAR(50)
);

CREATE TABLE MenuItem (
    ID serial PRIMARY KEY,
    RestaurantID INT NOT NULL REFERENCES Restaurant (ID),
    Name VARCHAR(255) NOT NULL,
    Price DECIMAL(10, 2) NOT NULL,
    Description TEXT
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE Address;
DROP TABLE MenuItem;
DROP TABLE Restaurant;

-- +goose StatementEnd
