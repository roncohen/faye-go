package memory

import (
	"github.com/roncohen/faye/protocol"
	"sync"
)

type ClientRegister struct {
	clients map[string]*protocol.Client
	mutex   sync.RWMutex
}

func NewClientRegister() ClientRegister {
	return ClientRegister{
		clients: make(map[string]*protocol.Client),
	}
}

func (cr ClientRegister) AddClient(client *protocol.Client) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	cr.clients[client.Id()] = client
}

func (cr ClientRegister) RemoveClient(clientId string) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	delete(cr.clients, clientId)
}

func (cr ClientRegister) GetClient(clientId string) *protocol.Client {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	return cr.clients[clientId]
}
