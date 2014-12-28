package faye

import (
	"fmt"
	"github.com/roncohen/faye/memory"
	"github.com/roncohen/faye/protocol"
	"github.com/roncohen/faye/transport"
	"github.com/roncohen/faye/utils"
	"strconv"
)

type Engine struct {
	ns      memory.MemoryNamespace
	clients *memory.ClientRegister
	logger  utils.Logger
}

func NewEngine(logger utils.Logger) Engine {
	return Engine{
		ns:      memory.NewMemoryNamespace(),
		clients: memory.NewClientRegister(logger),
		logger:  logger,
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

func (m Engine) NewClient(conn protocol.Connection) *protocol.Client {
	newClientId := m.ns.Generate()
	msgStore := memory.NewMemoryMsgStore()
	newClient := protocol.NewClient(newClientId, msgStore, m.logger)
	m.clients.AddClient(&newClient)
	return &newClient
}

func (m Engine) AddSubscription(clientId string, subscriptions []string) {
	m.logger.Infof("SUBSCRIBE %s subscription: %v", clientId, subscriptions)
	m.clients.AddSubscription(clientId, subscriptions)
}

func (m Engine) Handshake(request protocol.Message, conn protocol.Connection) string {
	newClientId := ""
	version := request["version"].(string)

	response := m.responseFromRequest(request)
	response["successful"] = false

	if version == protocol.BAYEUX_VERSION {
		newClientId = m.NewClient(conn).Id()

		response.Update(map[string]interface{}{
			"clientId":                 newClientId,
			"channel":                  protocol.META_PREFIX + protocol.META_HANDSHAKE_CHANNEL,
			"version":                  protocol.BAYEUX_VERSION,
			"advice":                   protocol.DEFAULT_ADVICE,
			"supportedConnectionTypes": []string{"websocket", "long-polling"},
			"successful":               true,
		})

	} else {
		response["error"] = fmt.Sprintf("Only supported version is '%s'", protocol.BAYEUX_VERSION)
	}

	// Answer directly
	conn.Send([]protocol.Message{response})
	return newClientId
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

func (m Engine) SubscribeService(chanOut chan<- protocol.Message, subscription []string) {
	conn := transport.InternalConnection{chanOut}
	newClient := m.NewClient(conn)
	newClient.Connect(-1, 0, nil, conn)
	newClient.SetConnection(conn)
	m.AddSubscription(newClient.Id(), subscription)
}

func (m Engine) SubscribeClient(request protocol.Message, client *protocol.Client) {
	response := m.responseFromRequest(request)
	response["successful"] = true

	subscription := request["subscription"]
	response["subscription"] = subscription

	var subs []string
	switch subscription.(type) {
	case []string:
		subs = subscription.([]string)
	case string:
		subs = []string{subscription.(string)}
	}

	for _, s := range subs {
		// Do not register clients subscribing to a service channel
		// They will be answered directly instead of through the normal subscription system
		if !protocol.NewChannel(s).IsService() {
			m.AddSubscription(client.Id(), []string{s})
		}
	}

	client.Queue(response)
}

func (m Engine) Disconnect(request protocol.Message, client *protocol.Client, conn protocol.Connection) {
	response := m.responseFromRequest(request)
	response["successful"] = true
	clientId := request.ClientId()
	m.logger.Debugf("Client %s disconnected", clientId)
}

func (m Engine) Publish(request protocol.Message) {
	requestingClient := m.clients.GetClient(request.ClientId())

	if requestingClient == nil {
		m.logger.Warnf("PUBLISH from unknown client %s", request)
	} else {
		response := m.responseFromRequest(request)
		response["successful"] = true
		data := request["data"]
		channel := request.Channel()

		m.clients.GetClient(request.ClientId()).Queue(response)

		go func() {
			// Prepare msg to send to subscribers
			msg := protocol.Message{}
			msg["channel"] = channel.Name()
			msg["data"] = data
			// TODO: Missing ID

			msg.SetClientId(request.ClientId())

			// Get clients with subscriptions
			recipients := m.clients.GetClients(channel.Expand())
			m.logger.Debugf("PUBLISH from %s on %s to %d recipients", request.ClientId(), channel, len(recipients))
			// Queue messages
			for _, c := range recipients {
				m.clients.GetClient(c).Queue(msg)
			}
		}()

	}

}

// Publish message directly to client
// msg should have "channel" which the client is expecting, e.g. "/service/echo"
func (m Engine) PublishFromService(recipientId string, msg protocol.Message) {
	// response["successful"] = true
	m.clients.GetClient(recipientId).Queue(msg)
}
