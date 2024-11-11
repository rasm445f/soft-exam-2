-- +goose Up
-- +goose StatementBegin
CREATE TABLE Restaurant (
    ID serial PRIMARY KEY,
    Name varchar(255) NOT NULL,
    Address text NOT NULL,
    Rating DECIMAL(2, 1)
);

CREATE TABLE MenuItem (
    ID serial PRIMARY KEY,
    RestaurantID int NOT NULL REFERENCES Restaurant (ID),
    Name varchar(255) NOT NULL,
    Price DECIMAL(10, 2) NOT NULL,
    Description text
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE Restaurant;

DROP TABLE MenuItem;

-- +goose StatementEnd
