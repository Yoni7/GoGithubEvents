package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func RunServer() {
	router := mux.NewRouter()

	router.HandleFunc("/github/event", GetEventTypes).Methods("GET")
	router.HandleFunc("/github/actors", GetActors).Methods("GET")
	router.HandleFunc("/github/repos", GetRepoUrls).Methods("GET")
	router.HandleFunc("/github/emails", GetEmails).Methods("GET")

	const port = 8080
	fmt.Printf("Starting server at port: %v\n", port)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
