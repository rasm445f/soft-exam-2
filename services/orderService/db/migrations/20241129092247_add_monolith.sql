-- +goose Up
-- +goose StatementBegin
CREATE TABLE Feedback (
    ID serial PRIMARY KEY,
    OrderID INT REFERENCES "Order" (ID) NOT NULL,
    CustomerID INT NOT NULL,
    DeliveryAgentRating INT,
    RestaurantRating INT,
    Comment TEXT
);

CREATE TABLE DeliveryAgent (
    ID serial PRIMARY KEY,
    FullName TEXT,
    ContactInfo TEXT,
    Availability BOOLEAN,
    Rating DECIMAL(10, 1)
);



-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE Feedback CASCADE
DROP TABLE DeliveryAgent CASCADE


-- +goose StatementEnd
