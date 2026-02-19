package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rcovery/go-url-shortener/internal/config"
	"github.com/rcovery/go-url-shortener/internal/http/handlers"
	infra_postgres "github.com/rcovery/go-url-shortener/internal/infra/postgres"
	"github.com/rcovery/go-url-shortener/shorturl"
	"github.com/rcovery/go-url-shortener/shorturl/postgres"
)

func main() {
	config.InitConfig()

	connectionString := infra_postgres.GetConnectionFromEnv()
	db, databaseErr := infra_postgres.NewDatabaseConnection(connectionString)
	if databaseErr != nil {
		panic(databaseErr)
	}

	repoInstance := postgres.NewRepository(db)
	serviceInstance := shorturl.NewService(repoInstance)
	handlers.HandleShortURL(serviceInstance)

	host := config.GetString("HOST")
	port := config.GetString("PORT")

	log.Printf("%v:%v\n", host, port)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
