package memory

import (
	// "container/list"
	// "github.com/roncohen/faye/engines"
	"github.com/roncohen/faye/protocol"
	// "log"
	"testing"
)

func TestEnqueAndGetBack(t *testing.T) {
	msgstore := NewMemoryMsgStore()
	msg := protocol.Message{"mymsg": "bleh"}

	msgstore.EnqueueMessages([]protocol.Message{msg})
	msgs := msgstore.GetAndClearMessages()

	if len(msgs) != 1 {
		t.Fatal("Should get one msgs back")
	}

	if msgs[0]["mymsg"] != msg["mymsg"] {
		t.Fatal("Got ", msgs[0], " expected ", msg)
	}
}
