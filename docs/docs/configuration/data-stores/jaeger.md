# Jaeger

```yaml
telemetry:
  dataStores:
    my_jaeger_instance_name:
      type: jaeger
      jaeger:
        endpoint: url-to-jaeger:16685
        tls:
          insecure: true

server:
    telemetry:
        dataStore: my_jaeger_instance_name
```
