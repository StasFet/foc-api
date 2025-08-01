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
	fmt.Println("foc_api.go started!")
	db := internal.InitDB("database/db.sqlite")
	fmt.Println("db initialised")
	defer db.Close()

	dbw := internal.CreateDBWrapper(db)

	createdPerformer, err := dbw.CreatePerformer(&internal.Performer{Name: "Edward Beaman", Email: "ebeam8@eq.edu.au"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("performer written into db!")

	p, err := dbw.GetPerformerById(createdPerformer.Id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", p)

	return
	if PORT == "" {
		PORT = "8000"
	}

	fmt.Println("Hello From foc_api.go!")
	mux := http.NewServeMux()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the FOC REST API!")
	})

	fmt.Printf("Listening on port %s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, mux))
}
