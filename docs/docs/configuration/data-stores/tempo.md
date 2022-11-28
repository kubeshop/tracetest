# Tempo

```yaml
telemetry:
  dataStores:
    tempo:
      type: tempo
      tempo:
        endpoint: tempo:9095
        tls:
          insecure: true

server:
    telemetry:
        dataStore: tempo
```
