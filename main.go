package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	// Address returns the address with which to access the server
	Address() string

	// IsAlive returns true if the server is alive and able to serve requests
	IsAlive() bool

	// Serve uses this server to process the request
	Serve(rw http.ResponseWriter, req *http.Request)
}

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

func (s *simpleServer) Address() string { return s.addr }

func (s *simpleServer) IsAlive() bool { return true }

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	handleErr(err)

	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

// handleErr prints the error and exits the program
// Note: this is not how one would want to handle errors in production, but
// serves well for demonstration purposes.
func handleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

// getNextServerAddr returns the address of the next available server to send a
// request to, using a simple round-robin algorithm
func (lb *LoadBalancer) getNextAvailableServer() Server {
	startIndex := lb.roundRobinCount % len(lb.servers)

	for i := startIndex; i < startIndex+len(lb.servers); i++ {
		server := lb.servers[i%len(lb.servers)]
		if server.IsAlive() {
			lb.roundRobinCount = i % len(lb.servers)
			return server
		}
	}

	// fallback to the first server if all servers are dead
	lb.roundRobinCount = 0
	return lb.servers[0]
}

func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, req *http.Request) {
	targetServer := lb.getNextAvailableServer()

	// could optionally log stuff about the request here!
	fmt.Printf("forwarding request to address %q\n", targetServer.Address())

	// could delete pre-existing X-Forwarded-For header to prevent IP spoofing
	targetServer.Serve(rw, req)
}

func main() {
	servers := []Server{
		newSimpleServer("https://www.facebook.com"),
		newSimpleServer("https://www.bing.com"),
		newSimpleServer("https://www.duckduckgo.com"),
	}

	lb := NewLoadBalancer("8000", servers)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serveProxy(rw, req)
	}

	// register a proxy handler to handle all requests
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("serving requests at 'localhost:%s'\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
