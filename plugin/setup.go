package plugin

import (
	"errors"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/forward"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyfile"
)

var (
	errForwardRuleRequired = errors.New("forward rule required")
)

func init() {
	caddy.RegisterPlugin("forwardif", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	f, err := parseForwardIf(c)
	if err != nil {
		return plugin.Error("forwardif", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		f.Next = next
		f.Forward.Next = next
		return f
	})

	c.OnStartup(func() error {
		return f.OnStartup()
	})

	c.OnStartup(func() error {
		return f.OnShutdown()
	})

	return nil
}

func (f *ForwardIf) OnStartup() error {
	return f.Forward.OnStartup()
}

func (f *ForwardIf) OnShutdown() error {
	return f.Forward.OnShutdown()
}

func parseForwardIf(c *caddy.Controller) (*ForwardIf, error) {
	if !c.Next() {
		return nil, errors.New("this plugin should have body")
	}

	f, err := ParseForwardIfStanza(&c.Dispenser)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func ParseForwardIfStanza(c *caddyfile.Dispenser) (*ForwardIf, error) {
	f := New()

	for c.NextBlock() {
		if err := parseBlock(c, f); err != nil {
			return f, err
		}
	}

	if f.Forward == nil {
		return f, plugin.Error("forwardif", errForwardRuleRequired)
	}

	return f, nil
}

func parseBlock(c *caddyfile.Dispenser, f *ForwardIf) error {
	switch c.Val() {
	case "forward":
		forward, err := forward.ParseForwardStanza(c)
		if err != nil {
			return err
		}
		f.Forward = forward
	}

	return nil
}
