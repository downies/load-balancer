package dummy_backend

import (
	"fmt"
	"log"
	"net/http"
)

func StartDummyBackend(port string) {
	// Create a new HTTP server
	server := http.Server{
		Addr: fmt.Sprint(":" + port),
	}

	http.Handle("/backend", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from backend running at: " + port))
	}))

	// Start the server
	log.Fatal(server.ListenAndServe())
}
