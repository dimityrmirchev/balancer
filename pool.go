package main

import (
	"sync"
)

type pool struct {
	backends []backend
	current  int
	mutex    *sync.RWMutex
}

func (pool *pool) next() *backend {
	size := len(pool.backends)
	if size == 0 {
		return nil
	}

	if size == 1 && pool.backends[0].isAlive {
		return &pool.backends[0]
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

			return &pool.backends[i]
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

			return &pool.backends[i]
		}

		pool.current++
	}

	return nil
}
