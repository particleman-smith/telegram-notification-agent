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

	//
	/// Set routes
	//

	router.HandleFunc("/test", api.Test).Methods("POST")

	// Bash
	router.HandleFunc("/bash-event/exec-failure", api.Error).Methods("POST")

	// ZFS
	router.HandleFunc("/zfs-event/data-error", api.Error).Methods("POST")
	router.HandleFunc("/zfs-event/zpool-state", api.Error).Methods("POST")

	// Backup
	router.HandleFunc("/backup-event/failure", api.Test).Methods("POST")

	// SMART

	// Launch server
	err := http.ListenAndServe(":9090", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
