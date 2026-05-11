package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"
	"github.com/rcovery/go-url-shortener/internal/config"
	"github.com/rcovery/go-url-shortener/internal/http/handlers"
	infra_postgres "github.com/rcovery/go-url-shortener/internal/infra/postgres"
	"github.com/rcovery/go-url-shortener/shorturl"
	"github.com/rcovery/go-url-shortener/shorturl/postgres"
)

func main() {
	config.InitConfig()

	baseContext, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	connectionString := infra_postgres.GetConnectionFromEnv()
	db, databaseErr := infra_postgres.NewDatabaseConnection(connectionString)
	if databaseErr != nil {
		panic(databaseErr)
	}

	otelShutdown, err := config.SetupOTelSDK(baseContext)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	repoInstance := postgres.NewRepository(db)
	serviceInstance := shorturl.NewService(repoInstance)
	handlers.HandleShortURL(serviceInstance)

	host := config.GetString("HOST")
	port := config.GetString("PORT")

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		BaseContext:  func(net.Listener) context.Context { return baseContext },
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- server.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		log.Fatal(err)
	case <-baseContext.Done():
		stop()
	}

	err = server.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
