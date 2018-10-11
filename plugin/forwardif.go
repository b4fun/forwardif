package plugin

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/forward"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type ForwardIf struct {
	Forward *forward.Forward

	Next plugin.Handler
}

func New() *ForwardIf {
	return &ForwardIf{}
}

func (f ForwardIf) Name() string { return "forwardif" }

func (f *ForwardIf) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	var next plugin.Handler
	if f.match(state) {
		next = f.Forward
	} else {
		next = f.Next
	}

	return plugin.NextOrFailure(f.Name(), next, ctx, w, r)
}

func (f ForwardIf) match(state request.Request) bool {
	return true
}
