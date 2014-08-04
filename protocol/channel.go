package protocol

import (
	"strings"
)

type MetaChannel interface{}

const (
	META_PREFIX             string = "/meta/"
	META_SERVICE                   = "/service"
	META_HANDSHAKE_CHANNEL         = "handshake"
	META_SUBSCRIBE_CHANNEL         = "subscribe"
	META_CONNECT_CHANNEL           = "connect"
	META_DISCONNECT_CHANNEL        = "disconnect"
	META_UNKNOWN_CHANNEL           = "unknown"
)

func NewChannel(name string) Channel {
	return Channel{name}
}

type Channel struct {
	name string
}

type Subscription Channel

func (c Channel) Name() string {
	return c.name
}

func (c Channel) IsMeta() bool {
	return strings.HasPrefix(c.name, META_PREFIX)
}

func (c Channel) IsService() bool {
	return strings.HasPrefix(c.name, META_SERVICE)
}

func (c Channel) MetaType() MetaChannel {
	if !c.IsMeta() {
		return nil
	} else {
		switch c.name[len(META_PREFIX):] {
		case META_CONNECT_CHANNEL:
			return META_CONNECT_CHANNEL
		case META_SUBSCRIBE_CHANNEL:
			return META_SUBSCRIBE_CHANNEL
		case META_DISCONNECT_CHANNEL:
			return META_DISCONNECT_CHANNEL
		case META_HANDSHAKE_CHANNEL:
			return META_HANDSHAKE_CHANNEL
		default:
			return META_UNKNOWN_CHANNEL
		}
	}
}

// Returns all the channels patterns that could match this channel
/*
For:
/foo/bar
We should return these:
/**
/foo/**
/foo/*
/foo/bar
*/
func (c Channel) Expand() []string {
	segments := strings.Split(c.name, "/")
	num_segments := len(segments)
	patterns := make([]string, num_segments+1)
	patterns[0] = "/**"
	for i := 1; i < len(segments); i = i + 2 {
		patterns[i] = strings.Join(segments[:i+1], "/") + "/**"
	}
	patterns[len(patterns)-2] = strings.Join(segments[:num_segments-1], "/") + "/*"
	patterns[len(patterns)-1] = c.name
	return patterns
}
