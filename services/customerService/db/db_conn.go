package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() (*pgx.Conn, error) {

	connStr := os.Getenv("DBSTRING")
	ctx := context.Background()

	db, err := pgx.Connect(ctx, connStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to the database!")

	return db, nil
}
