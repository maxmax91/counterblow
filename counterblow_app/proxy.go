package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

var started bool = false
var pageCount int
var roundRobinRules []RoundRobinRoutingRule

// left for debug purposes - only starts a http web server
func startHttpServer(bindAddr string, port string) {

	if started == false {
		addr := flag.String("addr", bindAddr+":"+port, "listen address")
		flag.Parse()

		http.HandleFunc("/",
			func(w http.ResponseWriter, req *http.Request) {
				var b strings.Builder

				fmt.Fprintf(&b, "%v\t%v\t%v\tHost: %v\n", req.RemoteAddr, req.Method, req.URL, req.Host)
				for name, headers := range req.Header {
					for _, h := range headers {
						fmt.Fprintf(&b, "%v: %v\n", name, h)
					}
				}
				log.Println(b.String())

				fmt.Fprintf(w, "hello %s\n", req.URL)
			})

		log.Println("Starting server on", *addr)

		started = true
		// before 'cause this will be blocking!
		if err := http.ListenAndServe(*addr, nil); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	} else {
		log.Println("Proxy already started!")
	}
}

func stopProxy() {

}

func elaborateRules(rules []RoutingRule) {
	// round robin rules!
	// this function transforms every rule preparing the structures that can be used by algorithms
	for _, rule := range rules {

		if rule.rule_type == 1 {
			var rrRule RoundRobinRoutingRule // new rrrr
			rrRule.rule_ipaddr = rule.rule_ipaddr
			rrRule.rule_subnetmask = rule.rule_subnetmask
			rrRule.rule_servers = []BackendServer{}
			rrRule.rule_source = rule.rule_source
			rrRule.rule_dest = rule.rule_dest
			for _, el := range strings.Split(rule.rule_servers, ",") {
				var backendServer BackendServer
				backendServer.address = parseToUrl(el)
				rrRule.rule_servers = append(rrRule.rule_servers, backendServer)
				fmt.Printf("Loaded backend server url %s\n", backendServer.address)
			}
			fmt.Printf("Loaded rule: %v\n", rule)
			roundRobinRules = append(roundRobinRules, rrRule)
		}
	}
}

func startReverseProxy(listeningAddr string, listeningPort int, rules []RoutingRule) error {

	elaborateRules(rules)

	fromAddr := listeningAddr + ":" + fmt.Sprint(listeningPort)

	proxy := loadBalancingReverseProxy()
	log.Println("Starting proxy server on", fromAddr)
	if err := http.ListenAndServe(fromAddr, proxy); err != nil {
		log.Fatal("ListenAndServe:", err)
		return err
	}
	return nil
}

func stopReverseProxy() {
	// todo: implement STOP
}

// parseToUrl parses a "to" address to url.URL value
func parseToUrl(addr string) *url.URL {
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}
	toUrl, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}
	return toUrl
}

func loadBalancingReverseProxy() *httputil.ReverseProxy {

	director := func(req *http.Request) {

		var target *url.URL
		// Simple round robin between the two targets

		// for each round-robin rule.
		// it will try to follow each rewrite rule
		// or simply follow the rule if rule_source is not set
		var newUri string
		for id, rrrr := range roundRobinRules {
			// cannot use the second value! it is a copy of the element
			// se la regola è giusta avanziamo server
			reqUrl := req.URL.EscapedPath()
			println(fmt.Sprintf("Requested URI %s", reqUrl))
			if len(rrrr.rule_source) != 0 { // there is something in source regex filter
				m1 := regexp.MustCompile(rrrr.rule_source)
				if m1.MatchString(reqUrl) {
					println(fmt.Sprintf("Matching regex for rule %d!", id))
					newUri = m1.ReplaceAllString(reqUrl, rrrr.rule_dest)
					println(fmt.Sprintf("Rewriting uri from %s to %s", reqUrl, newUri))
				} else {
					println(fmt.Sprintf("Not matching regex for rule %d", id))
					continue // continue without exiting the loop (try with next rule)
				}
			} else {
				println(fmt.Sprintf("Redirecting w/o changes in uri %s for roundrobin rule %d", reqUrl, id))
				newUri = reqUrl
			}
			// send the request to the right server!
			roundRobinRules[id].current_server += 1
			roundRobinRules[id].current_server = roundRobinRules[id].current_server % len(roundRobinRules[id].rule_servers)
			target = roundRobinRules[id].rule_servers[roundRobinRules[id].current_server].address
			println("Found server for target")
			// found target
			break
		}

		if len(newUri) == 0 {

		}

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(target, newUri)
		pageCount += 1
		UpdateServedPages(pageCount)
		TextAreaLog(fmt.Sprintf("Served Url:%v RawPath:%v to Host:%v!", req.URL.Path, req.URL.RawPath, req.Host))
		// For simplicity, we don't handle RawQuery or the User-Agent header here:
		// see the full code of NewSingleHostReverseProxy for an example of doing
		// that.
	}
	return &httputil.ReverseProxy{Director: director}
}

// singleJoiningSlash is taken from net/http/httputil
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// joinURLPath is taken from net/http/httputil
func joinURLPath(a *url.URL, b string) (path, rawpath string) {
	if a.RawPath == "" && b == "" {
		return singleJoiningSlash(a.Path, b), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b, apath + "/" + bpath
	}
	return a.Path + b, apath + bpath
}
