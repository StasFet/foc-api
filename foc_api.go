package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var PORT string = os.Getenv("PORT")

func main() {
	db := InitDB("./database/test")
	var dbw = CreateDBWrapper(db)
	testPerformance, err := dbw.CreatePerformance(&Performance{GroupName: "Test"})
	if err != nil {
		log.Fatal(err)
	}

	retrievePerformance, err := dbw.GetPerformanceById(testPerformance.Id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Performance: %+v", retrievePerformance)


	return
	if PORT == "" {
		PORT = "3000"
	}

	fmt.Println("Hello From main.go!")
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the FOC REST API!")
	})

	fmt.Printf("Listening on port %s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, mux))
}