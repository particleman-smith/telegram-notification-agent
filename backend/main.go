/*
Backend of the ScanMan web app. RESTful API only.
*/

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/particleman-smith/telegram-notification-agent/backend/api"
	"github.com/rs/cors"
)

func main() {
	// Using gorilla/mux for passing params in URL
	router := mux.NewRouter()

	// Using rs/cors
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:9090"},
	})
	handler := c.Handler(router)

	// Set routes
	router.HandleFunc("/zfs-event/test", api.Test).Methods("POST")

	// Launch server
	err := http.ListenAndServe(":8000", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
