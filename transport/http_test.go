package transport

import (
	"github.com/roncohen/faye/protocol"
	// "code.google.com/p/go-stringio"
	// "net/http/httptest"
	"testing"
)

func TestIsConnectedShouldReturnFalseOn2ndTry(t *testing.T) {
	// responseRecorder := httptest.NewRecorder()
	conn := NewLongPollingConnection()

	if !conn.IsConnected() {
		t.Fatal("new LongPollingConnection.IsConnected should be true")
	}

	conn.Send([]protocol.Message{protocol.Message{"channel": "/meta/connect"}})

	if conn.IsConnected() {
		t.Fatal("LongPollingConnection.IsConnected should be false after send")
	}
}
