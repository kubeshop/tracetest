# no-api-key-leak

Disallow leaked API keys for HTTP spans.

## Rule Details

This rule disallows the recording of API keys in HTTP spans.

### HTTP spans:

The following attributes are evaluated:

```
- http.response.header.authorization
- http.response.header.x-api-key
- http.request.header.authorization
- http.request.header.x-api-key
```

## Options

This rule has the following options:

- `"error"` requires no leaked API keys for HTTP spans
- `"disabled"` disables the no leaked API keys verification for HTTP spans
- `"warning"` verifies no leaked API keys for HTTPS spans but does not impact the analyzer score

## When Not To Use It

If you intentionally record API keys for HTTP spans then you can disable this rule.
