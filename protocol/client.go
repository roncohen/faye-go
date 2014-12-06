package protocol

import (
	"log"
	"reflect"
	"sync"
	"time"
)

type MsgStore interface {
	EnqueueMessages([]Message)
	GetAndClearMessages() []Message
}

// Connect requests starts a session
type Session struct {
	conn     Connection
	timeout  int
	response Message
	client   *Client
}

func NewSession(client *Client, conn Connection, timeout int, response Message) Session {
	session := Session{conn, timeout, response, client}
	log.Println("NewSession with timout", timeout)
	go func() {
		time.Sleep(time.Duration(timeout) * time.Millisecond)
		session.End()
	}()
	return session
}

func (s Session) End() {
	// c.mutex.Lock()
	// defer c.mutex.Unlock()

	if s.conn.IsConnected() {
		s.conn.Send([]Message{s.response})
	} else {
		log.Println("No longer connected ", s.client.clientId, reflect.TypeOf(s.conn))
	}
}

type Client struct {
	clientId     string
	connection   Connection
	msgStore     MsgStore
	is_connected bool
	responseMsg  Message
	mutex        sync.Mutex
}

func NewClient(clientId string, msgStore MsgStore) Client {
	client := Client{
		clientId:     clientId,
		msgStore:     msgStore,
		is_connected: false,
	}

	return client
}

func (c Client) Id() string {
	return c.clientId
}

func (c *Client) Connect(timeout int, interval int, responseMsg Message, connection Connection) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	NewSession(c, connection, timeout, responseMsg)

	c.responseMsg = responseMsg

	c.flushMsgs()
}

func (c *Client) SetConnection(connection Connection) {
	log.Println("Setting connection for", c.clientId, reflect.TypeOf(connection))
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.connection == nil || connection.Priority() > c.connection.Priority() {
		c.connection = connection
		c.is_connected = true
	}
}

func (c Client) Queue(msg Message) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.msgStore.EnqueueMessages([]Message{msg})
	c.flushMsgs()
}

func (c Client) QueueMany(msgs []Message) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.msgStore.EnqueueMessages(msgs)
	c.flushMsgs()
}

func (c Client) flushMsgs() {
	if c.is_connected && c.connection != nil && c.connection.IsConnected() {
		msgs := c.msgStore.GetAndClearMessages()
		if len(msgs) > 0 {

			var msgsWithConnect []Message
			if c.connection.IsSingleShot() {
				msgsWithConnect = append(msgs, c.responseMsg)

			} else {
				msgsWithConnect = msgs
			}
			log.Print("Sending ", len(msgsWithConnect), " msgs to ", c.clientId, " on ", reflect.TypeOf(c.connection))
			err := c.connection.Send(msgsWithConnect)

			// failed, so requeue
			if err != nil {
				log.Print("Was unable to send to ", c.clientId, ", requeued ", len(msgs), " messages")
				c.msgStore.EnqueueMessages(msgs)
			} else {
				c.responseMsg = nil
				c.is_connected = false
			}
		}
	} else {
		log.Print("Not connected for ", c.clientId)
	}
}
