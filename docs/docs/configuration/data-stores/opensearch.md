# OpenSearch

> :information_source: the parameter `index` is the name of the OpenSearch index where your traces are located. Usually you set it in the Data Prepper configuration.

```yaml
telemetry:
  dataStores:
    opensearch:
      type: opensearch
      opensearch:
        addresses:
          - http://opensearch:9200
        index: traces

server:
    telemetry:
        dataStore: opensearch
```
