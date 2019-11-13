/*
Backend of the ScanMan web app. RESTful API only.
*/

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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

	router.HandleFunc("/test", Test).Methods("POST")

	// Bash
	router.HandleFunc("/bash-event/exec-failure", Error).Methods("POST")

	// ZFS
	router.HandleFunc("/zfs-event/data-error", Error).Methods("POST")
	router.HandleFunc("/zfs-event/zpool-state", Error).Methods("POST")

	// Backup
	router.HandleFunc("/backup-event/failure", Test).Methods("POST")

	// SMART

	println("Launching server on port 9090")

	// Launch server
	err := http.ListenAndServe(":9090", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
