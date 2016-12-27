package main

import (
	"encoding/gob"
	"log"
	"net"
	"net/rpc"
	"time"
)

type (
	// Client that we use to connect to the server.
	Client struct {
		connection  *rpc.Client
		getCount    int
		putCount    int
		deleteCount int
		clearCount  int
	}
)

// CreateNewClient for the server.
func CreateNewClient(dsn string, timeout time.Duration) (*Client, error) {
	connection, err := net.DialTimeout("tcp", dsn, timeout)
	if err != nil {
		return nil, err
	}
	return &Client{connection: rpc.NewClient(connection)}, nil
}

// GetTopStories calls the NYT Top Stories API.
func (c *Client) GetTopStories() (map[string]interface{}, error) {
	// http://stackoverflow.com/a/21939636/1470257
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})

	params := "url"
	dat := make(map[string]interface{})
	log.Println("GetTopStories called and dat created")
	err := c.connection.Call("RPCForREST.GetTopStories", params, &dat)
	log.Println("rpc called!")
	return dat, err
}

// Get from the server.
func (c *Client) Get(key string) (*CacheItem, error) {
	c.getCount++
	var item *CacheItem
	err := c.connection.Call("RPC.Get", key, &item)
	return item, err
}

// Put into the server.
func (c *Client) Put(item *CacheItem) (bool, error) {
	c.putCount++
	var added bool
	err := c.connection.Call("RPC.Put", item, &added)
	return added, err
}

// Delete from the server.
func (c *Client) Delete(key string) (bool, error) {
	c.deleteCount++
	var deleted bool
	err := c.connection.Call("RPC.Delete", key, &deleted)
	return deleted, err
}

// Clear all in the server.
func (c *Client) Clear() (bool, error) {
	c.clearCount++
	var cleared bool
	err := c.connection.Call("RPC.Clear", true, &cleared)
	return cleared, err
}

// Stats all the things.
func (c *Client) Stats() (*Requests, error) {
	requests := &Requests{}
	err := c.connection.Call("RPC.Stats", true, requests)
	return requests, err
}

// Reset all counters.
func (c *Client) Reset() (bool, error) {
	var reset bool
	err := c.connection.Call("RPC.Reset", true, &reset)
	return reset, err
}
