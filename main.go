package main

import (
	"load-balancer/dummy_backend"
	"load-balancer/load_balancer"
)

func main() {
	channel := make(chan int)

	go dummy_backend.StartDummyBackend()
	go load_balancer.StartLoadBalancer()

	<-channel
}
