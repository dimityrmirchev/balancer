package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var pool backendPool

func main() {
	port := flag.Int("port", 3001, "port to listen on")
	flag.Parse()
	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	backends := parseBackends(config.Backends)
	pool = backendPool{backends: backends, mutex: new(sync.RWMutex)}
	server := http.Server{
		Addr:    ":" + fmt.Sprint(*port),
		Handler: http.HandlerFunc(balance),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func balance(w http.ResponseWriter, req *http.Request) {
	proxy := pool.next()

	if proxy == nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	proxy.ErrorHandler = handleError

	proxy.ServeHTTP(w, req)
}

func handleError(w http.ResponseWriter, req *http.Request, err error) {
	log.Printf("Host %s is unavailable on port %s. Error: %s", req.URL.Host, req.URL.Port(), err.Error())

	pool.markBackendStatus(req.URL, false)
	balance(w, req)
}

func parseBackends(urls []string) []backend {
	backends := []backend{}
	for _, conf := range urls {
		res, err := url.Parse(conf)
		if err == nil {
			backend := backend{}
			backend.isAlive = true
			backend.url = res
			backend.proxy = httputil.NewSingleHostReverseProxy(res)
			backends = append(backends, backend)
		}

	}
	return backends
}
