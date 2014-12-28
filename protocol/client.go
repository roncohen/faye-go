package protocol

import (
	"github.com/roncohen/faye/utils"
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
	started  time.Time
	logger   utils.Logger
}

func NewSession(client *Client, conn Connection, timeout int, response Message, logger utils.Logger) *Session {
	session := Session{conn, timeout, response, client, time.Now(), logger}
	if timeout > 0 {
		go func() {
			time.Sleep(time.Duration(timeout) * time.Millisecond)
			session.End()
		}()
	}
	return &session
}

func (s Session) End() {
	// c.mutex.Lock()
	// defer c.mutex.Unlock()

	if s.conn.IsConnected() {
		s.conn.Send([]Message{s.response})
	} else {
		s.logger.Debugf("No longer connected %s", s.client.clientId)
	}
}

type Client struct {
	clientId    string
	connection  Connection
	msgStore    MsgStore
	isConnected bool
	responseMsg Message
	mutex       sync.Mutex
	lastSession *Session
	created     time.Time
	logger      utils.Logger
}

func NewClient(clientId string, msgStore MsgStore, logger utils.Logger) Client {
	client := Client{
		clientId:    clientId,
		msgStore:    msgStore,
		isConnected: false,
		created:     time.Now(),
		logger:      logger,
	}

	return client
}

func (c Client) Id() string {
	return c.clientId
}

func (c *Client) Connect(timeout int, interval int, responseMsg Message, connection Connection) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.lastSession = NewSession(c, connection, timeout, responseMsg, c.logger)
	c.responseMsg = responseMsg

	c.flushMsgs()
}

func (c *Client) SetConnection(connection Connection) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.connection == nil || connection.Priority() > c.connection.Priority() {
		c.connection = connection
		c.isConnected = true
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

func (c Client) IsExpired() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if time.Now().Sub(c.created) > time.Duration(1*time.Minute) {
		if c.lastSession != nil &&
			time.Now().Sub(c.lastSession.started) > time.Duration(2*time.Hour) {
			return true
		}
	}
	return false
}

func (c Client) flushMsgs() {
	if c.isConnected && c.connection != nil && c.connection.IsConnected() {
		msgs := c.msgStore.GetAndClearMessages()
		if len(msgs) > 0 {

			var msgsWithConnect []Message
			if c.connection.IsSingleShot() {
				msgsWithConnect = append(msgs, c.responseMsg)

			} else {
				msgsWithConnect = msgs
			}
			c.logger.Debugf("Sending %d msgs to %s on %s", len(msgsWithConnect), c.clientId, reflect.TypeOf(c.connection))

			err := c.connection.Send(msgsWithConnect)

			// failed, so requeue
			if err != nil {
				c.logger.Debugf("Was unable to send to %s requeued %d messages", c.clientId, len(msgs))
				c.msgStore.EnqueueMessages(msgs)
			} else {
				c.responseMsg = nil
				c.isConnected = false
			}
		}
	} else {
		c.logger.Debugf("Not connected for %s", c.clientId)
	}
}
