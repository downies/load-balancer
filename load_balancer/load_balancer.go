package load_balancer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type backend struct {
	healthy bool
	url     string
}

var backends []backend

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

func StartLoadBalancer(loadBalancerPort int, healthCheckPeriodInSeconds int, backendPort int, backendUrls ...string) {
	// Populate the list of backends
	for _, url := range backendUrls {
		newBackend := backend{
			url:     url,
			healthy: true,
		}
		backends = append(backends, newBackend)
	}
	go RunHealthCheck(healthCheckPeriodInSeconds)
	// Create a new HTTP server
	server := http.Server{
		Addr: fmt.Sprintf(":%d", loadBalancerPort),
	}

	http.Handle("/", logRequestDetails(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new HTTP client
		client := http.Client{}
		url := getBackendUrl(backendPort)
		fmt.Println(url)
		resp, _ := client.Get(url)

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

func RunHealthCheck(healthCheckPeriodInSeconds int) {
	for i, _ := range backends {
		go runSchedule(i, healthCheckPeriodInSeconds)
	}
}

func runSchedule(i int, healthCheckPeriodInSeconds int) {
	for {
		verifyHealthy(i)
		time.Sleep(time.Duration(healthCheckPeriodInSeconds) * time.Second)
	}
}

func verifyHealthy(i int) {

	url := "http://localhost:" + "7070" + "/" + backends[i].url
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	status := res.StatusCode
	if status != http.StatusOK {
		fmt.Printf("removing pod: %s\n", backends[i].url)
		backends[i].healthy = false
	}
}

func getBackendUrl(backendPort int) string {
	count = (count + 1) % len(backends)
	return "http://localhost:" + fmt.Sprintf("%d", backendPort) + "/" + backends[count].url
}
