-- name: CreateCustomer :one
INSERT INTO customer (name, email, phonenumber, address, password)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetCustomerByID :one
SELECT * FROM customer WHERE id = $1;

-- name: GetAllCustomers :many
SELECT * FROM customer ORDER BY name;

-- name: UpdateCustomer :one
UPDATE customer
SET name = $2, email = $3, phonenumber = $4, address = $5, password = $6
WHERE id = $1
RETURNING *;

-- name: DeleteCustomer :exec
DELETE FROM customer WHERE id = $1;