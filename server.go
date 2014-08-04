package faye

import (
	"encoding/json"
	"github.com/roncohen/faye/protocol"
	"io"
	"log"
)

type Server struct {
	engine Engine
}

func NewServer(engine Engine) Server {
	return Server{engine}
}

func (s Server) GetClient(request protocol.Message, conn protocol.Connection) *protocol.Client {
	clientId := request.ClientId()
	client := s.engine.GetClient(clientId)
	if client == nil {
		log.Printf("Message %v from unknown client %v", request.Channel(), clientId)
		response := request
		response["successful"] = false
		response["advice"] = map[string]interface{}{"reconnect": "handshake", "interval": "1000"}
		conn.Send([]protocol.Message{response})
		conn.Close()
	}
	return client
}

func (s Server) HandleRequest(msges interface{}, conn protocol.Connection) {
	switch msges.(type) {
	case []interface{}:
		msg_list := msges.([]interface{})
		for _, msg := range msg_list {
			s.HandleMessage(msg.(map[string]interface{}), conn)
		}
	case map[string]interface{}:
		s.HandleMessage(msges.(map[string]interface{}), conn)
	}
}

func (s Server) HandleMessage(msg protocol.Message, conn protocol.Connection) {
	channel := msg.Channel()
	if channel.IsMeta() {
		s.HandleMeta(msg, conn)
	} else if channel.IsService() {
		s.HandleService(msg)
	} else {
		s.HandlePublish(msg)
	}
}

func (s Server) HandleService(msg protocol.Message) protocol.Message {
	return nil
}

func (s Server) HandleMeta(msg protocol.Message, conn protocol.Connection) protocol.Message {
	meta_channel := msg.Channel().MetaType()

	if meta_channel == protocol.META_HANDSHAKE_CHANNEL {
		s.engine.Handshake(msg, conn)
	} else {
		client := s.GetClient(msg, conn)
		if client != nil {
			client.SetConnection(conn)

			switch meta_channel {
			case protocol.META_HANDSHAKE_CHANNEL:
				s.engine.Handshake(msg, conn)
			case protocol.META_CONNECT_CHANNEL:
				s.engine.Connect(msg, client, conn)

			case protocol.META_DISCONNECT_CHANNEL:
				s.engine.Disconnect(msg, client, conn)

			case protocol.META_SUBSCRIBE_CHANNEL:
				s.engine.Subscribe(msg, client)

			case protocol.META_UNKNOWN_CHANNEL:
				log.Panicf("Message with unknown meta channel received")

			}
		}
	}

	return nil
}

func (s Server) HandlePublish(msg protocol.Message) {
	// Publish
	clientId := msg.ClientId()
	// if _client, is_connected := s.engine.GetClient(clientId); !is_connected {
	// 	// TODO: Howto answer if not connected.
	// 	return nil
	// }

	log.Printf("Client %s publishing to %s", clientId, msg.Channel())
	s.engine.Publish(msg)
}

func JSONWrite(w io.Writer, obj interface{}) error {
	msg, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	w.Write(msg)
	return nil
}
