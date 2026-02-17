package postgres

import (
	"database/sql"
	"fmt"
	"os"
)

func GetConnectionFromEnv() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBDATABASE"), os.Getenv("DBSSLMODE"))
}

func NewDatabaseConnection(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
