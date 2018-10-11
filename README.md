# forwardif

[coredns][] plugin for conditional forward rules.

[coredns]: https://coredns.io/

## Example

```
. {
  forwardif {
    adb_file adb:///path/to/file

    // forward rule on one of the conditions matched
    forward . 10.0.0.1
  }

  // fallback rule
  forward . 8.8.8.8
}
```
