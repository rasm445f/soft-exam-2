-- name: FetchAllRestaurants :many
SELECT r.id, r.name, r.rating, r.category, r.street, r.zip_code
FROM restaurant r
JOIN address a ON r.zip_code = a.zip_code;

-- name: GetRestaurantById :one
SELECT id, name, rating, category, street, zip_code
FROM restaurant
WHERE id = $1;

-- name: GetMenuItemByRestaurantAndId :one
SELECT id, restaurantid, name, price, description
FROM menuitem
WHERE restaurantid = $1 AND id = $2;

-- name: FetchMenuItemsByRestaurantId :many
SELECT id, restaurantid, name, price, description
from menuitem
WHERE restaurantid = $1;

-- name: CreateAddress :exec
INSERT INTO address (zip_code, city)
VALUES ($1, $2)
ON CONFLICT (zip_code) DO NOTHING;

-- name: CreateRestaurant :one
INSERT INTO restaurant (name, rating, category, street, zip_code)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: CreateMenuItem :one
INSERT INTO menuitem (restaurantid, name, price, description)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: FetchAllCategories :many
SELECT DISTINCT category
FROM restaurant
WHERE category IS NOT NULL
ORDER BY category;

-- name: FilterRestaurantsByCategory :many
SELECT id, name, rating, category, street, zip_code
FROM restaurant
WHERE category ILIKE $1;

