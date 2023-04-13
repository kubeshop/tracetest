# Measured Core

The core measured library that has the Metric interfaces and implementations.

[![npm](https://img.shields.io/npm/v/measured-core.svg)](https://www.npmjs.com/package/measured-core) 

## Install

```
npm install measured-core
```

## What is in this package

### Metric Implemenations

The core library has the following metrics classes:

#### [Gauge](https://yaorg.github.io/node-measured/packages/measured-core/Gauge.html)
Values that can be read instantly via a supplied call back.

#### [SettableGauge](https://yaorg.github.io/node-measured/packages/measured-core/SettableGauge.html)
Just like a Gauge but its value is set directly rather than supplied by a callback.

#### [CachedGauge](https://yaorg.github.io/node-measured/packages/measured-core/CachedGauge.html)
Like a mix of the regular and settable Gauge it takes a call back that returns a promise that will resolve the cached value and an interval that it should call the callback on to update its cached value.

#### [Counter](https://yaorg.github.io/node-measured/packages/measured-core/Counter.html)
Counters are things that increment or decrement.

#### [Timer](https://yaorg.github.io/node-measured/packages/measured-core/Timer.html)
Timers are a combination of Meters and Histograms. They measure the rate as well as distribution of scalar events.

#### [Histogram](https://yaorg.github.io/node-measured/packages/measured-core/Histogram.html)
Keeps a reservoir of statistically relevant values to explore their distribution.

#### [Meter](https://yaorg.github.io/node-measured/packages/measured-core/Meter.html)
Things that are measured as events / interval.

### Registry

The core library comes with a basic registry class 

#### [Collection](https://yaorg.github.io/node-measured/packages/measured-core/Collection.html)

that is not aware of dimensions / tags and leaves reporting up to you.

#### See the [measured-reporting](../measured-reporting/) module for more advanced and featured registries.

### Other

See The [measured-core](https://yaorg.github.io/node-measured/packages/measured-core/module-measured-core.html) modules for the full list of exports for require('measured-core').

## Usage

**Step 1:** Add measurements to your code. For example, lets track the
requests/sec of a http server:

```js
var http  = require('http');
var stats = require('measured').createCollection();

http.createServer(function(req, res) {
  stats.meter('requestsPerSecond').mark();
  res.end('Thanks');
}).listen(3000);
```

**Step 2:** Show the collected measurements (more advanced examples follow later):

```js
setInterval(function() {
  console.log(stats.toJSON());
}, 1000);
```

This will output something like this every second:

```
{ requestsPerSecond:
   { mean: 1710.2180279856818,
     count: 10511,
     'currentRate': 1941.4893498239829,
     '1MinuteRate': 168.08263156623656,
     '5MinuteRate': 34.74630977619571,
     '15MinuteRate': 11.646507524106095 } }
```

**Step 3:** Aggregate the data into your backend of choice.
Here are a few time series data aggregators.
- [Graphite](http://graphite.wikidot.com/)
    - A free and open source, self hosted and managed solution for time series data.
- [SignalFx](https://signalfx.com/)
    - An enterprise SASS offering for time series data.
- [Datadog](https://www.datadoghq.com/)
    - An enterprise SASS offering for time series data.
