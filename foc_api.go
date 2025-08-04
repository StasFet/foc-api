package main

import (
	"fmt"
	internal "foc_api/internal"
	"log"
	"net/http"
	"os"
)

var PORT string = os.Getenv("PORT")

func main() {
	db := internal.InitDB("database/db.sqlite")
	defer db.Close()

	wrapper := internal.CreateDBWrapper(db)
	api := internal.NewAPI(wrapper)

	if PORT == "" {
		PORT = "8000"
	}

	fmt.Println("Hello From foc_api.go!")
	mux := http.NewServeMux()

	mux.HandleFunc("/performances", api.PerformanceHandler)
	mux.HandleFunc("/performances/", api.PerformanceHandler)

	mux.HandleFunc("/performers", api.PerformerHandler)
	mux.HandleFunc("/performers/", api.PerformerHandler)

	mux.HandleFunc("/test", testRequest)
	mux.HandleFunc("/test/", testRequest)

	mux.HandleFunc("/junctions", api.JunctionHandler)
	mux.HandleFunc("/junctions/", api.JunctionHandler)

	fmt.Printf("Listening on port %s\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, mux))
}

func testRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the FOC REST API! request path: %s", r.URL.RawPath)
}
