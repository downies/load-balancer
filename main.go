package main

import (
	"load-balancer/dummy_backend"
	"load-balancer/load_balancer"
)

func main() {
	channel := make(chan int)

	backendPods := []string{"pod1", "pod2", "pod3"}
	backendPodsFail := []string{"pod1", "pod2", "pod3", "pod4"}
	backendPort := 7070
	loadBalancerPort := 80
	healthCheckPeriodInSeconds := 5

	go dummy_backend.StartDummyBackends(backendPort, backendPods...)
	load_balancer.StartLoadBalancer(loadBalancerPort, healthCheckPeriodInSeconds, backendPort, backendPodsFail...)

	<-channel
}
