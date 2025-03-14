package main

import (
	"fmt"
	"go_microservice/internal/db"
	"go_microservice/internal/handlers"
	"log"
	"net/http"
)

func main() {
	db.InitDB()

	http.HandleFunc("/", handlers.Handler)
	http.HandleFunc("/api/status", handlers.StatusHandler)
	http.HandleFunc("/api/info", handlers.InfoHandler)
	http.HandleFunc("/api/geo", handlers.GeoHandler)
	http.HandleFunc("/api/save", handlers.FetchAndSaveData)

	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
