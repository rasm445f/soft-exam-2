-- Create a new Order
-- name: CreateOrder :one
INSERT INTO "Order" (TotalAmount, VATAmount, Status, Timestamp, Comment, CustomerID, RestaurantID, DeliveryAgentID, PaymentID, BonusID, FeeID)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING
    ID;

-- Fetch an Order by ID
-- name: GetOrderById :one
SELECT
    ID,
    TotalAmount,
    VATAmount,
    Status,
    Timestamp,
    Comment,
    CustomerID,
    RestaurantID,
    DeliveryAgentID,
    PaymentID,
    BonusID,
    FeeID
FROM
    "Order"
WHERE
    ID = $1;

-- Fetch all Orders
-- name: GetAllOrders :many
SELECT
    ID,
    TotalAmount,
    VATAmount,
    Status,
    Timestamp,
    Comment,
    CustomerID,
    RestaurantID,
    DeliveryAgentID,
    PaymentID,
    BonusID,
    FeeID
FROM
    "Order"
ORDER BY
    Timestamp DESC;

-- Update an Order's status
-- name: UpdateOrderStatus :exec
UPDATE
    "Order"
SET
    Status = $1
WHERE
    ID = $2;

-- Update an Order's status and deliveryAgent
-- name: UpdateOrderStatusAndDeliveryAgent :exec
UPDATE
    "Order"
SET
    Status = $1,
    DeliveryAgentID = $2
WHERE
    ID = $3;

-- Delete an Order
-- name: DeleteOrder :exec
DELETE FROM "Order"
WHERE ID = $1;

-- Create a new Order Item
-- name: CreateOrderItem :one
INSERT INTO OrderItem (OrderID, Name, Price, Quantity)
    VALUES ($1, $2, $3, $4)
RETURNING
    ID;

-- Fetch all Order Items by Order ID
-- name: GetOrderItemsByOrderId :many
SELECT
    ID,
    OrderID,
    Name,
    Price,
    Quantity
FROM
    OrderItem
WHERE
    OrderID = $1;

-- Delete all Order Items by Order ID
-- name: DeleteOrderItemsByOrderId :exec
DELETE FROM OrderItem
WHERE OrderID = $1;

-- Create a Payment
-- name: CreatePayment :one
INSERT INTO Payment (PaymentStatus, PaymentMethod)
    VALUES ($1, $2)
RETURNING
    ID;

-- Fetch a Payment by ID
-- name: GetPaymentById :one
SELECT
    ID,
    PaymentStatus,
    PaymentMethod
FROM
    Payment
WHERE
    ID = $1;

-- Create a Fee
-- name: CreateFee :one
INSERT INTO Fee (Percentage, Amount, Description)
    VALUES ($1, $2, $3)
RETURNING
    ID;

-- Fetch a Fee by ID
-- name: GetFeeById :one
SELECT
    ID,
    Percentage,
    Amount,
    Description
FROM
    Fee
WHERE
    ID = $1;

-- Create a Bonus
-- name: CreateBonus :one
INSERT INTO Bonus (Description, EarlyLateAmount, Percentage)
    VALUES ($1, $2, $3)
RETURNING
    ID;

-- Fetch a Bonus by ID
-- name: GetBonusById :one
SELECT
    ID,
    Description,
    EarlyLateAmount,
    Percentage
FROM
    Bonus
WHERE
    ID = $1;

-- Fetch all Feedbacks
-- name: GetAllFeedbacks :many
SELECT
    *
FROM
    Feedback
ORDER BY
    ID DESC;

-- Create Feedback
-- name: CreateFeedback :one
INSERT INTO Feedback (OrderID, CustomerID, DeliveryAgentRating, RestaurantRating, Comment)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    ID;

-- Fetch a Feedback by ID
-- name: GetFeedbackById :one
SELECT
    *
FROM
    Feedback
WHERE
    ID = $1;

-- Fetch a Feedback by OrderID
-- name: GetFeedbackByOrderId :one
SELECT
    *
FROM
    Feedback
WHERE
    OrderID = $1;

-- Fetch all DeliveryAgents
-- name: GetAllDeliveryAgents :many
SELECT
    *
FROM
    DeliveryAgent
ORDER BY
    ID DESC;

-- Create DeliveryAgent
-- name: CreateDeliveryAgent :one
INSERT INTO DeliveryAgent (FullName, ContactInfo, Availability, Rating)
    VALUES ($1, $2, $3, $4)
RETURNING
    ID;

-- Fetch an Order by ID
-- name: GetDeliveryAgentById :one
SELECT
    *
FROM
    DeliveryAgent
WHERE
    ID = $1;

-- Update Delivery Agent Availability
-- name: UpdateDeliveryAgentAvailability :exec
UPDATE
    DeliveryAgent
SET
    Availability = $1
WHERE
    ID = $2;

-- Update Delivery Agent Rating
-- name: UpdateDeliveryAgentRating :exec
UPDATE
    DeliveryAgent
SET
    Rating = $1
WHERE
    ID = $2;

-- Get all Feedbacks for a Delivery Agent
-- name: GetAllFeedbacksFromDeliveryAgentByOrderId :many
SELECT
    f.ID,
    f.OrderID,
    f.CustomerID,
    f.DeliveryAgentRating,
    f.RestaurantRating,
    f.Comment,
    o.DeliveryAgentID
FROM
    Feedback f
    JOIN "Order" o ON f.OrderID = o.ID
WHERE
    o.DeliveryAgentID = (
        SELECT
            DeliveryAgentID
        FROM
            "Order"
        WHERE
            OrderID = $1
        LIMIT 1);

-- Update an Order's bonus
-- name: UpdateOrderBonus :exec
UPDATE
    "Order"
SET
    BonusID = $1
WHERE
    ID = $2;


