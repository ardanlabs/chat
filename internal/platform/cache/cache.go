package cache

import (
	"fmt"
	"net"
	"sync"
)

// Client represents a connected person in the server.
type Client struct {
	ID      string
	TCPAddr *net.TCPAddr
}

// Cache maintains client connections
type Cache struct {
	clients   map[string]Client
	addresses map[string]string
	mu        sync.Mutex
}

// New returns a cache value ready for use.
func New() *Cache {
	return &Cache{
		clients:   make(map[string]Client),
		addresses: make(map[string]string),
	}
}

// Get returns the current set of clients.
func (c *Cache) Get(id string) []Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	clients := make([]Client, 0, len(c.clients))
	for _, client := range c.clients {
		if client.ID != id {
			clients = append(clients, client)
		}
	}

	return clients
}

// Add adds a client value to the cache.
func (c *Cache) Add(id string, tcpAddr *net.TCPAddr) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.clients[id]; exists {
		return fmt.Errorf("client [ %s ] already exists", id)
	}

	client := Client{
		ID:      id,
		TCPAddr: tcpAddr,
	}

	c.clients[id] = client
	c.addresses[tcpAddr.String()] = id

	return nil
}

// GetID find the client value by id.
func (c *Cache) GetID(id string) (Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	client, exists := c.clients[id]
	if !exists {
		return Client{}, fmt.Errorf("client [ %s ] does not exist", id)
	}

	return client, nil
}

// GetAddress find the client value by address.
func (c *Cache) GetAddress(address string) (Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id, exists := c.addresses[address]
	if !exists {
		return Client{}, fmt.Errorf("client [ %s ] does not exist", address)
	}

	client, exists := c.clients[id]
	if !exists {
		return Client{}, fmt.Errorf("client [ %s ] does not exist", id)
	}

	return client, nil
}

// Remove removes a client value from the cache.
func (c *Cache) Remove(address string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	id, exists := c.addresses[address]
	if !exists {
		return fmt.Errorf("client [ %s ] does not exist", address)
	}

	delete(c.addresses, address)
	delete(c.clients, id)

	return nil
}
