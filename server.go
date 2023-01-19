package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

func RunServer() {
	router := mux.NewRouter()

	router.HandleFunc("/github/event", getItemsData).Methods("GET")
	router.HandleFunc("/github/actors", getItemsData).Methods("GET")
	router.HandleFunc("/github/repos", getItemsData).Methods("GET")
	router.HandleFunc("/github/emails", getItemsData).Methods("GET")

	fmt.Printf("Starting server at port: 8080\n")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

func getItemsData(w http.ResponseWriter, r *http.Request) {
	var items interface{}
	switch r.RequestURI {
		case "/github/event":
			items = GetEventsDocs()

		case "/github/actors":
			general := GetActorsDocs()
			items = general.Actors

		case "/github/repos":
			items = GetRepoUrls()

		case "/github/emails":
			items = GetEmailsDocs()

		default:
			fmt.Printf("Error: invalid URL: %v\n", r.RequestURI)
			w.WriteHeader(http.StatusBadRequest)
			return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

