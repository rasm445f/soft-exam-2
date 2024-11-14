-- name: FetchAllRestaurants :many
SELECT id, name, address, rating
FROM restaurant;

-- name: GetRestaurantById :one
SELECT id, name, address, rating
FROM restaurant
WHERE id = $1;

-- name: GetMenuItemByRestaurantAndId :one
SELECT id, name, description, price, restaurantid
FROM menuitem
WHERE restaurantid = $1 AND id = $2;

-- name: FetchMenuItemsByRestaurantId :many
SELECT id, name, description, price, restaurantid
from menuitem
WHERE restaurantid = $1;

-- name: CreateRestaurant :one
INSERT INTO restaurant (name, address, rating)
VALUES ($1, $2, $3)
RETURNING id;

-- name: CreateMenuItem :one
INSERT INTO menuitem (name, description, price, restaurantid)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: FetchAllCategories :many
SELECT DISTINCT category
FROM restaurant
WHERE category IS NOT NULL
ORDER BY category;

-- name: FilterRestaurantsByCategory :many
SELECT id, name, address, rating, category
FROM restaurant
WHERE category ILIKE $1;

