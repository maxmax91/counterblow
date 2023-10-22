package main

import "net/http/httputil"

// implemented:
// ip hash
// round robin
// for other policies is more complicated
// we need to keep track of server health

type RoutingRule struct {
	rule_id         string
	rule_type       string
	rule_ipaddr     string
	rule_subnetmask int
	rule_servers    string
}

// reflect the database structure
type BackendServer struct {
	addess       string
	priority     int
	errors       int
	timeouts     int
	ReverseProxy *httputil.ReverseProxy
}

func DecideRouting(from string) {
	// new connection from a client!
	// he want to connect with a server

}

func StartProxies() {

}
