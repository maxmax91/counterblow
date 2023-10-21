package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func startProxy(port string) {
	addr := flag.String("addr", "127.0.0.1:"+port, "listen address")
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
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func stopProxy() {

}
