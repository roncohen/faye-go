package faye

import (
	"fmt"
	"github.com/roncohen/faye/memory"
	"github.com/roncohen/faye/protocol"
	"log"
	"strconv"
)

type Engine struct {
	ns            memory.MemoryNamespace
	subscriptions memory.SubscriptionRegister
	clients       memory.ClientRegister
}

func NewEngine() Engine {
	return Engine{
		ns:            memory.NewMemoryNamespace(),
		clients:       memory.NewClientRegister(),
		subscriptions: memory.NewSubscriptionRegister(),
	}
}

func (m Engine) responseFromRequest(request protocol.Message) protocol.Message {
	response := protocol.Message{}
	response["channel"] = request.Channel().Name()
	if reqId, ok := request["id"]; ok {
		response["id"] = reqId.(string)
	}

	return response
}

func (m Engine) GetClient(clientId string) *protocol.Client {
	return m.clients.GetClient(clientId)
}

func (m Engine) Handshake(request protocol.Message, conn protocol.Connection) {
	// supportedConnectionType := request["supportedConnectionTypes"].([]interface{})
	version := request["version"].(string)

	response := m.responseFromRequest(request)
	response["successful"] = false

	if version == protocol.BAYEUX_VERSION {
		newClientId := m.ns.Generate()

		response.Update(map[string]interface{}{
			"clientId":                 newClientId,
			"channel":                  protocol.META_PREFIX + protocol.META_HANDSHAKE_CHANNEL,
			"version":                  protocol.BAYEUX_VERSION,
			"advice":                   protocol.DEFAULT_ADVICE,
			"supportedConnectionTypes": []string{"websocket", "long-polling"},
			"successful":               true,
		})

		msgStore := memory.NewMemoryMsgStore()
		newClient := protocol.NewClient(newClientId, msgStore)
		m.clients.AddClient(&newClient)
	} else {
		response["error"] = fmt.Sprintf("Only supported version is '%s'", protocol.BAYEUX_VERSION)
	}
	// Answer directly
	conn.Send([]protocol.Message{response})
}

func (m Engine) Connect(request protocol.Message, client *protocol.Client, conn protocol.Connection) {
	response := m.responseFromRequest(request)
	response["successful"] = true

	timeout, _ := strconv.Atoi(protocol.DEFAULT_ADVICE["timeout"])

	response.Update(protocol.Message{
		"advice": protocol.DEFAULT_ADVICE,
	})
	client.Connect(timeout, 0, response, conn)
}

func (m Engine) Subscribe(request protocol.Message, client *protocol.Client) {
	response := m.responseFromRequest(request)
	response["successful"] = true

	subscription := request["subscription"]
	response["subscription"] = subscription

	log.Printf("SUBSCRIBE %s subscription: %v", client.Id(), subscription)

	switch subscription.(type) {
	case []string:
		m.subscriptions.AddSubscription(client.Id(), subscription.([]string))
	case string:
		m.subscriptions.AddSubscription(client.Id(), []string{subscription.(string)})
	}

	client.Queue(response)
}

func (m Engine) Disconnect(request protocol.Message, client *protocol.Client, conn protocol.Connection) {
	response := m.responseFromRequest(request)
	response["successful"] = true
	clientId := request.ClientId()
	log.Printf("Client %s disconnected", clientId)
}

func (m Engine) Publish(request protocol.Message) {
	response := m.responseFromRequest(request)
	response["successful"] = true
	data := request["data"]
	channel := request.Channel()

	go func() {
		// Prepare msg to send to subscribers
		msg := protocol.Message{}
		msg["channel"] = channel
		msg["data"] = data
		// TODO: Missing ID

		// Get clients with subscriptions
		recipients := m.subscriptions.GetClients(channel.Expand())
		log.Print("PUBLISH to ", len(recipients), " clients")
		// Queue messages
		for _, c := range recipients {
			m.clients.GetClient(c).Queue(response)
		}
	}()
	m.clients.GetClient(request.ClientId()).Queue(response)
}
