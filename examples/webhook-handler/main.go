package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	moyasar "github.com/yazeed1s/moyasar-go"
)

func main() {
	http.HandleFunc("/moyasar/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		var event moyasar.WebhookEvent
		if err := json.Unmarshal(body, &event); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		// Return 2xx fast. Do heavy work after this in your app.
		w.WriteHeader(http.StatusOK)

		fmt.Println("event:", event.ID, event.Type)
		fmt.Println("data:", string(event.Data))
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
