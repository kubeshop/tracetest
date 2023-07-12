# secure-https-protocol

Enforce usage of secure protocol for HTTP server spans.

## Rule Details

This rule enforces usage of a secure protocol for an HTTP server span. The URI scheme that identifies the used protocol should be `"https"`.

### HTTP spans:

If span kind is `"server"`, the following attributes are evaluated:

```
- http.scheme = "https"
- http.url = "https"
```

## Options

This rule has the following options:

- `"error"` requires secure protocol for HTTP server spans
- `"disabled"` disables the secure protocol verification for HTTP server spans
- `"warning"` verifies secure protocol for HTTPS server spans but does not impact the analyzer score

## When Not To Use It

If you intentionally use non secure protocol for HTTP server spans then you can disable this rule.
