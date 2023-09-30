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

func StartLoadBalancer(ports ...string) {
	backends = ports
	// Create a new HTTP server
	server := http.Server{
		Addr: ":8080",
	}

	http.Handle("/", logRequestDetails(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new HTTP client
		client := http.Client{}
		resp, _ := client.Get(getbackendurl())
		getbackendurl()
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

func getbackendurl() string {
	count = (count + 1) % len(backends)
	fmt.Println(backends[count])
	return "http://localhost:" + backends[count] + "/backend-" + backends[count]
}
