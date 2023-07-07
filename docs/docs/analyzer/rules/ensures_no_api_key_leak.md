# ensures_no_api_key_leak

Disallow leaked api keys for HTTP spans

## Rule Details

This rule disallows the record of api keys in HTTP spans.

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

- `"error"` requires not leaked api keys for HTTP spans
- `"disabled"` disables the no leaked api keys verification for HTTP spans
- `"warning"` verifies not leaked api keys for HTTPS spans but does not impact the analyzer score

## When Not To Use It

If you intentionally record api keys for HTTP spans then you can disable this rule.
