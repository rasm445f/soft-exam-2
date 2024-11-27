-- name: CreateCustomer :exec
WITH new_address AS (
    INSERT INTO address (street_address, zip_code)
    VALUES ($1, $2)
    RETURNING id AS address_id
)
INSERT INTO customer (name, email, phonenumber, addressid, password)
VALUES ($3, $4, $5, (SELECT address_id FROM new_address), $6);


-- name: GetCustomerByID :one
SELECT 
    c.id,
    c.name,
    c.email,
    c.phonenumber,
    c.password,
    a.street_address,
    z.zip_code,
    z.city
FROM customer c
LEFT JOIN address a ON c.addressid = a.id
LEFT JOIN zipcode z ON a.zip_code = z.zip_code
WHERE c.id = $1;

-- name: GetAllCustomers :many
SELECT 
    c.id,
    c.name,
    c.email,
    c.phonenumber,
    c.password,
    a.street_address AS street_address,
    z.zip_code,
    z.city
FROM customer c
LEFT JOIN address a ON c.addressid = a.id
LEFT JOIN zipcode z ON a.zip_code = z.zip_code
ORDER BY c.name;


-- name: UpdateCustomer :exec
UPDATE customer
SET 
    name = COALESCE($2, name),
    email = COALESCE($3, email),
    phonenumber = COALESCE($4, phonenumber),
    password = COALESCE($5, password)
WHERE id = $1;


-- name: UpdateAddress :exec
UPDATE address
SET 
    street_address = COALESCE($2, street_address),
    zip_code = COALESCE($3, zip_code)
WHERE id = $1;



-- name: CreateZipCode :exec
INSERT INTO zipcode (zip_code, city)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteCustomer :exec
DELETE FROM customer WHERE id = $1;