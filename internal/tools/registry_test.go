package tools

import (
	"testing"
)

func TestNewDefaultRegistry_HasAllDefaultTools(t *testing.T) {
	t.Parallel()

	r := NewDefaultRegistry()

	for _, name := range DefaultToolNames() {
		if _, ok := r.Get(name); !ok {
			t.Fatalf("expected tool %q to be registered", name)
		}
	}
}

func TestMustGet_Unknown(t *testing.T) {
	t.Parallel()

	r := NewRegistry()
	if _, err := r.MustGet("unknown"); err != ErrToolNotFound {
		t.Fatalf("expected ErrToolNotFound, got %v", err)
	}
}
