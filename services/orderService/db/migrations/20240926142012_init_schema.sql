-- +goose Up
-- +goose StatementBegin
CREATE TABLE Bonus (
    ID serial PRIMARY KEY,
    Description text,
    EarlyLateAmount DECIMAL(10, 2),
    Percentage DECIMAL(5, 2)
);

CREATE TABLE Payment (
    ID serial PRIMARY KEY,
    PaymentStatus varchar(50) NOT NULL,
    PaymentMethod varchar(50) NOT NULL
);

CREATE TABLE Fee (
    ID serial PRIMARY KEY,
    Amount DECIMAL(10, 2) NOT NULL,
    Description text
);

CREATE TABLE "Order" (
    ID serial PRIMARY KEY,
    TotalAmount DECIMAL(10, 2) NOT NULL,
    VATAmount DECIMAL(10, 2),
    Status varchar(50) NOT NULL,
    Timestamp timestamp DEFAULT NOW(),
    Comment text,
    CustomerID int,
    RestaurantID int,
    DeliveryAgentID int,
    PaymentID int REFERENCES Payment (ID),
    BonusID int REFERENCES Bonus (ID),
    FeedbackID int
);

CREATE TABLE OrderItem (
    ID serial PRIMARY KEY,
    OrderID int NOT NULL REFERENCES "Order" (ID),
    Name varchar(255) NOT NULL,
    Price DECIMAL(10, 2) NOT NULL,
    Quantity DECIMAL(10, 2) NOT NULL
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE OrderItem;

DROP TABLE Fee;

DROP TABLE Payment;

DROP TABLE Bonus;

DROP TABLE "Order";

-- +goose StatementEnd
