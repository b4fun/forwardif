package plugin

import (
	"context"

	"github.com/b4fun/forwardif"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/forward"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("forwardif")

type ForwardIf struct {
	Forward *forward.Forward

	matcher forwardif.MatchFunc

	Next plugin.Handler
}

func New() *ForwardIf {
	return &ForwardIf{
		matcher: forwardif.FalseMatcher,
	}
}

func (f ForwardIf) Name() string { return "forwardif" }

func (f *ForwardIf) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	var next plugin.Handler
	if f.matcher(state) {
		next = f.Forward
		log.Infof("matched for %s", state.Name())
	} else {
		next = f.Next
		log.Infof("missed for %s", state.Name())
	}

	return plugin.NextOrFailure(f.Name(), next, ctx, w, r)
}
