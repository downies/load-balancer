package dummy_backend

import (
	"fmt"
	"log"
	"net/http"
)

func StartDummyBackends(port int, urls ...string) {
	// Create a new HTTP server
	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	// Register a handler for each URL pattern
	for _, url := range urls {
		path := "/" + url
		message := "hello from backend: " + url

		http.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(message))
		}))
	}

	// Start the server
	fmt.Println("Serving dummy backends on port", port)
	log.Fatal(server.ListenAndServe())
}
