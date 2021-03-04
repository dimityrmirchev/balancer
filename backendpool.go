package main

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type backendPool struct {
	backends []backend
	current  int
	mutex    *sync.RWMutex
}

func (pool *backendPool) next() *httputil.ReverseProxy {
	size := len(pool.backends)
	if size == 0 {
		return nil
	}

	if size == 1 && pool.backends[0].isAlive {
		return pool.backends[0].proxy
	}

	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	ind := pool.current

	// traverse from current to end
	for i := ind; i < size; i++ {

		if pool.backends[i].isAlive {
			if pool.current+1 < size {
				pool.current = pool.current + 1
			} else {
				pool.current = 0
			}

			return pool.backends[i].proxy
		}
		pool.current++
	}

	// traverse from start to current in case we did not find an alive backend service
	pool.current = 0
	for i := 0; i < ind; i++ {
		if pool.backends[i].isAlive {
			if pool.current+1 < size {
				pool.current = pool.current + 1
			} else {
				pool.current = 0
			}

			return pool.backends[i].proxy
		}

		pool.current++
	}

	return nil
}

func (pool *backendPool) markBackendStatus(url *url.URL, status bool) {
	for i, backend := range pool.backends {
		if backend.url.Host == url.Host && backend.url.Port() == url.Port() {
			pool.mutex.Lock()
			defer pool.mutex.Unlock()
			pool.backends[i].isAlive = status
			return
		}
	}
}
