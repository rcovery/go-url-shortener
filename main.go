package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rcovery/go-url-shortener/internal/infra/postgres"
)

func main() {
	connectionString := postgres.GetConnectionFromEnv()
	db, err := postgres.NewDatabaseConnection(connectionString)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/api/url", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			{
				rawBody := r.Body
				body, err := io.ReadAll(rawBody)
				if err != nil {
					fmt.Println("Erro ao receber body!")
					return
				}
				if len(body) == 0 {
					fmt.Println("Body está vazio!")
					return
				}

				var createURLBody CreateURL
				err = json.Unmarshal(body, &createURLBody)
				if err != nil {
					fmt.Println("Não foi possível decodificar o JSON!")
					return
				}

				_, URLErr := createNewURL(db, createURLBody)
				if URLErr != nil {
					log.Println(URLErr)
					w.WriteHeader(400)
					break
				}

				w.WriteHeader(200)
				break
			}
		default:
			{
				w.WriteHeader(405)
			}
		}
	})
	http.HandleFunc("/{url_name}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			{
				urlName := r.PathValue("url_name")
				urlFromDatabase, _ := getURLByName(db, urlName)
				if urlFromDatabase == "" {
					log.Println(err)
					w.WriteHeader(404)
					break
				}

				w.Header().Add("Location", urlFromDatabase)
				w.WriteHeader(303)
				break
			}
		default:
			{
				w.WriteHeader(405)
			}
		}
	})

	if err := http.ListenAndServe("0.0.0.0:9000", nil); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func getURLByName(db *sql.DB, name string) (string, string) {
	result, err := db.Query(`
		SELECT id, url, idempotency_key
		FROM urls
		WHERE name = $1
			AND expires_at > NOW()
		LIMIT 1
		`, name)
	if err != nil {
		log.Println(err, "Get URL")
	}
	defer result.Close()
	result.Next()

	var url, id, idempotency string
	if err = result.Scan(&id, &url, &idempotency); err != nil {
		log.Println(err, "Get URL - Scan error")
		return "", ""
	}

	return url, idempotency
}

func createNewURL(db *sql.DB, URLCreationData CreateURL) (string, error) {
	alreadyExistsURL, idempotency := getURLByName(db, URLCreationData.Name)

	if alreadyExistsURL != "" {
		if idempotency == URLCreationData.IdempotencyKey {
			return URLCreationData.URL, nil
		}

		return "", fmt.Errorf("URL name already been taken")
	}

	urlUUID, urlUUIDErr := uuid.NewV7()
	if urlUUIDErr != nil {
		return "", urlUUIDErr
	}

	_, insertionErr := db.Query(`
		INSERT INTO urls
			(id, url, name, idempotency_key)
		VALUES
			($1, $2, $3, $4)
		`, urlUUID, URLCreationData.URL, URLCreationData.Name, URLCreationData.IdempotencyKey)
	if insertionErr != nil {
		log.Println(insertionErr, "Create URL - Insert error")
		return "", insertionErr
	}

	return URLCreationData.URL, nil
}

type CreateURL struct {
	URL            string `json:"url"`
	Name           string `json:"name"`
	IdempotencyKey string `json:"idempotency_key"`
}
