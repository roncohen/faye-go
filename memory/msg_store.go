package memory

import (
	"container/list"
	"github.com/roncohen/faye-go/protocol"
)

type MemoryMsgStore struct {
	msgs *list.List
}

func NewMemoryMsgStore() *MemoryMsgStore {
	return &MemoryMsgStore{list.New()}
}

func (m *MemoryMsgStore) EnqueueMessages(msgs []protocol.Message) {
	for _, msg := range msgs {
		m.msgs.PushBack(msg)
	}
}

func (m *MemoryMsgStore) GetAndClearMessages() []protocol.Message {
	var msgArray = make([]protocol.Message, m.msgs.Len())
	i := 0
	for e := m.msgs.Front(); e != nil; e = e.Next() {
		msgArray[i] = e.Value.(protocol.Message)
		i = i + 1
	}
	m.msgs = &list.List{}
	return msgArray
}
