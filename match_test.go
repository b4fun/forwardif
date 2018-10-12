package forwardif

import (
	"testing"

	"github.com/coredns/coredns/request"
)

func TestFalseMatcher(t *testing.T) {
	state := request.Request{}

	if m := FalseMatcher(state); m != false {
		t.Errorf("unexpected: %v", m)
	}
}
