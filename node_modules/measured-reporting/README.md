# Measured Reporting

The registry and reporting library that has the classes needed to create a dimension aware, self reporting metrics registry.

[![npm](https://img.shields.io/npm/v/measured-reporting.svg)](https://www.npmjs.com/package/measured-reporting) 

## Install

```
npm install measured-reporting
```

## What is in this package

### [Self Reporting Metrics Registry](https://yaorg.github.io/node-measured/SelfReportingMetricsRegistry.html)
A dimensional aware self-reporting metrics registry, just supply this class with a reporter implementation at instantiation and this is all you need to instrument application level metrics in your app.

See the [SelfReportingMetricsRegistryOptions](http://yaorg.github.io/node-measured/build/docs/packages/measured-reporting/global.html#SelfReportingMetricsRegistryOptions) for advanced configuration.

```javascript
const { SelfReportingMetricsRegistry, LoggingReporter } = require('measured-reporting');
const registry = new SelfReportingMetricsRegistry(new LoggingReporter({
  defaultDimensions: {
    hostname: os.hostname()
  }
}));

// The metric will flow through LoggingReporter#_reportMetrics(metrics) every 10 seconds by default
const myCounter = registry.getOrCreateCounter('my-counter');

```

### [Reporter Abstract Class](https://yaorg.github.io/node-measured/Reporter.html)
Extend this class and override the [_reportMetrics(metrics)](https://yaorg.github.io/node-measured/Reporter.html#_reportMetrics__anchor) method to create a vendor specific reporter implementation. 

See the [ReporterOptions](http://yaorg.github.io/node-measured/build/docs/packages/measured-reporting/global.html#ReporterOptions) for advanced configuration.

#### Current Implementations
- [SignalFx Reporter](https://yaorg.github.io/node-measured/SignalFxMetricsReporter.html) in the `measured-signalfx-reporter` package.
  - reports metrics to SignalFx.
- [Logging Reporter](https://yaorg.github.io/node-measured/LoggingReporter.html) in the `measured-reporting` package.
  - A reporter impl that simply logs the metrics via the Logger

#### Creating an anonymous Implementation
You can technically create an anonymous instance of this, see the following example.
```javascript
const os = require('os');
const process = require('process');
const { SelfReportingMetricsRegistry, Reporter } = require('measured-reporting');

// Create a self reporting registry with an anonymous Reporter instance;
const registry = new SelfReportingMetricsRegistry(
  new class extends Reporter {
    constructor() {
      super({
        defaultDimensions: {
          hostname: os.hostname(),
          env: process.env['NODE_ENV'] ? process.env['NODE_ENV'] : 'unset'
        }
      })
    }

    _reportMetrics(metrics) {
      metrics.forEach(metric => {
        console.log(JSON.stringify({
          metricName: metric.name,
          dimensions: this._getDimensions(metric),
          data: metric.metricImpl.toJSON()
        }))
      });
    }
  }()
);

// create a gauge that reports the process uptime every second
const processUptimeGauge = registry.getOrCreateGauge('node.process.uptime', () => process.uptime(), {}, 1);
```

Example output:
```bash
APP5HTD6ACCD8C:foo jfiel2$ NODE_ENV=development node index.js
{"metricName":"node.process.uptime","dimensions":{"hostname":"APP5HTD6ACCD8C","env":"development"},"data":0.092}
{"metricName":"node.process.uptime","dimensions":{"hostname":"APP5HTD6ACCD8C","env":"development"},"data":1.099}
{"metricName":"node.process.uptime","dimensions":{"hostname":"APP5HTD6ACCD8C","env":"development"},"data":2.104}
{"metricName":"node.process.uptime","dimensions":{"hostname":"APP5HTD6ACCD8C","env":"development"},"data":3.105}
{"metricName":"node.process.uptime","dimensions":{"hostname":"APP5HTD6ACCD8C","env":"development"},"data":4.106}
```


Consider creating a proper class and contributing it back to Measured if it is generic and sharable.

### [Logging Reporter Class](https://yaorg.github.io/node-measured/LoggingReporter.html)
A simple reporter that logs the metrics via the Logger.

See the [ReporterOptions](http://yaorg.github.io/node-measured/build/docs/packages/measured-reporting/global.html#ReporterOptions) for advanced configuration.

```javascript
const { SelfReportingMetricsRegistry, LoggingReporter } = require('measured-reporting');
const registry = new SelfReportingMetricsRegistry(new LoggingReporter({
  logger: myLogerImpl, // defaults to new console logger if not supplied
  defaultDimensions: {
    hostname: require('os').hostname()
  }
}));
```

## What are dimensions?
As described by Signal Fx:
    
*A dimension is a key/value pair that, along with the metric name, is part of the identity of a time series. 
You can filter and aggregate time series by those dimensions across SignalFx.*
    
DataDog has a [nice blog post](https://www.datadoghq.com/blog/the-power-of-tagged-metrics/) about how they are used in their aggregator api.

Graphite also supports the concept via [tags](http://graphite.readthedocs.io/en/latest/tags.html).

