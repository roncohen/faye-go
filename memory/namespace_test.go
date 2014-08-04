package memory

import "testing"

func TestGenerateShouldReturn32Chars(t *testing.T) {
	ns := NewMemoryNamespace()
	newId := ns.Generate()

	if len(newId) != 32 {
		t.Error("New ids should be 32 chars")
	}

	if !ns.IsUsed(newId) {
		t.Error("IsUsed on new ID should return true")
	}
}

func TestExpireNonExisting(t *testing.T) {
	ns := NewMemoryNamespace()

	id := "abcabc"

	ns.Expire(id)

	if ns.IsUsed(id) {
		t.Error("IsUsed should return false")
	}
}

func TestExpire(t *testing.T) {
	ns := NewMemoryNamespace()

	newId := ns.Generate()
	if !ns.IsUsed(newId) {
		t.Error("IsUsed on new ID should return true")
	}

	ns.Expire(newId)

	if ns.IsUsed(newId) {
		t.Error("IsUsed should return false")
	}
}
