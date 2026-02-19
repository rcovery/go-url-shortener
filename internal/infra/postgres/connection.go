package postgres

import (
	"database/sql"
	"fmt"

	"github.com/rcovery/go-url-shortener/internal/config"
)

func GetConnectionFromEnv() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s",
		config.GetString("DBHOST"),
		config.GetString("DBUSER"),
		config.GetString("DBPASS"),
		config.GetString("DBDATABASE"),
		config.GetString("DBSSLMODE"),
	)
}

func NewDatabaseConnection(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
