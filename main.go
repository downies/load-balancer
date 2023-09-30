package main

import (
	"load-balancer/dummy_backend"
	"load-balancer/load_balancer"
)

func main() {
	channel := make(chan int)

	backendPods := []string{"pod1", "pod2", "pod3"}
	backendPort := 7070
	loadBalancerPort := 80

	go dummy_backend.StartDummyBackends(backendPort, backendPods...)
	go load_balancer.StartLoadBalancer(loadBalancerPort, backendPort, backendPods...)

	<-channel
}
