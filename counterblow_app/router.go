package main

import (
	"net/url"
)

// implemented:
// ip hash
// round robin
// for other policies is more complicated
// we need to keep track of server health

type RoutingRule struct {
	rule_id         string
	rule_type       int
	rule_ipaddr     string
	rule_subnetmask int
	rule_servers    string
	rule_source     string // requested page. Can be (.*)
	rule_dest       string // rewrited page. Can be $1
}

type RoundRobinRoutingRule struct { // extends RoutingRule -- how to do it in go?
	rule_id         string
	rule_type       string
	rule_ipaddr     string
	rule_subnetmask int
	rule_servers    []BackendServer
	current_server  int    // defines the current server
	rule_source     string // requested page. Can be (.*)
	rule_dest       string // rewrited page. Can be $1
}

// reflect the database structure
type BackendServer struct {
	address  *url.URL
	priority int
	errors   int
	timeouts int
	active   bool // maybe it must be skipped if it is not active
	//ReverseProxy *httputil.ReverseProxy
}

func DecideRouting(from string) {
	// new connection from a client!
	// he want to connect with a server

}

func StartProxies() {

}
