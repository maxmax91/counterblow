package main

import (
	"net/http/httputil"
)

// implemented:
// ip hash
// round robin
// for other policies is more complicated
// we need to keep track of server health

type BackendServer struct {
	addess       string
	priority     int
	errors       int
	timeouts     int
	ReverseProxy *httputil.ReverseProxy
}

type RoutingRule struct {
	from   string
	to     string
	policy string
}

func DecideRouting(from string) {
	// new connection from a client!
	// he want to connect with a server

}

func StartProxies() {

}
