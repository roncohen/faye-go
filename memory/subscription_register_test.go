package memory

import (
	. "github.com/roncohen/faye-go/utils"

	"testing"
)

func TestAddSubscription(t *testing.T) {
	subreg := NewSubscriptionRegister()
	clientId := "clientId"
	patterns := "/foo/**"
	subreg.AddSubscription(clientId, []string{patterns})

	expected_clients := []string{clientId}
	actual_clients := subreg.GetClients([]string{patterns})

	if !CompareStringSlices(actual_clients, expected_clients) {
		t.Error("Expected client differ from actual clients. Expected: ", expected_clients, " actual: ", actual_clients)
	}
}

func TestRemoveSubscription(t *testing.T) {
	subreg := NewSubscriptionRegister()
	clientId := "clientId"
	patterns := "/foo/**"
	subreg.AddSubscription(clientId, []string{patterns})
	subreg.RemoveSubscription(clientId, []string{patterns})

	expected_clients := []string{}
	actual_clients := subreg.GetClients([]string{patterns})

	if !CompareStringSlices(actual_clients, expected_clients) {
		t.Error("Expected client differ from actual clients. Expected: ", expected_clients, " actual: ", actual_clients)
	}
}

func TestRemoveClient(t *testing.T) {
	subreg := NewSubscriptionRegister()
	clientId := "clientId"
	patterns := "/foo/**"
	subreg.AddSubscription(clientId, []string{patterns})
	// subreg.RemoveSubscription(clientId, []string{patterns})
	subreg.RemoveClient(clientId)

	expected_clients := []string{}
	actual_clients := subreg.GetClients([]string{patterns})

	if !CompareStringSlices(actual_clients, expected_clients) {
		t.Error("Expected client differ from actual clients. Expected: ", expected_clients, " actual: ", actual_clients)
	}
}

// func TestExpireNonExisting(t *testing.T) {
// 	ns := NewMemoryNamespace()

// 	id := "abcabc"

// 	ns.Expire(id)

// 	if ns.IsUsed(id) {
// 		t.Error("IsUsed should return false")
// 	}
// }

// func TestExpire(t *testing.T) {
// 	ns := NewMemoryNamespace()

// 	newId := ns.Generate()
// 	if !ns.IsUsed(newId) {
// 		t.Error("IsUsed on new ID should return true")
// 	}

// 	ns.Expire(newId)

// 	if ns.IsUsed(newId) {
// 		t.Error("IsUsed should return false")
// 	}
// }
