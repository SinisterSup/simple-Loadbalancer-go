# Simple Load Balancer

This project implements a basic load balancer in Go, which forwards incoming HTTP requests to a list of predefined backend servers using a round-robin algorithm. The load balancer listens on a specified port and can be easily extended to include health checks or other load balancing algorithms.

## Code Overview

The `main.go` file consists of:

- `Server`: An interface representing a backend server, with methods for retrieving its address, checking its availability, and serving an HTTP request.
- `simpleServer`: A basic implementation of the `Server` interface that uses `httputil.ReverseProxy` to serve requests.
- `LoadBalancer`: A struct representing the load balancer, which keeps track of a list of backend servers and implements a round-robin algorithm to distribute incoming requests.
- `main()`: The entry point of the program, which initializes the backend servers, sets up the load balancer, and starts listening for incoming requests.

## Usage

To use the load balancer, simply compile and run the `main.go` file:

``` bash
go build main.go
./main
```

By default, the load balancer will listen on port 8000 and forward requests to the following backend servers:

- https://www.facebook.com
- https://www.bing.com
- https://www.duckduckgo.com

You can modify the list of backend servers and the listening port by updating the corresponding variables in the `main()` function.

To test the load balancer, send requests to `localhost:8000` using your preferred HTTP client, such as a web browser or `curl`:

```bash
curl http://localhost:8000
```

The load balancer will forward the request to one of the backend servers and return their response.
