# Round Robin Load Balancer in Golang

This is a simple Go project that demonstrates how a Round Robin Load Balancer works. The load balancer distributes incoming requests to a set of dummy pods in a round-robin fashion. This project is intended for learning purposes.

```bash
go run db.go # start dummy pods
go run lb.go # start load balancer
```


```bash
curl --location 'http://localhost' # load balancer
curl --location 'http://localhost:7070/<pod_name>' # pod
```
