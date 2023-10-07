package main

import (
	"load-balancer/dummy_backend"
)

func main() {
	backendPods := []string{"pod1", "pod2", "pod3"}
	backendPort := 7070

	dummy_backend.StartDummyBackends(backendPort, backendPods...)
}
