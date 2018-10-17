# forwardif

[coredns][] plugin for conditional forward rules.

[![Build Status](https://travis-ci.org/b4fun/forwardif.svg)](https://travis-ci.org/b4fun/forwardif)
[![](https://godoc.org/github.com/b4fun/forwardif?status.svg)](http://godoc.org/github.com/b4fun/forwardif)

[coredns]: https://coredns.io/

## Example

```
. {
  forwardif {
    adb_file adb:///path/to/file

    // forward rule on one of the conditions matched
    forward_to 10.0.0.1
  }

  // fallback rule
  forward . 8.8.8.8
}
```

## coredns build

See [coredns_build](https://github.com/b4fun/coredns_build).

## LICENSE

MIT
