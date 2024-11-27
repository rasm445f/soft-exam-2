-- Create a new Order
-- name: CreateOrder :one
INSERT INTO "Order" (TotalAmount, VATAmount, Status, Timestamp, Comment, CustomerID, RestaurantID, DeliveryAgentID, PaymentID, BonusID, FeeID)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING ID;

-- Fetch an Order by ID
-- name: GetOrderById :one
SELECT ID, TotalAmount, VATAmount, Status, Timestamp, Comment, CustomerID, RestaurantID, DeliveryAgentID, PaymentID, BonusID, FeeID
FROM "Order"
WHERE ID = $1;

-- Fetch all Orders
-- name: GetAllOrders :many
SELECT ID, TotalAmount, VATAmount, Status, Timestamp, Comment, CustomerID, RestaurantID, DeliveryAgentID, PaymentID, BonusID, FeeID
FROM "Order"
ORDER BY Timestamp DESC;

-- Update an Order's status
-- name: UpdateOrderStatus :exec
UPDATE "Order"
SET Status = $1
WHERE ID = $2;

-- Delete an Order
-- name: DeleteOrder :exec
DELETE FROM "Order"
WHERE ID = $1;

-- Create a new Order Item
-- name: CreateOrderItem :one
INSERT INTO OrderItem (OrderID, Name, Price, Quantity)
VALUES ($1, $2, $3, $4)
RETURNING ID;

-- Fetch all Order Items by Order ID
-- name: GetOrderItemsByOrderId :many
SELECT ID, OrderID, Name, Price, Quantity
FROM OrderItem
WHERE OrderID = $1;

-- Delete all Order Items by Order ID
-- name: DeleteOrderItemsByOrderId :exec
DELETE FROM OrderItem
WHERE OrderID = $1;

-- Create a Payment
-- name: CreatePayment :one
INSERT INTO Payment (PaymentStatus, PaymentMethod)
VALUES ($1, $2)
RETURNING ID;

-- Fetch a Payment by ID
-- name: GetPaymentById :one
SELECT ID, PaymentStatus, PaymentMethod
FROM Payment
WHERE ID = $1;

-- Create a Fee
-- name: CreateFee :one
INSERT INTO Fee (Percentage, Amount, Description)
VALUES ($1, $2, $3)
RETURNING ID;

-- Fetch a Fee by ID
-- name: GetFeeById :one
SELECT ID, Percentage, Amount, Description
FROM Fee
WHERE ID = $1;

-- Create a Bonus
-- name: CreateBonus :one
INSERT INTO Bonus (Description, EarlyLateAmount, Percentage)
VALUES ($1, $2, $3)
RETURNING ID;

-- Fetch a Bonus by ID
-- name: GetBonusById :one
SELECT ID, Description, EarlyLateAmount, Percentage
FROM Bonus
WHERE ID = $1;
