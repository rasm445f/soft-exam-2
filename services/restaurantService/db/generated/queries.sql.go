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
INSERT INTO menuitem (name, description, price, restaurantid)
VALUES ($1, $2, $3, $4)
RETURNING id
`

type CreateMenuItemParams struct {
	Name         string         `json:"name"`
	Description  *string        `json:"description"`
	Price        pgtype.Numeric `json:"price"`
	Restaurantid int32          `json:"restaurantid"`
}

func (q *Queries) CreateMenuItem(ctx context.Context, arg CreateMenuItemParams) (int32, error) {
	row := q.db.QueryRow(ctx, createMenuItem,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.Restaurantid,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const createRestaurant = `-- name: CreateRestaurant :one
INSERT INTO restaurant (name, address, rating)
VALUES ($1, $2, $3)
RETURNING id
`

type CreateRestaurantParams struct {
	Name    string         `json:"name"`
	Address string         `json:"address"`
	Rating  pgtype.Numeric `json:"rating"`
}

func (q *Queries) CreateRestaurant(ctx context.Context, arg CreateRestaurantParams) (int32, error) {
	row := q.db.QueryRow(ctx, createRestaurant, arg.Name, arg.Address, arg.Rating)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const fetchAllRestaurants = `-- name: FetchAllRestaurants :many
SELECT id, name, address, rating
FROM restaurant
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
			&i.Address,
			&i.Rating,
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
SELECT id, name, description, price, restaurantid
from menuitem
WHERE restaurantid = $1
`

type FetchMenuItemsByRestaurantIdRow struct {
	ID           int32          `json:"id"`
	Name         string         `json:"name"`
	Description  *string        `json:"description"`
	Price        pgtype.Numeric `json:"price"`
	Restaurantid int32          `json:"restaurantid"`
}

func (q *Queries) FetchMenuItemsByRestaurantId(ctx context.Context, restaurantid int32) ([]FetchMenuItemsByRestaurantIdRow, error) {
	rows, err := q.db.Query(ctx, fetchMenuItemsByRestaurantId, restaurantid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchMenuItemsByRestaurantIdRow
	for rows.Next() {
		var i FetchMenuItemsByRestaurantIdRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Restaurantid,
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
SELECT id, name, description, price, restaurantid
FROM menuitem
WHERE restaurantid = $1 AND id = $2
`

type GetMenuItemByRestaurantAndIdParams struct {
	Restaurantid int32 `json:"restaurantid"`
	ID           int32 `json:"id"`
}

type GetMenuItemByRestaurantAndIdRow struct {
	ID           int32          `json:"id"`
	Name         string         `json:"name"`
	Description  *string        `json:"description"`
	Price        pgtype.Numeric `json:"price"`
	Restaurantid int32          `json:"restaurantid"`
}

func (q *Queries) GetMenuItemByRestaurantAndId(ctx context.Context, arg GetMenuItemByRestaurantAndIdParams) (GetMenuItemByRestaurantAndIdRow, error) {
	row := q.db.QueryRow(ctx, getMenuItemByRestaurantAndId, arg.Restaurantid, arg.ID)
	var i GetMenuItemByRestaurantAndIdRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.Restaurantid,
	)
	return i, err
}

const getRestaurantById = `-- name: GetRestaurantById :one
SELECT id, name, address, rating
FROM restaurant
WHERE id = $1
`

func (q *Queries) GetRestaurantById(ctx context.Context, id int32) (Restaurant, error) {
	row := q.db.QueryRow(ctx, getRestaurantById, id)
	var i Restaurant
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Rating,
	)
	return i, err
}
