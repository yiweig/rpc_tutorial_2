package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

type (
	// RPC server struct.
	RPC struct {
		cache    map[string]string
		requests *Requests
		mu       *sync.RWMutex
	}

	// RPCForREST will call a REST endpoint.
	RPCForREST struct {
		URL    string
		apiKey string
	}

	// CacheItem for our cache.
	CacheItem struct {
		Key   string
		Value string
	}

	// Requests object for making calls.
	Requests struct {
		Get    int
		Put    int
		Delete int
		Clear  int
	}
)

var (
	// ErrNotFound for cache misses.
	ErrNotFound = errors.New("Cache key not found")
)

// CreateNewRPC constructor.
func CreateNewRPC() *RPC {
	return &RPC{
		cache:    make(map[string]string),
		requests: &Requests{},
		mu:       &sync.RWMutex{},
	}
}

func CreateNewRPCForREST() *RPCForREST {
	return &RPCForREST{
		// "https://api.nytimes.com/svc/topstories/v2/home.json?api-key=1865cb7220b54b56aedd9e8ab53130c9"
		URL:    "https://api.nytimes.com/svc/topstories/v2/home.json",
		apiKey: "1865cb7220b54b56aedd9e8ab53130c9",
	}
}

// func (r *RPCForREST) GetTopStories(_ interface{}, resp *map[string]interface{}) (err error) {
// 	client := &http.Client{Timeout: 10 * time.Second}
// 	var url = r.URL + "?api-key=" + r.apiKey
// 	response, err := client.Get(url)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer response.Body.Close()

// 	return json.NewDecoder(response.Body).Decode(resp)
// }

func (r *RPCForREST) GetTopStories(url string, target *map[string]interface{}) error {
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})
	var myClient = &http.Client{Timeout: 10 * time.Second}
	resp, err := myClient.Get(r.URL + "?api-key=" + r.apiKey)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// Get from the cache
func (r *RPC) Get(key string, resp *CacheItem) (err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cacheValue, found := r.cache[key]

	if !found {
		return ErrNotFound
	}

	*resp = CacheItem{key, cacheValue}
	r.requests.Get++
	return nil
}

// Put into the cache
func (r *RPC) Put(item *CacheItem, ack *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cache[item.Key] = item.Value
	*ack = true

	r.requests.Put++
	return nil
}

// Delete from the cache
func (r *RPC) Delete(key string, ack *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var found bool
	_, found = r.cache[key]

	if !found {
		return ErrNotFound
	}

	delete(r.cache, key)
	*ack = true

	r.requests.Delete++
	return nil
}

// Clear the cache
func (r *RPC) Clear(skip bool, ack *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cache = make(map[string]string)
	*ack = true

	r.requests.Clear++
	return nil
}

// Stats will get the stats.
func (r *RPC) Stats(skip bool, requests *Requests) error {
	*requests = *r.requests
	return nil
}

// Reset will reset all values on the server to 0.
func (r *RPC) Reset(_ bool, ack *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.requests.Get = 0
	r.requests.Put = 0
	r.requests.Delete = 0
	r.requests.Clear = 0

	*ack = true

	return nil
}
