package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/rcovery/go-url-shortener/internal/http/handlers"
	infra_postgres "github.com/rcovery/go-url-shortener/internal/infra/postgres"
	"github.com/rcovery/go-url-shortener/shorturl"
	"github.com/rcovery/go-url-shortener/shorturl/postgres"
)

func main() {
	connectionString := infra_postgres.GetConnectionFromEnv()
	db, err := infra_postgres.NewDatabaseConnection(connectionString)
	if err != nil {
		panic(err)
	}

	repoInstance := postgres.NewRepository(db)
	serviceInstance := shorturl.NewService(repoInstance)
	log.Println("aoeiawoie")
	handlers.HandleShortURL(serviceInstance)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if host == "" || port == "" {
		log.Printf("Check your configuration for HOST %s - PORT %s", host, port)
	}

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	log.Printf("Starting server at %s:%s", host, port)
}
