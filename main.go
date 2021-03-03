package main

import (
	"log"
	"net/url"
)

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	backends := []backend{}
	for _, conf := range config.Backends {
		res, err := url.Parse(conf)
		if err == nil {
			backend := backend{}
			backend.isAlive = true
			backend.url = res
			backends = append(backends, backend)
		}

	}

	//pool := pool{backends: backends, mutex: new(sync.RWMutex)}

}
