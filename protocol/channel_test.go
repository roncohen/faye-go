package protocol_test

import (
	. "github.com/roncohen/faye-go/protocol"
	"github.com/roncohen/faye-go/utils"
	"testing"
)

func TestExpandSimpleChannel(t *testing.T) {
	chan1 := NewChannel("/foo/bar")

	expected := []string{
		"/**",
		"/foo/**",
		"/foo/*",
		"/foo/bar",
	}

	patterns := chan1.Expand()

	if !utils.CompareStringSlices(expected, patterns) {
		t.Fatal("Expected ", expected, " got ", patterns)
	}
}
