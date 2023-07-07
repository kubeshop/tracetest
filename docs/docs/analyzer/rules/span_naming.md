# span_naming

Enforce span names that identify a class of Spans, rather than individual Span instances

## Rule Details

The span name concisely identifies the work represented by the Span, for example, an RPC method name, a function name, or the name of a subtask or stage within a larger computation. The span name SHOULD be the most general string that identifies a class of Spans, rather than individual Span instances while still being human-readable.

The following OTel semantic conventions for span names are defined:

### HTTP spans:

If span kind is `"server"`, the name should follow this format:

```
{http.method} {http.route}
```

If span kind is `"client"`, the name should follow this format:

```
{http.method}
```

### Database spans:

```
{db.operation} {db.name}.{db.sql.table}
```

If `db.sql.table` is not available, the name should follow this format:

```
{db.operation} {db.name}
```

### RPC spans:

```
{package}.{service}/{method}
```

### Messaging spans:

```
{destination name} {operation name}
```

## Options

This rule has the following options:

- `"error"` requires span names to follow the OTel semantic convention
- `"disabled"` disables the span name verification
- `"warning"` verifies span names to follow the OTel semantic convention but does not impact the analyzer score

## When Not To Use It

If you don’t want to enforce OTel span names, don’t enable this rule.
