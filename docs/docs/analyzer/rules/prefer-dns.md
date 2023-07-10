# prefer-dns

Enforce usage of DNS instead of IP addresses.

## Rule Details

When connecting to remote servers, ensure the usage of DNS instead of IP addresses to avoid issues.

The following attributes are evaluated:

```
- http.url
- db.connection_string
```

If span kind is `"client"`, the following attributes are evaluated:

```
- net.peer.name
```

## Options

This rule has the following options:

- `"error"` requires DNS over IP addresses
- `"disabled"` disables the DNS over IP addresses verification
- `"warning"` verifies DNS over IP addresses but does not impact the analyzer score

## When Not To Use It

If you intentionally use and record IP addresses then you can disable this rule.
