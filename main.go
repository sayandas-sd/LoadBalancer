package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
	isRunning() bool
	Serve(w http.ResponseWriter, r *http.Request)
}

type NewServer struct {
	addr  string
	proxy httputil.ReverseProxy
}

func simpleServer(addr string) *NewServer {
	server, err := url.Parse(addr)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	if err != nil {

	}

	return &NewServer{
		addr:  addr,
		proxy: *httputil.NewSingleHostReverseProxy(server),
	}

}

type LoadBalancer struct {
	port       string
	roundRobin int
	servers    []Server
}

func newLoadBalancer(port string, servers []Server) *LoadBalancer {

	return &LoadBalancer{
		port:       port,
		roundRobin: 0,
		servers:    servers,
	}
}
func (s *NewServer) Address() string { return s.addr }

func (s *NewServer) isRunning() bool { return true }

func (s *NewServer) Serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}

func (lb *LoadBalancer) getServer() Server {
	server := lb.servers[lb.roundRobin%len(lb.servers)]

	for !server.isRunning() {
		lb.roundRobin++
		server = lb.servers[lb.roundRobin%len(lb.servers)]
	}

	lb.roundRobin++
	return server
}

func (lb *LoadBalancer) serverProxy(w http.ResponseWriter, r *http.Request) {
	target := lb.getServer()
	fmt.Printf("forwaring request to address %q\n", target.Address())
	target.Serve(w, r)
}

func main() {
	servers := []Server{
		simpleServer("https://www.facebook.com"),
		simpleServer("https://www.duckduckgo.com/"),
		simpleServer("https://www.bing.com"),
	}

	lb := newLoadBalancer("3000", servers)

	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		lb.serverProxy(w, r)
	}

	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Server running at http://localhost:%s\n", lb.port)

	err := http.ListenAndServe(":"+lb.port, nil)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
