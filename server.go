package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/hello", Hello)
	log.Printf("hello-nomock listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Hello, World!!!\"}"))
}
