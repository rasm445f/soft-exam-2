-- name: CreateCustomer :one
INSERT INTO customer (name, email, phonenumber, address, password)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, email, phonenumber, address, password;

-- name: GetCustomerByID :one
SELECT * FROM customer WHERE id = $1;

-- name: GetAllCustomers :many
SELECT * FROM customer ORDER BY name;

-- name: UpdateCustomer :exec
UPDATE customer
SET 
    name = COALESCE($2, name),
    email = COALESCE($3, email),
    phonenumber = COALESCE($4, phonenumber),
    address = COALESCE($5, address),
    password = COALESCE($6, password)
WHERE id = $1;

-- name: DeleteCustomer :exec
DELETE FROM customer WHERE id = $1;
