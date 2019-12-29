package db

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
)

var myDB = map[string]map[string]bool{
	"ynet.co.il": {
		"/page=1": true,
		"/page=2": true,
	},
	"shesh.co.il": {
		"/chat": true,
	},
}

const domainsMaxCacheSize = 10
const pathsMaxCacheSize = 10

var (
	_cache *lru.Cache
)

func init() {
	_cache, _ = lru.NewWithEvict(domainsMaxCacheSize, abortDomainHandler)
}

func createNewCache(size int) *lru.Cache {
	var (
		cache *lru.Cache
		err error
	)

	if cache, err = lru.New(size); err != nil {
		panic("Failed initializing cache!")
	}

	return cache
}

type DomainPathRequest struct {
	Path string
	Reply chan bool
}

type cacheEntry struct {
	requests chan *DomainPathRequest
	abort chan struct{}
}

func getDomainPath(domain string, requests chan *DomainPathRequest, abort chan struct{}) {
	cache := createNewCache(pathsMaxCacheSize)

	for {
		select {
		case <-abort:
			// Die
			return
		case request := <-requests:
			// Return from cache
			if pathFromDB, exists := cache.Get(request.Path); exists {
				request.Reply <- pathFromDB.(bool)
				continue
			}

			// Go to DB and cache the response
			exists := expensivelyCheckDb(domain, request.Path)
			cache.Add(request.Path, exists)
			request.Reply <- exists
		}
	}
}

func abortDomainHandler(key interface{}, value interface{}) {
	entry := value.(*cacheEntry)
	close(entry.abort)
}

func expensivelyCheckDb(domain, path string) bool {
	fmt.Println("<DB HIT>")
	return myDB[domain][path]
}

func sendRequest(domain, path string) (result bool, sent bool) {
	domainData, exists := _cache.Get(domain)
	if !exists {
		return false, false
	}

	// Prepare request
	reply := make(chan bool)
	request := &DomainPathRequest{
		Path:  path,
		Reply: reply,
	}

	// Send it
	select {
	case domainData.(*cacheEntry).requests <- request:
		result = <-reply
		return result, true
	case <-domainData.(*cacheEntry).abort:
		return false, false
	}
}

func Get(domain, path string) bool {
	result, sent := sendRequest(domain, path)
	if sent {
		return result
	}

	for !sent {
		// Creating new goroutine to handle this request
		domainData := cacheEntry{
			requests: make(chan *DomainPathRequest),
			abort:    make(chan struct{}),
		}

		_cache.Add(domain, &domainData)
		go getDomainPath(domain, domainData.requests, domainData.abort)
		result, sent = sendRequest(domain, path)
	}

	return result
}
