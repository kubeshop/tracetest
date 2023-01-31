const apm = require('elastic-apm-node').start({
  serviceName: 'sample-app',
  serverUrl: 'http://apm-server:8200',
})
