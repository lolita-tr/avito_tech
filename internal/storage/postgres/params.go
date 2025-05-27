package postgres

import (
	"fmt"
	"os"
)

type DBParams struct {
	URL string
}

func NewDBParams() *DBParams {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DB")

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, database)

	return &DBParams{url}
}
