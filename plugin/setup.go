package plugin

import (
	"errors"

	"github.com/b4fun/forwardif"
	"github.com/b4fun/forwardif/adb"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/forward"
	"github.com/coredns/coredns/plugin/pkg/parse"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyfile"
)

var (
	errForwardRuleRequired = errors.New("forward rule required")
	errForwardToRequired   = errors.New("forward_to required")
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
		return f, errForwardRuleRequired
	}

	return f, nil
}

func parseBlock(c *caddyfile.Dispenser, f *ForwardIf) (err error) {
	var forwardOpt *forward.Opt
	getForwardOpt := func() *forward.Opt {
		if forwardOpt == nil {
			forwardOpt = forward.NewDefault()
		}

		return forwardOpt
	}

	switch c.Val() {
	case "forward_from":
		if !c.NextArg() {
			return c.ArgErr()
		}
		getForwardOpt().From = c.Val()
	case "forward_to":
		to := c.RemainingArgs()
		if len(to) < 1 {
			return c.ArgErr()
		}

		getForwardOpt().ToHosts, err = parse.HostPortOrFile(to...)
		if err != nil {
			return err
		}
	case "pattern":
		if !c.NextArg() {
			return c.ArgErr()
		}

		f.matcher, err = forwardif.NewRegexMatch(c.Val())
		if err != nil {
			return err
		}
	case "adb":
		if !c.NextArg() {
			return c.ArgErr()
		}

		adb, err := adb.NewAdBlockFromFilePath(c.Val(), false)
		if err != nil {
			return err
		}
		f.matcher = adb.Match
	case "adb_exception":
		if !c.NextArg() {
			return c.ArgErr()
		}

		adb, err := adb.NewAdBlockFromFilePath(c.Val(), true)
		if err != nil {
			return err
		}
		f.matcher = adb.Match
	}

	if forwardOpt != nil {
		if len(forwardOpt.ToHosts) < 1 {
			return errForwardToRequired
		}

		f.Forward = forwardOpt.Build()
	}

	return nil
}
