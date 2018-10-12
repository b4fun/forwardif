package forwardif

import (
	"testing"

	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func TestRegexMatch_Basic(t *testing.T) {
	makeState := func(question string) request.Request {
		state := request.Request{W: &test.ResponseWriter{}, Req: new(dns.Msg)}
		state.Req.SetQuestion(question, dns.TypeA)
		return state
	}

	cases := []struct {
		regex   string
		state   request.Request
		matched bool
	}{
		{
			"example",
			makeState("example.org."),
			true,
		},
		{
			"not_example",
			makeState("example.org."),
			false,
		},
		{
			"\\.org",
			makeState("example.org."),
			true,
		},
		{
			"\\.io",
			makeState("example.org."),
			false,
		},
	}
	for _, c := range cases {
		m, err := NewRegexMatch(c.regex)
		if err != nil {
			t.Error(err)
			continue
		}
		if actual := m(c.state); actual != c.matched {
			t.Errorf("expected = %v, actual = %v", c.matched, actual)
			continue
		}
	}
}
