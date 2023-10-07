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

var globalBackends []backend

var count int

func logRequestDetails(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request details
		log.Printf("Received request from %s\n\n", r.RemoteAddr)
		log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
		log.Printf("Host: %s\n", r.Host)
		log.Printf("User-Agent: %s\n", r.UserAgent())
		log.Printf("Accept: %s\n\n", r.Header.Get("Accept"))

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func StartLoadBalancer(loadBalancerPort int, healthCheckPeriodInSeconds int, backendPort int, backendUrls ...string) {
	populateBackends(backendUrls, backendPort)

	go startHealthChecks(healthCheckPeriodInSeconds)

	startLoadBalancer(loadBalancerPort, backendPort)
}

func startLoadBalancer(loadBalancerPort int, backendPort int) {
	server := http.Server{
		Addr: fmt.Sprintf(":%d", loadBalancerPort),
	}

	http.Handle("/", logRequestDetails(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new HTTP client
		client := http.Client{}
		url := getBackendUrl()
		log.Println(url)
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

	log.Fatal(server.ListenAndServe())
}

func populateBackends(backendUrls []string, backendPort int) {
	for _, url := range backendUrls {
		newBackend := backend{
			url:     fmt.Sprintf("http://localhost:" + fmt.Sprintf("%d", backendPort) + "/" + url),
			healthy: true,
		}
		globalBackends = append(globalBackends, newBackend)
	}
}

func startHealthChecks(healthCheckPeriodInSeconds int) {
	for i := range globalBackends {
		go startScheduledHealthChecks(&globalBackends[i], healthCheckPeriodInSeconds)
	}
}

func startScheduledHealthChecks(b *backend, healthCheckPeriodInSeconds int) {
	for {
		verifyBackendHealthy(b)
		time.Sleep(time.Duration(healthCheckPeriodInSeconds) * time.Second)
	}
}

func verifyBackendHealthy(b *backend) {
	url := b.url
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		markBackendHealthStatus(b, false)
		return
	}
	status := res.StatusCode
	if status == http.StatusOK {
		if !b.healthy {
			markBackendHealthStatus(b, true)
		}
		return
	}

	if b.healthy {
		markBackendHealthStatus(b, false)
	}
}

func markBackendHealthStatus(b *backend, healthy bool) {
	if healthy {
		log.Printf("Backend %s is healthy again, adding back to the pool\n", b.url)
	} else {
		log.Printf("Backend %s is unhealthy, removing from the pool\n", b.url)
	}
	b.healthy = healthy
}

func getBackendUrl() string {
	count = incrementCounter()
	for !globalBackends[count].healthy {
		count = incrementCounter()
	}
	return globalBackends[count].url
}

func incrementCounter() int {
	return (count + 1) % len(globalBackends)
}
