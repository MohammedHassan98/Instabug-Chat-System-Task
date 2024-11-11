// cmd/server/main.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
    // Set up router
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	
    // Start server
    log.Println("Starting server on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}
