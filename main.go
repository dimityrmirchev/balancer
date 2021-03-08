package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

var pool backendPool
var ticker *time.Ticker
var quit chan struct{}

func main() {
	port := flag.Int("port", 3001, "Port to listen on")
	backendsRawVal := flag.String("backends", "", "A list of backends to balance between separated with ','")
	flag.Parse()

	if *backendsRawVal == "" {
		log.Fatal("Please provide value for the backends flag")
	}

	splitBackends := strings.Split(*backendsRawVal, ",")
	backends := parseBackends(splitBackends)

	if len(backends) == 0 {
		log.Fatal("Please provide valid backend service urls")
	}

	pool = backendPool{backends: backends, mutex: new(sync.RWMutex)}
	server := http.Server{
		Addr:    ":" + fmt.Sprint(*port),
		Handler: http.HandlerFunc(balance),
	}

	registerHealthChecks()

	fmt.Println("Balancer listening on port " + fmt.Sprint(*port))
	err := server.ListenAndServe()
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

func registerHealthChecks() {
	ticker = time.NewTicker(5 * time.Second)
	quit = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, v := range pool.backends {
					timeout := time.Second
					conn, err := net.DialTimeout("tcp", v.url.Host, timeout)
					if err != nil {
						log.Printf("Cannot establish connection. Error: %s", err.Error())
						pool.markBackendStatus(v.url, true)
					} else if conn != nil {
						defer conn.Close()
						pool.markBackendStatus(v.url, true)
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
