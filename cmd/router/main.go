package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gitbao/gitbao/router"
)

func main() {

	// h := httputil.NewSingleHostReverseProxy(u)
	handler := Handler()
	s := &http.Server{
		Addr:           ":8002",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Router up and running on port 8002")
	log.Fatal(s.ListenAndServe())
}

func Handler() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		var subdomain string
		host := req.Host
		host_parts := strings.Split(host, ".")
		destination := "gitbao.com"

		if len(host_parts) > 2 {
			subdomain = host_parts[0]
			var err error
			destination, err = router.GetDestinaton(subdomain)
			if err != nil {
				log.Println(err)
				destination = "gitbao.com"
			}
		}
		target, err := url.Parse("http://" + destination)
		if err != nil {
			log.Println("Bad destination.")
			return
		}

		targetQuery := target.RawQuery
		log.Printf("%s %s %s %s", req.Method, req.Host, subdomain, req.URL)

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

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
