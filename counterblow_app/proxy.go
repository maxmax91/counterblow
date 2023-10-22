package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var started bool = false
var pageCount int

// for debug purposes?
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

func startReverseProxy(listeningAddr string, listeningPort int) error {
	fromAddr := flag.String("from", listeningAddr+":"+fmt.Sprint(listeningPort), "proxy's listening address")
	toAddr1 := flag.String("to1", "google.it:80", "the first address this proxy will forward to")
	toAddr2 := flag.String("to2", "microsoft.it:80", "the second address this proxy will forward to")
	flag.Parse()

	toUrl1 := parseToUrl(*toAddr1)
	toUrl2 := parseToUrl(*toAddr2)

	proxy := loadBalancingReverseProxy(toUrl1, toUrl2)
	log.Println("Starting proxy server on", *fromAddr)
	if err := http.ListenAndServe(*fromAddr, proxy); err != nil {
		log.Fatal("ListenAndServe:", err)
		return err
	}
	return nil
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

func loadBalancingReverseProxy(target1, target2 *url.URL) *httputil.ReverseProxy {
	var targetNum = 1

	director := func(req *http.Request) {
		var target *url.URL
		// Simple round robin between the two targets
		if targetNum == 1 {
			target = target1
			targetNum = 2
		} else {
			target = target2
			targetNum = 1
		}

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
		pageCount += 1
		UpdateServedPages(pageCount)
		TextAreaLog(fmt.Sprintf("Served Url:%v RawPath:%v from Host:%v!", req.URL.Path, req.URL.RawPath, req.Host))
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
func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}
