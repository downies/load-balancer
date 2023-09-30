package load_balancer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var backends []string

var count int

func logRequestDetails(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request details
		fmt.Printf("Received request from %s\n", r.RemoteAddr)
		fmt.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
		fmt.Printf("Host: %s\n", r.Host)
		fmt.Printf("User-Agent: %s\n", r.UserAgent())
		fmt.Printf("Accept: %s\n\n", r.Header.Get("Accept"))

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func StartLoadBalancer(loadBalancerPort int, backendPort int, backendUrls ...string) {
	// Populate the list of backends
	backends = backendUrls

	// Create a new HTTP server
	server := http.Server{
		Addr: fmt.Sprintf(":%d", loadBalancerPort),
	}

	http.Handle("/", logRequestDetails(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new HTTP client
		client := http.Client{}
		resp, _ := client.Get(getBackendUrl(backendPort))

		// Read the response body into a string
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading backend response", http.StatusInternalServerError)
			return
		}

		// Write the response body string to the client
		w.WriteHeader(http.StatusOK)
		w.Write(body)

	})))

	// Start the server
	log.Fatal(server.ListenAndServe())
}

func getBackendUrl(backendPort int) string {
	count = (count + 1) % len(backends)
	return "http://localhost:" + fmt.Sprintf("%d", backendPort) + "/" + backends[count]
}
