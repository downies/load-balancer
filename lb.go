package main

import (
	"load-balancer/load_balancer"
)

func main() {
	backendPods := []string{"pod1", "pod2", "pod3"}
	backendPodsFail := append(backendPods, "pod4")

	backendPort := 7070
	loadBalancerPort := 80
	healthCheckPeriodInSeconds := 5

	load_balancer.StartLoadBalancer(loadBalancerPort, healthCheckPeriodInSeconds, backendPort, backendPodsFail...)
}
