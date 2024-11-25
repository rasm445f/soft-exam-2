// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package generated

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createMenuItem = `-- name: CreateMenuItem :one
INSERT INTO menuitem (restaurantid, name, price, description)
VALUES ($1, $2, $3, $4)
RETURNING id
`

type CreateMenuItemParams struct {
	Restaurantid int32          `json:"restaurantid"`
	Name         string         `json:"name"`
	Price        pgtype.Numeric `json:"price"`
	Description  *string        `json:"description"`
}

func (q *Queries) CreateMenuItem(ctx context.Context, arg CreateMenuItemParams) (int32, error) {
	row := q.db.QueryRow(ctx, createMenuItem,
		arg.Restaurantid,
		arg.Name,
		arg.Price,
		arg.Description,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const createRestaurant = `-- name: CreateRestaurant :one
INSERT INTO restaurant (name, rating, category, address, zip_code)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

type CreateRestaurantParams struct {
	Name     string         `json:"name"`
	Rating   pgtype.Numeric `json:"rating"`
	Category *string        `json:"category"`
	Address  *string        `json:"address"`
	ZipCode  *int32         `json:"zip_code"`
}

func (q *Queries) CreateRestaurant(ctx context.Context, arg CreateRestaurantParams) (int32, error) {
	row := q.db.QueryRow(ctx, createRestaurant,
		arg.Name,
		arg.Rating,
		arg.Category,
		arg.Address,
		arg.ZipCode,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const createZipCode = `-- name: CreateZipCode :exec
INSERT INTO zipcode (zip_code, city)
VALUES ($1, $2)
ON CONFLICT (zip_code) DO NOTHING
`

type CreateZipCodeParams struct {
	ZipCode int32  `json:"zip_code"`
	City    string `json:"city"`
}

func (q *Queries) CreateZipCode(ctx context.Context, arg CreateZipCodeParams) error {
	_, err := q.db.Exec(ctx, createZipCode, arg.ZipCode, arg.City)
	return err
}

const fetchAllCategories = `-- name: FetchAllCategories :many
SELECT DISTINCT category
FROM restaurant
WHERE category IS NOT NULL
ORDER BY category
`

func (q *Queries) FetchAllCategories(ctx context.Context) ([]*string, error) {
	rows, err := q.db.Query(ctx, fetchAllCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*string
	for rows.Next() {
		var category *string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		items = append(items, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchAllRestaurants = `-- name: FetchAllRestaurants :many
SELECT r.id, r.name, r.rating, r.category, r.address, r.zip_code
FROM restaurant r
JOIN zipcode a ON r.zip_code = a.zip_code
`

func (q *Queries) FetchAllRestaurants(ctx context.Context) ([]Restaurant, error) {
	rows, err := q.db.Query(ctx, fetchAllRestaurants)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Restaurant
	for rows.Next() {
		var i Restaurant
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Rating,
			&i.Category,
			&i.Address,
			&i.ZipCode,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchMenuItemsByRestaurantId = `-- name: FetchMenuItemsByRestaurantId :many
SELECT id, restaurantid, name, price, description
from menuitem
WHERE restaurantid = $1
`

func (q *Queries) FetchMenuItemsByRestaurantId(ctx context.Context, restaurantid int32) ([]Menuitem, error) {
	rows, err := q.db.Query(ctx, fetchMenuItemsByRestaurantId, restaurantid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Menuitem
	for rows.Next() {
		var i Menuitem
		if err := rows.Scan(
			&i.ID,
			&i.Restaurantid,
			&i.Name,
			&i.Price,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const filterRestaurantsByCategory = `-- name: FilterRestaurantsByCategory :many
SELECT id, name, rating, category, address, zip_code
FROM restaurant
WHERE category ILIKE $1
`

func (q *Queries) FilterRestaurantsByCategory(ctx context.Context, category *string) ([]Restaurant, error) {
	rows, err := q.db.Query(ctx, filterRestaurantsByCategory, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Restaurant
	for rows.Next() {
		var i Restaurant
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Rating,
			&i.Category,
			&i.Address,
			&i.ZipCode,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMenuItemByRestaurantAndId = `-- name: GetMenuItemByRestaurantAndId :one
SELECT id, restaurantid, name, price, description
FROM menuitem
WHERE restaurantid = $1 AND id = $2
`

type GetMenuItemByRestaurantAndIdParams struct {
	Restaurantid int32 `json:"restaurantid"`
	ID           int32 `json:"id"`
}

func (q *Queries) GetMenuItemByRestaurantAndId(ctx context.Context, arg GetMenuItemByRestaurantAndIdParams) (Menuitem, error) {
	row := q.db.QueryRow(ctx, getMenuItemByRestaurantAndId, arg.Restaurantid, arg.ID)
	var i Menuitem
	err := row.Scan(
		&i.ID,
		&i.Restaurantid,
		&i.Name,
		&i.Price,
		&i.Description,
	)
	return i, err
}

const getRestaurantById = `-- name: GetRestaurantById :one
SELECT id, name, rating, category, address, zip_code
FROM restaurant
WHERE id = $1
`

func (q *Queries) GetRestaurantById(ctx context.Context, id int32) (Restaurant, error) {
	row := q.db.QueryRow(ctx, getRestaurantById, id)
	var i Restaurant
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Rating,
		&i.Category,
		&i.Address,
		&i.ZipCode,
	)
	return i, err
}
