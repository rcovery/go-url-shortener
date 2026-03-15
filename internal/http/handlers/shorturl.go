package handlers

import (
	"context"
	"encoding/json"
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
					writeJSONError(w, http.StatusBadRequest, "invalid_content_type")
					return
				}

				rawBody := http.MaxBytesReader(w, r.Body, 1*MB)
				body, err := io.ReadAll(rawBody)
				if err != nil {
					log.Println("failed reading body:", err)
					writeJSONError(w, http.StatusBadRequest, "invalid_body")
					return
				}
				if len(body) == 0 {
					writeJSONError(w, http.StatusBadRequest, "empty_body")
					return
				}

				var createURLBody shorturl.ShortURL
				err = json.Unmarshal(body, &createURLBody)
				if err != nil {
					log.Println("failed decoding json:", err)
					writeJSONError(w, http.StatusBadRequest, "invalid_json")
					return
				}

				createdLink, URLErr := service.Create(ctx, createURLBody.ID, createURLBody.IdempotencyKey, createURLBody.Name, createURLBody.Link)
				if URLErr != nil {
					log.Println(URLErr)
					writeJSONError(w, http.StatusBadRequest, "create_failed")
					break
				}
				if createdLink == "" {
					log.Println("Created an empty URL")
					writeJSONError(w, http.StatusBadRequest, "create_failed")
					break
				}

				w.WriteHeader(200)
				break
			}
		default:
			{
				writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed")
			}
		}
	})

	http.HandleFunc("/{url_name}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			{
				ctx, ctxCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer ctxCancel()

				urlName := r.PathValue("url_name")
				urlFromDatabase, selectionError := service.Select(ctx, urlName)
				if selectionError != nil || urlFromDatabase == "" {
					log.Println(selectionError)
					writeJSONError(w, http.StatusNotFound, "not_found")
					break
				}

				w.Header().Add("Location", urlFromDatabase)
				w.WriteHeader(303)
				break
			}
		default:
			{
				writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed")
			}
		}
	})
}

func writeJSONError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encodeErr := json.NewEncoder(w).Encode(map[string]string{
		"error": code,
	})
	if encodeErr != nil {
		log.Println("failed encoding json error:", encodeErr)
	}
}

const (
	B  int64 = 1
	KB       = 1024 * B
	MB       = 1024 * KB
)
