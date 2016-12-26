package main

import (
	"net"
	"net/rpc"
	"time"
)

type (
	Client struct {
		connection  *rpc.Client
		getCount    int
		putCount    int
		deleteCount int
		clearCount  int
	}
)

func NewClient(dsn string, timeout time.Duration) (*Client, error) {
	connection, err := net.DialTimeout("tcp", dsn, timeout)
	if err != nil {
		return nil, err
	}
	return &Client{connection: rpc.NewClient(connection)}, nil
}

func (c *Client) Get(key string) (*CacheItem, error) {
	c.getCount++
	var item *CacheItem
	err := c.connection.Call("RPC.Get", key, &item)
	return item, err
}

func (c *Client) Put(item *CacheItem) (bool, error) {
	c.putCount++
	var added bool
	err := c.connection.Call("RPC.Put", item, &added)
	return added, err
}

func (c *Client) Delete(key string) (bool, error) {
	c.deleteCount++
	var deleted bool
	err := c.connection.Call("RPC.Delete", key, &deleted)
	return deleted, err
}

func (c *Client) Clear() (bool, error) {
	c.clearCount++
	var cleared bool
	err := c.connection.Call("RPC.Clear", true, &cleared)
	return cleared, err
}

func (c *Client) Stats() (*Requests, error) {
	requests := &Requests{}
	err := c.connection.Call("RPC.Stats", true, requests)
	return requests, err
}

func (c *Client) Reset() (bool, error) {
	var reset bool
	err := c.connection.Call("RPC.Reset", true, &reset)
	return reset, err
}