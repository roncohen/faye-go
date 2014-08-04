package memory

import (
	"sync"
)

type ClientSet struct {
	clients map[string]bool
}

func NewClientSet() ClientSet {
	return ClientSet{make(map[string]bool)}
}

func (c ClientSet) Add(clientId string) {
	c.clients[clientId] = true
}

func (c ClientSet) AddMany(clientIds []string) {
	for _, clientId := range clientIds {
		c.clients[clientId] = true
	}
}

func (c ClientSet) Remove(clientId string) {
	delete(c.clients, clientId)
}

func (c ClientSet) GetAll() []string {
	all := make([]string, len(c.clients))
	i := 0
	for k := range c.clients {
		all[i] = k
		i = i + 1
	}
	return all
}

type SubscriptionRegister struct {
	subscriptions map[string]ClientSet
	mutex         sync.RWMutex
}

func NewSubscriptionRegister() SubscriptionRegister {
	return SubscriptionRegister{
		subscriptions: make(map[string]ClientSet),
	}
}

func (sr SubscriptionRegister) AddSubscription(clientId string, patterns []string) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	for _, pattern := range patterns {
		_, ok := sr.subscriptions[pattern]
		if !ok {
			sr.subscriptions[pattern] = NewClientSet()
		}
		sr.subscriptions[pattern].Add(clientId)
	}
}

func (sr SubscriptionRegister) RemoveSubscription(clientId string, patterns []string) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	for _, pattern := range patterns {
		sr.subscriptions[pattern].Remove(clientId)
	}
}

func (sr SubscriptionRegister) GetClients(patterns []string) []string {
	clientSet := NewClientSet()
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()

	for _, pattern := range patterns {
		clientSet.AddMany(sr.subscriptions[pattern].GetAll())
	}
	return clientSet.GetAll()
}
