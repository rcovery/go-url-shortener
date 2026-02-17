package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rcovery/go-url-shortener/shorturl"
)

func HandleShortURL(service *shorturl.Service) {
	http.HandleFunc("/api/url", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			{
				ctx, ctxCancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer ctxCancel()

				contentType := r.Header.Get("Content-Type")
				if contentType != "application/json" {
					w.WriteHeader(400)
					return
				}

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

				var createURLBody shorturl.ShortURL
				err = json.Unmarshal(body, &createURLBody)
				if err != nil {
					fmt.Println("Não foi possível decodificar o JSON!")
					return
				}

				createdLink, URLErr := service.Create(ctx, createURLBody.ID, createURLBody.IdempotencyKey, createURLBody.Name, createURLBody.Link)
				if URLErr != nil {
					log.Println(URLErr)
					w.WriteHeader(400)
					break
				}
				if createdLink == "" {
					log.Println("Created an empty URL")
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
				ctx, ctxCancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer ctxCancel()

				time.Sleep(10 * time.Second)

				urlName := r.PathValue("url_name")
				urlFromDatabase, selectionError := service.Select(ctx, urlName)
				if selectionError != nil || urlFromDatabase == "" {
					log.Println(selectionError)
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
}
