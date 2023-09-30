package dummy_backend

import (
	"log"
	"net/http"
)

func StartDummyBackend() {
	// Create a new HTTP server
	server := http.Server{
		Addr: ":7070",
	}

	http.Handle("/backend", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from backend!"))
	}))

	// Start the server
	log.Fatal(server.ListenAndServe())
}
