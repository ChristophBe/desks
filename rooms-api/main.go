package main

import (
	"log"
	"net/http"
)
import "github.com/gorilla/mux"

func main() {
	router := mux.NewRouter()
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Failed to start Api", err)
	}
}
