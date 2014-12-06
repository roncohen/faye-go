package protocol_test

import (
	"github.com/roncohen/faye/memory"
	. "github.com/roncohen/faye/protocol"
	"log"
	"testing"
)

var test_ns = memory.NewMemoryNamespace()

func getNewClient() Client {
	msgStore := memory.NewMemoryMsgStore()
	return NewClient(test_ns.Generate(), msgStore)
}

func getMsgWithId(responses []Message, msgId string) Message {
	for _, m := range responses {
		if id, ok := m["id"]; ok && id == msgId {
			return m
		}
	}
	return nil
}

type FakeSingleShotConnection struct {
	SentMessages chan []Message
	closed       bool
}

func (f *FakeSingleShotConnection) Send(msgs []Message) error {
	f.SentMessages <- msgs
	log.Printf("Fakeconnection got msgs: %+v", msgs)
	return nil
}

func (f FakeSingleShotConnection) Close() {
	f.closed = true
}

func (f FakeSingleShotConnection) IsConnected() bool {
	return !f.closed
}

func (f FakeSingleShotConnection) IsSingleShot() bool {
	return true
}

func (f FakeSingleShotConnection) Priority() int {
	return 1
}

func TestEnqueAndReleaseMsgs(t *testing.T) {
	conn := FakeSingleShotConnection{closed: false, SentMessages: make(chan []Message, 10)}
	client := getNewClient()
	client.Queue(Message{"channel": "/meta/subscribe", "id": "1"})
	client.SetConnection(&conn)
	client.Connect(100, 0, Message{"channel": "/meta/connect", "id": "2"}, &conn)

	sentMessages := <-conn.SentMessages

	if len(sentMessages) != 2 {
		t.Fatal("Should release 2 messages, got:", sentMessages)
	}

	if sentMessages[0]["channel"] != "/meta/subscribe" {
		t.Fatal("First sent message should be subscribe response")
	}

	if sentMessages[1]["channel"] != "/meta/connect" {
		t.Fatal("Second sent message should be connect response, response", sentMessages[1])
	}
}

// func TestTimesoutsMsgs(t *testing.T) {
// 	conn := FakeSingleShotConnection{closed: false, SentMessages: make(chan []Message, 10)}
// 	client := getNewClient()

// 	client.Connect(100, 0, Message{"channel": "/meta/connect", "id": "2"}, &conn)
// 	client.SetConnection(&conn)

// 	sentMessages := <-conn.SentMessages

// 	if len(sentMessages) != 1 {
// 		t.Fatal("Should release 1 messages, got:", sentMessages)
// 	}

// 	if sentMessages[0]["channel"] != "/meta/connect" {
// 		t.Fatal("Message should be connect response")
// 	}

// 	if sentMessages[0]["id"] != "2" {
// 		t.Fatal("Connect response has invalid id", sentMessages[0])
// 	}
// }

/*
func TestTwoConnectionsAndTimeouts(t *testing.T) {
	conn1 := FakeConnection{closed: false, SentMessages: make(chan []Message, 10)}
	conn2 := FakeConnection{closed: false, SentMessages: make(chan []Message, 10)}
	client := getNewClient()
	client.SetConnection(0, 0, Message{"channel": "/meta/connect", "id": "2"}, &conn1)
	client.SetConnection(0, 0, Message{"channel": "/meta/connect", "id": "3"}, &conn2)

	sentMessages1 := <-conn1.SentMessages
	sentMessages2 := <-conn2.SentMessages

	if sentMessages1[0]["id"] != "2" {
		t.Fatal("Connect response expected response with id: 2, msg was:", sentMessages1[0])
	}

	if sentMessages2[0]["id"] != "3" {
		t.Fatal("Connect response expected response with id: 3, msg was:", sentMessages2[0])
	}
}

func TestTwoConnectionsAndNewMessage(t *testing.T) {
	conn1 := FakeConnection{closed: false, SentMessages: make(chan []Message, 10)}
	conn2 := FakeConnection{closed: false, SentMessages: make(chan []Message, 10)}
	client := getNewClient()
	client.SetConnection(1000, 0, Message{"channel": "/meta/connect", "id": "2"}, &conn1)
	client.SetConnection(1000, 0, Message{"channel": "/meta/connect", "id": "3"}, &conn2)

	// sends on first conn
	client.Queue(Message{"channel": "/meta/subscribe", "id": "4"})

	// sends on next conn
	client.Queue(Message{"channel": "/meta/subscribe", "id": "5"})

	sentMessages1 := <-conn1.SentMessages
	sentMessages2 := <-conn2.SentMessages

	if sentMessages1[0]["id"] != "2" {
		t.Fatal("Connect response expected response with id: 2, msg was:", sentMessages1[0])
	}

	if sentMessages2[0]["id"] != "3" {
		t.Fatal("Connect response expected response with id: 3, msg was:", sentMessages2[0])
	}
}
*/
