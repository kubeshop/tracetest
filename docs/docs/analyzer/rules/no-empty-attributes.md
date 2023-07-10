# no-empty-attributes

Disallow empty attribute values.

## Rule Details

An `Attribute` is a key-value pair, which is encapsulated as part of a span. The attribute value should not be empty to be considered valid.

## Options

This rule has the following options:

- `"error"` requires attribute values to not be empty
- `"disabled"` disables the attribute values verification
- `"warning"` verifies attribute values to not be empty but does not impact the analyzer score

## When Not To Use It

If you intentionally use empty attribute values then you can disable this rule.
