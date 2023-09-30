package main

import (
	"load-balancer/dummy_backend"
	"load-balancer/load_balancer"
)

func main() {
	channel := make(chan int)

	go dummy_backend.StartDummyBackend("7070")
	go dummy_backend.StartDummyBackend("7071")
	go load_balancer.StartLoadBalancer("7070", "7071")

	<-channel
}
