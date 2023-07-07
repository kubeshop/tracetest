# attribute_naming

Enforce attribute keys to follow common specifications

## Rule Details

An Attribute is a key-value pair, which is encapsulated as part of a span. The attribute key should follow a set of common specifications to be considered valid.

The following OTel semantic conventions for attribute keys are defined:

- It must be a non-null and non-empty string.
- It must be a valid Unicode sequence.
- It should use namespacing to avoid name clashes. Delimit the namespaces using a dot character. For example `service.version` denotes the service version where `service` is the namespace and `version` is an attribute in that namespace.
- Namespaces can be nested. For example `telemetry.sdk` is a namespace inside top-level `telemetry` namespace and `telemetry.sdk.name` is an attribute inside `telemetry.sdk` namespace.
- For each multi-word separate the words by underscores (use snake_case). For example `http.status_code` denotes the status code in the http namespace.
- Names should not coincide with namespaces. For example if `service.instance.id` is an attribute name then it is no longer valid to have an attribute named `service.instance` because `service.instance` is already a namespace.

## Options

This rule has the following options:

- `"error"` requires attribute keys to follow the OTel semantic convention
- `"disabled"` disables the attribute keys verification
- `"warning"` verifies attribute keys to follow the OTel semantic convention but does not impact the analyzer score

## When Not To Use It

If you don’t want to enforce OTel attribute keys, don’t enable this rule.
