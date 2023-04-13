/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

'use strict'

var fs = require('fs')
var path = require('path')

var ElasticAPMHttpClient = require('elastic-apm-http-client')
var truncate = require('unicode-byte-truncate')

const REDACTED = require('./constants').REDACTED
var logging = require('./logging')
var version = require('../package').version

const { WildcardMatcher } = require('./wildcard-matcher')
const { CloudMetadata } = require('./cloud-metadata')
const { NoopTransport } = require('./noop-transport')
const { isLambdaExecutionEnvironment } = require('./lambda')

let confFile = loadConfigFile()

const INTAKE_STRING_MAX_SIZE = 1024
const CAPTURE_ERROR_LOG_STACK_TRACES_NEVER = 'never'
const CAPTURE_ERROR_LOG_STACK_TRACES_MESSAGES = 'messages'
const CAPTURE_ERROR_LOG_STACK_TRACES_ALWAYS = 'always'
const CONTEXT_MANAGER_PATCH = 'patch'
const CONTEXT_MANAGER_ASYNCHOOKS = 'asynchooks'
const CONTEXT_MANAGER_ASYNCLOCALSTORAGE = 'asynclocalstorage'
const TRACE_CONTINUATION_STRATEGY_CONTINUE = 'continue'
const TRACE_CONTINUATION_STRATEGY_RESTART = 'restart'
const TRACE_CONTINUATION_STRATEGY_RESTART_EXTERNAL = 'restart_external'

var DEFAULTS = {
  abortedErrorThreshold: '25s',
  active: true,
  addPatch: undefined,
  apiRequestSize: '768kb',
  apiRequestTime: '10s',
  breakdownMetrics: true,
  captureBody: 'off',
  captureErrorLogStackTraces: CAPTURE_ERROR_LOG_STACK_TRACES_MESSAGES,
  captureExceptions: true,
  captureHeaders: true,
  centralConfig: true,
  cloudProvider: 'auto',
  containerId: undefined,
  // 'contextManager' and 'asyncHooks' are explicitly *not* included in DEFAULTS
  // because normalizeContextManager() needs to know if a value was provided by
  // the user.
  contextPropagationOnly: false,
  disableInstrumentations: [],
  disableSend: false,
  elasticsearchCaptureBodyUrls: [
    '*/_search', '*/_search/template', '*/_msearch', '*/_msearch/template',
    '*/_async_search', '*/_count', '*/_sql', '*/_eql/search'
  ],
  environment: process.env.NODE_ENV || 'development',
  errorOnAbortedRequests: false,
  exitSpanMinDuration: '0ms',
  filterHttpHeaders: true,
  globalLabels: undefined,
  ignoreMessageQueues: [],
  instrument: true,
  instrumentIncomingHTTPRequests: true,
  kubernetesNamespace: undefined,
  kubernetesNodeName: undefined,
  kubernetesPodName: undefined,
  kubernetesPodUID: undefined,
  logLevel: 'info',
  logUncaughtExceptions: false, // TODO: Change to `true` in the v4.0.0
  longFieldMaxLength: 10000,
  // Rough equivalent of the Java Agent's max_queue_size:
  // https://www.elastic.co/guide/en/apm/agent/java/current/config-reporter.html#config-max-queue-size
  maxQueueSize: 1024,
  metricsInterval: '30s',
  metricsLimit: 1000,
  opentelemetryBridgeEnabled: false,
  sanitizeFieldNames: [
    // These patterns are specified in the shared APM specs:
    // https://github.com/elastic/apm/blob/main/specs/agents/sanitization.md
    'password', 'passwd', 'pwd', 'secret', '*key', '*token*',
    '*session*', '*credit*', '*card*', '*auth*', 'set-cookie', '*principal*',
    // These are default patterns only in the Node.js APM agent, historically
    // from when the "is-secret" dependency was used.
    'pw', 'pass', 'connect.sid'
  ],
  serviceNodeName: undefined,
  serverTimeout: '30s',
  serverUrl: 'http://127.0.0.1:8200',
  sourceLinesErrorAppFrames: 5,
  sourceLinesErrorLibraryFrames: 5,
  sourceLinesSpanAppFrames: 0,
  sourceLinesSpanLibraryFrames: 0,
  spanCompressionEnabled: true,
  spanCompressionExactMatchMaxDuration: '50ms',
  spanCompressionSameKindMaxDuration: '0ms',
  // 'spanStackTraceMinDuration' is explicitly *not* included in DEFAULTS
  // because normalizeSpanStackTraceMinDuration() needs to know if a value
  // was provided by the user.
  stackTraceLimit: 50,
  traceContinuationStrategy: TRACE_CONTINUATION_STRATEGY_CONTINUE,
  transactionIgnoreUrls: [],
  transactionMaxSpans: 500,
  transactionSampleRate: 1.0,
  useElasticTraceparentHeader: true,
  usePathAsTransactionName: false,
  verifyServerCert: true
}

var ENV_TABLE = {
  abortedErrorThreshold: 'ELASTIC_APM_ABORTED_ERROR_THRESHOLD',
  active: 'ELASTIC_APM_ACTIVE',
  addPatch: 'ELASTIC_APM_ADD_PATCH',
  apiKey: 'ELASTIC_APM_API_KEY',
  apiRequestSize: 'ELASTIC_APM_API_REQUEST_SIZE',
  apiRequestTime: 'ELASTIC_APM_API_REQUEST_TIME',
  asyncHooks: 'ELASTIC_APM_ASYNC_HOOKS',
  breakdownMetrics: 'ELASTIC_APM_BREAKDOWN_METRICS',
  captureBody: 'ELASTIC_APM_CAPTURE_BODY',
  captureErrorLogStackTraces: 'ELASTIC_APM_CAPTURE_ERROR_LOG_STACK_TRACES',
  captureExceptions: 'ELASTIC_APM_CAPTURE_EXCEPTIONS',
  captureHeaders: 'ELASTIC_APM_CAPTURE_HEADERS',
  captureSpanStackTraces: 'ELASTIC_APM_CAPTURE_SPAN_STACK_TRACES',
  centralConfig: 'ELASTIC_APM_CENTRAL_CONFIG',
  cloudProvider: 'ELASTIC_APM_CLOUD_PROVIDER',
  containerId: 'ELASTIC_APM_CONTAINER_ID',
  contextManager: 'ELASTIC_APM_CONTEXT_MANAGER',
  contextPropagationOnly: 'ELASTIC_APM_CONTEXT_PROPAGATION_ONLY',
  disableInstrumentations: 'ELASTIC_APM_DISABLE_INSTRUMENTATIONS',
  disableSend: 'ELASTIC_APM_DISABLE_SEND',
  environment: 'ELASTIC_APM_ENVIRONMENT',
  exitSpanMinDuration: 'ELASTIC_APM_EXIT_SPAN_MIN_DURATION',
  ignoreMessageQueues: ['ELASTIC_IGNORE_MESSAGE_QUEUES', 'ELASTIC_APM_IGNORE_MESSAGE_QUEUES'],
  elasticsearchCaptureBodyUrls: 'ELASTIC_APM_ELASTICSEARCH_CAPTURE_BODY_URLS',
  errorMessageMaxLength: 'ELASTIC_APM_ERROR_MESSAGE_MAX_LENGTH',
  errorOnAbortedRequests: 'ELASTIC_APM_ERROR_ON_ABORTED_REQUESTS',
  filterHttpHeaders: 'ELASTIC_APM_FILTER_HTTP_HEADERS',
  frameworkName: 'ELASTIC_APM_FRAMEWORK_NAME',
  frameworkVersion: 'ELASTIC_APM_FRAMEWORK_VERSION',
  globalLabels: 'ELASTIC_APM_GLOBAL_LABELS',
  hostname: 'ELASTIC_APM_HOSTNAME',
  instrument: 'ELASTIC_APM_INSTRUMENT',
  instrumentIncomingHTTPRequests: 'ELASTIC_APM_INSTRUMENT_INCOMING_HTTP_REQUESTS',
  kubernetesNamespace: ['ELASTIC_APM_KUBERNETES_NAMESPACE', 'KUBERNETES_NAMESPACE'],
  kubernetesNodeName: ['ELASTIC_APM_KUBERNETES_NODE_NAME', 'KUBERNETES_NODE_NAME'],
  kubernetesPodName: ['ELASTIC_APM_KUBERNETES_POD_NAME', 'KUBERNETES_POD_NAME'],
  kubernetesPodUID: ['ELASTIC_APM_KUBERNETES_POD_UID', 'KUBERNETES_POD_UID'],
  logLevel: 'ELASTIC_APM_LOG_LEVEL',
  logUncaughtExceptions: 'ELASTIC_APM_LOG_UNCAUGHT_EXCEPTIONS',
  longFieldMaxLength: 'ELASTIC_APM_LONG_FIELD_MAX_LENGTH',
  maxQueueSize: 'ELASTIC_APM_MAX_QUEUE_SIZE',
  metricsInterval: 'ELASTIC_APM_METRICS_INTERVAL',
  metricsLimit: 'ELASTIC_APM_METRICS_LIMIT',
  opentelemetryBridgeEnabled: 'ELASTIC_APM_OPENTELEMETRY_BRIDGE_ENABLED',
  payloadLogFile: 'ELASTIC_APM_PAYLOAD_LOG_FILE',
  sanitizeFieldNames: ['ELASTIC_SANITIZE_FIELD_NAMES', 'ELASTIC_APM_SANITIZE_FIELD_NAMES'],
  serverCaCertFile: 'ELASTIC_APM_SERVER_CA_CERT_FILE',
  secretToken: 'ELASTIC_APM_SECRET_TOKEN',
  serverTimeout: 'ELASTIC_APM_SERVER_TIMEOUT',
  serverUrl: 'ELASTIC_APM_SERVER_URL',
  serviceName: 'ELASTIC_APM_SERVICE_NAME',
  serviceNodeName: 'ELASTIC_APM_SERVICE_NODE_NAME',
  serviceVersion: 'ELASTIC_APM_SERVICE_VERSION',
  sourceLinesErrorAppFrames: 'ELASTIC_APM_SOURCE_LINES_ERROR_APP_FRAMES',
  sourceLinesErrorLibraryFrames: 'ELASTIC_APM_SOURCE_LINES_ERROR_LIBRARY_FRAMES',
  sourceLinesSpanAppFrames: 'ELASTIC_APM_SOURCE_LINES_SPAN_APP_FRAMES',
  sourceLinesSpanLibraryFrames: 'ELASTIC_APM_SOURCE_LINES_SPAN_LIBRARY_FRAMES',
  spanCompressionEnabled: 'ELASTIC_APM_SPAN_COMPRESSION_ENABLED',
  spanCompressionExactMatchMaxDuration: 'ELASTIC_APM_SPAN_COMPRESSION_EXACT_MATCH_MAX_DURATION',
  spanCompressionSameKindMaxDuration: 'ELASTIC_APM_SPAN_COMPRESSION_SAME_KIND_MAX_DURATION',
  spanStackTraceMinDuration: 'ELASTIC_APM_SPAN_STACK_TRACE_MIN_DURATION',
  spanFramesMinDuration: 'ELASTIC_APM_SPAN_FRAMES_MIN_DURATION',
  stackTraceLimit: 'ELASTIC_APM_STACK_TRACE_LIMIT',
  traceContinuationStrategy: 'ELASTIC_APM_TRACE_CONTINUATION_STRATEGY',
  transactionIgnoreUrls: 'ELASTIC_APM_TRANSACTION_IGNORE_URLS',
  transactionMaxSpans: 'ELASTIC_APM_TRANSACTION_MAX_SPANS',
  transactionSampleRate: 'ELASTIC_APM_TRANSACTION_SAMPLE_RATE',
  useElasticTraceparentHeader: 'ELASTIC_APM_USE_ELASTIC_TRACEPARENT_HEADER',
  usePathAsTransactionName: 'ELASTIC_APM_USE_PATH_AS_TRANSACTION_NAME',
  verifyServerCert: 'ELASTIC_APM_VERIFY_SERVER_CERT'
}

var CENTRAL_CONFIG = {
  log_level: 'logLevel',
  transaction_sample_rate: 'transactionSampleRate',
  transaction_max_spans: 'transactionMaxSpans',
  capture_body: 'captureBody',
  transaction_ignore_urls: 'transactionIgnoreUrls',
  sanitize_field_names: 'sanitizeFieldNames',
  ignore_message_queues: 'ignoreMessageQueues',
  span_stack_trace_min_duration: 'spanStackTraceMinDuration',
  trace_continuation_strategy: 'traceContinuationStrategy',
  exit_span_min_duration: 'exitSpanMinDuration'
}

var BOOL_OPTS = [
  'active',
  'asyncHooks',
  'breakdownMetrics',
  'captureExceptions',
  'captureHeaders',
  'captureSpanStackTraces',
  'centralConfig',
  'contextPropagationOnly',
  'disableSend',
  'errorOnAbortedRequests',
  'filterHttpHeaders',
  'instrument',
  'instrumentIncomingHTTPRequests',
  'logUncaughtExceptions',
  'opentelemetryBridgeEnabled',
  'spanCompressionEnabled',
  'usePathAsTransactionName',
  'verifyServerCert'
]

var NUM_OPTS = [
  'longFieldMaxLength',
  'maxQueueSize',
  'metricsLimit',
  'sourceLinesErrorAppFrames',
  'sourceLinesErrorLibraryFrames',
  'sourceLinesSpanAppFrames',
  'sourceLinesSpanLibraryFrames',
  'stackTraceLimit',
  'transactionMaxSpans',
  'transactionSampleRate'
]

var DURATION_OPTS = [
  {
    name: 'abortedErrorThreshold',
    defaultUnit: 's',
    allowedUnits: ['ms', 's', 'm'],
    allowNegative: false
  },
  {
    name: 'apiRequestTime',
    defaultUnit: 's',
    allowedUnits: ['ms', 's', 'm'],
    allowNegative: false
  },
  {
    name: 'exitSpanMinDuration',
    defaultUnit: 'ms',
    allowedUnits: ['us', 'ms', 's', 'm'],
    allowNegative: false
  },
  {
    name: 'metricsInterval',
    defaultUnit: 's',
    allowedUnits: ['ms', 's', 'm'],
    allowNegative: false
  },
  {
    name: 'serverTimeout',
    defaultUnit: 's',
    allowedUnits: ['ms', 's', 'm'],
    allowNegative: false
  },
  {
    name: 'spanCompressionExactMatchMaxDuration',
    defaultUnit: 'ms',
    allowedUnits: ['ms', 's', 'm'],
    allowNegative: false
  },
  {
    name: 'spanCompressionSameKindMaxDuration',
    defaultUnit: 'ms',
    allowedUnits: ['ms', 's', 'm'],
    allowNegative: false
  },
  {
    // Deprecated: use `spanStackTraceMinDuration`.
    name: 'spanFramesMinDuration',
    defaultUnit: 's',
    allowedUnits: ['ms', 's', 'm'],
    allowNegative: true
  },
  {
    name: 'spanStackTraceMinDuration',
    defaultUnit: 'ms',
    allowedUnits: ['ms', 's', 'm'],
    allowNegative: true
  }
]

var BYTES_OPTS = [
  'apiRequestSize',
  'errorMessageMaxLength'
]

var MINUS_ONE_EQUAL_INFINITY = [
  'transactionMaxSpans'
]

var ARRAY_OPTS = [
  'disableInstrumentations',
  'elasticsearchCaptureBodyUrls',
  'sanitizeFieldNames',
  'transactionIgnoreUrls',
  'ignoreMessageQueues'
]

var KEY_VALUE_OPTS = [
  'addPatch',
  'globalLabels'
]

// Configure a logger for the agent.
//
// This is separate from `createConfig` to allow the agent to have an early
// logger before `agent.start()` is called.
function configLogger (opts) {
  const logLevel = (
    process.env[ENV_TABLE.logLevel] ||
    (opts && opts.logLevel) ||
    (confFile && confFile.logLevel) ||
    DEFAULTS.logLevel
  )

  // `ELASTIC_APM_LOGGER=false` is provided as a mechanism to *disable* a
  // custom logger for troubleshooting because a wrapped custom logger does
  // not get structured log data.
  // https://www.elastic.co/guide/en/apm/agent/nodejs/current/troubleshooting.html#debug-mode
  let customLogger = null
  if (process.env.ELASTIC_APM_LOGGER !== 'false') {
    customLogger = (
      (opts && opts.logger) ||
      (confFile && confFile.logger)
    )
  }

  return logging.createLogger(logLevel, customLogger)
}

// Create an initial configuration from DEFAULTS. This is used as a stand-in
// for Agent configuration until `agent.start(...)` is called.
function initialConfig (logger) {
  const cfg = Object.assign({}, DEFAULTS)

  // Reproduce the generated properties for `Config`.
  cfg.ignoreUrlStr = []
  cfg.ignoreUrlRegExp = []
  cfg.ignoreUserAgentStr = []
  cfg.ignoreUserAgentRegExp = []
  cfg.elasticsearchCaptureBodyUrlsRegExp = []
  cfg.transactionIgnoreUrlRegExp = []
  cfg.sanitizeFieldNamesRegExp = []
  cfg.ignoreMessageQueuesRegExp = []
  normalize(cfg, logger)

  cfg.transport = new NoopTransport()

  return cfg
}

function createConfig (opts, logger) {
  return new Config(opts, logger)
}

class Config {
  constructor (opts, logger) {
    this.ignoreUrlStr = []
    this.ignoreUrlRegExp = []
    this.ignoreUserAgentStr = []
    this.ignoreUserAgentRegExp = []
    this.elasticsearchCaptureBodyUrlsRegExp = []
    this.transactionIgnoreUrlRegExp = []
    this.sanitizeFieldNamesRegExp = []
    this.ignoreMessageQueuesRegExp = []

    const isLambda = isLambdaExecutionEnvironment()

    // If we didn't find a config file on process boot, but a path to one is
    // provided as a config option, let's instead try to load that
    if (confFile === null && opts && opts.configFile) {
      confFile = loadConfigFile(opts.configFile)
    }

    Object.assign(
      this,
      DEFAULTS, // default options
      confFile, // options read from config file
      opts, // options passed in to agent.start()
      readEnv() // options read from environment variables
    )

    // The logger is used later in this function, so create/update it first.
    // Unless a new custom `logger` was provided, we use the one created earlier
    // in `configLogger()`.
    const customLogger = (process.env.ELASTIC_APM_LOGGER === 'false' ? null : this.logger)
    if (!customLogger && logger) {
      logging.setLogLevel(logger, this.logLevel)
      this.logger = logger
    } else {
      this.logger = logging.createLogger(this.logLevel, customLogger)
    }

    // Fallback and validation handling for `serviceName` and `serviceVersion`.
    if (this.serviceName) {
      // A value here means an explicit value was given. Error out if invalid.
      try {
        validateServiceName(this.serviceName)
      } catch (err) {
        this.logger.error('serviceName "%s" is invalid: %s', this.serviceName, err.message)
        this.serviceName = null
      }
    } else if (isLambda) {
      this.serviceName = process.env.AWS_LAMBDA_FUNCTION_NAME
    } else {
      // Zero-conf support: use package.json#name, else
      // `unknown-${service.agent.name}-service`.
      try {
        this.serviceName = serviceNameFromPackageJson()
      } catch (err) {
        this.logger.warn(err.message)
      }
      if (!this.serviceName) {
        this.serviceName = 'unknown-nodejs-service'
      }
    }
    if (this.serviceVersion) {
      // pass
    } else if (isLambda) {
      this.serviceVersion = process.env.AWS_LAMBDA_FUNCTION_VERSION
    } else {
      // Zero-conf support: use package.json#version, if possible.
      try {
        this.serviceVersion = serviceVersionFromPackageJson()
      } catch (err) {
        // pass
      }
    }

    normalize(this, this.logger)

    if (isLambda) {
      // Override some config in AWS Lambda environment.
      this.metricsInterval = 0
      this.cloudProvider = 'none'
      this.centralConfig = false
    }
    if (this.metricsInterval === 0) {
      this.breakdownMetrics = false
    }

    if (this.disableSend || this.contextPropagationOnly) {
      this.transport = function createNoopTransport (conf, agent) {
        return new NoopTransport()
      }
    } else if (typeof this.transport !== 'function') {
      this.transport = function httpTransport (conf, agent) {
        const config = getBaseClientConfig(conf, agent)
        var transport = new ElasticAPMHttpClient(config)

        transport.on('config', remoteConf => {
          agent.logger.debug({ remoteConf }, 'central config received')
          const conf = {}
          const unknown = []

          for (const [key, value] of Object.entries(remoteConf)) {
            const newKey = CENTRAL_CONFIG[key]
            if (newKey) {
              conf[newKey] = value
            } else {
              unknown.push(key)
            }
          }

          if (unknown.length > 0) {
            agent.logger.warn(`Central config warning: unsupported config names: ${unknown.join(', ')}`)
          }

          if (Object.keys(conf).length > 0) {
            normalize(conf, agent.logger)

            for (const [key, value] of Object.entries(conf)) {
              const oldValue = agent._conf[key]
              agent._conf[key] = value
              if (key === 'logLevel' && value !== oldValue && !logging.isLoggerCustom(agent.logger)) {
                logging.setLogLevel(agent.logger, value)
                agent.logger.info(`Central config success: updated logger with new logLevel: ${value}`)
              }
              agent.logger.info(`Central config success: updated ${key}: ${value}`)
            }
          }
        })

        transport.on('error', err => {
          agent.logger.error('APM Server transport error: %s', err.stack)
        })

        transport.on('request-error', err => {
          const haveAccepted = Number.isFinite(err.accepted)
          const haveErrors = Array.isArray(err.errors)
          let msg

          if (err.code === 404) {
            msg = 'APM Server responded with "404 Not Found". ' +
              'This might be because you\'re running an incompatible version of the APM Server. ' +
              'This agent only supports APM Server v6.5 and above. ' +
              'If you\'re using an older version of the APM Server, ' +
              'please downgrade this agent to version 1.x or upgrade the APM Server'
          } else if (err.code) {
            msg = `APM Server transport error (${err.code}): ${err.message}`
          } else {
            msg = `APM Server transport error: ${err.message}`
          }

          if (haveAccepted || haveErrors) {
            if (haveAccepted) msg += `\nAPM Server accepted ${err.accepted} events in the last request`
            if (haveErrors) {
              for (const error of err.errors) {
                msg += `\nError: ${error.message}`
                if (error.document) msg += `\n  Document: ${error.document}`
              }
            }
          } else if (err.response) {
            msg += `\n${err.response}`
          }

          agent.logger.error(msg)
        })

        return transport
      }
    }
  }

  // Return a reasonably loggable object for this Config instance.
  // Exclude undefined fields and complex objects like `logger`.
  toJSON () {
    const EXCLUDE_FIELDS = {
      logger: true,
      transport: true
    }
    const REDACT_FIELDS = {
      apiKey: true,
      secretToken: true,
      serverUrl: true
    }
    const NICE_REGEXPS_FIELDS = {
      ignoreUrlRegExp: true,
      ignoreUserAgentRegExp: true,
      transactionIgnoreUrlRegExp: true,
      sanitizeFieldNamesRegExp: true,
      ignoreMessageQueuesRegExp: true
    }
    const loggable = {}
    for (const k in this) {
      if (EXCLUDE_FIELDS[k] || this[k] === undefined) {
        // pass
      } else if (REDACT_FIELDS[k]) {
        loggable[k] = REDACTED
      } else if (NICE_REGEXPS_FIELDS[k] && Array.isArray(this[k])) {
        // JSON.stringify() on a RegExp is "{}", which isn't very helpful.
        loggable[k] = this[k].map(r => r instanceof RegExp ? r.toString() : r)
      } else {
        loggable[k] = this[k]
      }
    }
    return loggable
  }
}

function readEnv () {
  var opts = {}
  for (const key of Object.keys(ENV_TABLE)) {
    let env = ENV_TABLE[key]
    if (!Array.isArray(env)) env = [env]
    for (const envKey of env) {
      if (envKey in process.env) {
        opts[key] = process.env[envKey]
      }
    }
  }
  return opts
}

function validateServiceName (s) {
  if (typeof s !== 'string') {
    throw new Error('not a string')
  } else if (!/^[a-zA-Z0-9 _-]+$/.test(s)) {
    throw new Error('contains invalid characters (allowed: a-z, A-Z, 0-9, _, -, <space>)')
  }
}

// findPkgInfo() looks up from the script dir (or cwd) for a "package.json" file
// from which to load the name and version. It returns:
//    {
//      startDir: "<full path to starting dir>",
//      path: "/the/full/path/to/package.json",  // may be null
//      data: {
//        name: "<the package name>",            // may be missing
//        version: "<the package version>"       // may be missing
//      }
//    }
let pkgInfoCache
function findPkgInfo () {
  if (pkgInfoCache === undefined) {
    // Determine a good starting dir from which to look for a "package.json".
    let startDir = require.main && require.main.filename && path.dirname(require.main.filename)
    if (!startDir && process.argv[1]) {
      // 'require.main' is undefined if the agent is preloaded with `node
      // --require elastic-apm-node/... script.js`.
      startDir = path.dirname(process.argv[1])
    }
    if (!startDir) {
      process.cwd()
    }
    pkgInfoCache = {
      startDir,
      path: null,
      data: {}
    }

    // Look up from the starting dir for a "package.json".
    const { root } = path.parse(startDir)
    let dir = startDir
    while (true) {
      const pj = path.resolve(dir, 'package.json')
      if (fs.existsSync(pj)) {
        pkgInfoCache.path = pj
        break
      }
      if (dir === root) {
        break
      }
      dir = path.dirname(dir)
    }

    // Attempt to load "name" and "version" from the package.json.
    if (pkgInfoCache.path) {
      try {
        const data = JSON.parse(fs.readFileSync(pkgInfoCache.path))
        if (data.name) {
          // For backward compatibility, maintain the trimming done by
          // https://github.com/npm/normalize-package-data#what-normalization-currently-entails
          pkgInfoCache.data.name = data.name.trim()
        }
        if (data.version) {
          pkgInfoCache.data.version = data.version
        }
      } catch (_err) {
        // Silently leave data empty.
      }
    }
  }
  return pkgInfoCache
}

function serviceNameFromPackageJson () {
  const pkg = findPkgInfo()
  if (!pkg.path) {
    throw new Error(`could not infer serviceName: could not find package.json up from ${pkg.startDir}`)
  }
  if (!pkg.data.name) {
    throw new Error(`could not infer serviceName: "${pkg.path}" does not contain a "name"`)
  }
  if (typeof pkg.data.name !== 'string') {
    throw new Error(`could not infer serviceName: "name" in "${pkg.path}" is not a string`)
  }
  let serviceName = pkg.data.name

  // Normalize a namespaced npm package name, '@ns/name', to 'ns-name'.
  const match = /^@([^/]+)\/([^/]+)$/.exec(serviceName)
  if (match) {
    serviceName = match[1] + '-' + match[2]
  }

  // Sanitize, by replacing invalid service name chars with an underscore.
  const SERVICE_NAME_BAD_CHARS = /[^a-zA-Z0-9 _-]/g
  serviceName = serviceName.replace(SERVICE_NAME_BAD_CHARS, '_')

  // Disallow some weird sanitized values. For example, it is better to
  // have the fallback "unknown-{service.agent.name}-service" than "_" or
  // "____" or " ".
  const ALL_NON_ALPHANUMERIC = /^[ _-]*$/
  if (ALL_NON_ALPHANUMERIC.test(serviceName)) {
    serviceName = null
  }
  if (!serviceName) {
    throw new Error(`could not infer serviceName from name="${pkg.data.name}" in "${pkg.path}"`)
  }

  return serviceName
}

function serviceVersionFromPackageJson () {
  const pkg = findPkgInfo()
  if (!pkg.path) {
    throw new Error(`could not infer serviceVersion: could not find package.json up from ${pkg.startDir}`)
  }
  if (!pkg.data.version) {
    throw new Error(`could not infer serviceVersion: "${pkg.path}" does not contain a "version"`)
  }
  if (typeof pkg.data.version !== 'string') {
    throw new Error(`could not infer serviceVersion: "version" in "${pkg.path}" is not a string`)
  }
  return pkg.data.version
}

function normalize (opts, logger) {
  normalizeKeyValuePairs(opts)
  normalizeNumbers(opts)
  normalizeBytes(opts)
  normalizeArrays(opts)
  normalizeDurationOptions(opts, logger)
  normalizeBools(opts, logger)
  normalizeIgnoreOptions(opts)
  normalizeElasticsearchCaptureBodyUrls(opts)
  normalizeSanitizeFieldNames(opts)
  normalizeContextManager(opts, logger) // Must be after normalizeBools().
  normalizeCloudProvider(opts, logger)
  normalizeTransactionSampleRate(opts, logger)
  normalizeTraceContinuationStrategy(opts, logger)

  // This must be after `normalizeDurationOptions()` and `normalizeBools()`
  // because it synthesizes the deprecated `spanFramesMinDuration` and
  // `captureSpanStackTraces` options into `spanStackTraceMinDuration`.
  normalizeSpanStackTraceMinDuration(opts, logger)

  truncateOptions(opts)
}

const ALLOWED_TRACE_CONTINUATION_STRATEGY = {
  [TRACE_CONTINUATION_STRATEGY_CONTINUE]: true,
  [TRACE_CONTINUATION_STRATEGY_RESTART]: true,
  [TRACE_CONTINUATION_STRATEGY_RESTART_EXTERNAL]: true
}
function normalizeTraceContinuationStrategy (opts, logger) {
  if ('traceContinuationStrategy' in opts &&
      !(opts.traceContinuationStrategy in ALLOWED_TRACE_CONTINUATION_STRATEGY)) {
    logger.warn('Invalid "traceContinuationStrategy" config value %j, falling back to default %j',
      opts.traceContinuationStrategy, DEFAULTS.traceContinuationStrategy)
    opts.traceContinuationStrategy = DEFAULTS.traceContinuationStrategy
  }
}

// Normalize provided values for `spanFramesMinDuration` (deprecated),
// `captureSpanStackTraces` (deprecated) and `spanStackTraceMinDuration` into
// a final value for `spanStackTraceMinDuration` that is used by the agent.
//
// This function expects `normalizeDurationOptions()` and `normalizeBools()`
// to have already been called.
//
// | spanStackTraceMinDuration | captureSpanStackTraces | spanFramesMinDuration   | resultant spanStackTraceMinDuration |
// | ------------------------- | ---------------------- | ----------------------- | ----------------------------------- |
// | -                         | -                      | -                       | `-1ms` (no span stacktraces)        |
// | `-1ms` (negative value)   | (any value is ignored) | (any value is ignored)  | `-1ms` (no span stacktraces)        |
// | `0ms` (zero value)        | (any value is ignored) | (any value is ignored)  | `0ms` (stacktraces for all spans)   |
// | `5ms` (positive value)    | (any value is ignored) | (any value is ignored)  | `5ms`                               |
// | -                         | `false`                | (any value)             | `-1ms` (no span stacktraces)        |
// | -                         | `true`                 | -                       | `10ms` (backwards compatible value) |
// | -                         | `true` or unspecified  | `0ms` (zero value)      | `-1ms` (no span stacktraces)        |
// | -                         | `true` or unspecified  | `-1ms` (negative value) | `0ms` (stacktraces for all spans)   |
// | -                         | `true` or unspecified  | `5ms` (positive value)  | `5ms`                               |
function normalizeSpanStackTraceMinDuration (opts, logger) {
  const before = {}
  if (opts.captureSpanStackTraces !== undefined) before.captureSpanStackTraces = opts.captureSpanStackTraces
  if (opts.spanFramesMinDuration !== undefined) before.spanFramesMinDuration = opts.spanFramesMinDuration
  if (opts.spanStackTraceMinDuration !== undefined) before.spanStackTraceMinDuration = opts.spanStackTraceMinDuration

  if ('spanStackTraceMinDuration' in opts) {
    // If the new option was specified, then use it and ignore the old two.
  } else if (opts.captureSpanStackTraces === false) {
    opts.spanStackTraceMinDuration = -1 // Turn off span stacktraces.
  } else if ('spanFramesMinDuration' in opts) {
    if (opts.spanFramesMinDuration === 0) {
      opts.spanStackTraceMinDuration = -1 // Turn off span stacktraces.
    } else if (opts.spanFramesMinDuration < 0) {
      opts.spanStackTraceMinDuration = 0 // Stacktraces for all spans.
    } else {
      opts.spanStackTraceMinDuration = opts.spanFramesMinDuration
    }
  } else if (opts.captureSpanStackTraces === true) {
    // For backwards compat, use the default `spanFramesMinDuration` value
    // from before `spanStackTraceMinDuration` was introduced.
    opts.spanStackTraceMinDuration = 10 / 1e3 // 10ms
  } else {
    // None of the three options was specified.
    opts.spanStackTraceMinDuration = -1 // Turn off span stacktraces.
  }
  delete opts.captureSpanStackTraces
  delete opts.spanFramesMinDuration

  // Log if something potentially interesting happened here.
  if (Object.keys(before).length > 0) {
    const after = { spanStackTraceMinDuration: opts.spanStackTraceMinDuration }
    logger.trace({ before, after }, 'normalizeSpanStackTraceMinDuration')
  }
}

const ALLOWED_CONTEXT_MANAGER = {
  [CONTEXT_MANAGER_PATCH]: true,
  [CONTEXT_MANAGER_ASYNCHOOKS]: true,
  [CONTEXT_MANAGER_ASYNCLOCALSTORAGE]: true
}

/**
 * Normalize and validate the given values for `contextManager`, and the
 * deprecated `asyncHooks` that it replaces.
 *
 * - `contextManager=patch` means use the "patch-async" technique. I.e., do
 *   limited monkey patching of Node.js core async methods to do limited context
 *   tracking).
 * - `contextManager=asynchooks` means use the "async_hooks.createHook()"
 *   technique. This works in all supported versions of node, but can have
 *   significant performance overhead for Promise-heavy apps.
 * - `contextManager=asynclocalstorage` means use the "AsyncLocalStorage"
 *   technique *if supported in the version of node* (>=14.5 || ^12.19.0).
 *   Otherwise, this will warn and fallback to "asynchooks".
 * - The `asyncHooks` config var is now deprecated. It is translated to the
 *   equivalent `contextManager` value.
 *    - `asyncHooks=false` -> `contextManager=patch`
 *    - `asyncHooks=true` -> leaves the `contextManager` value empty to get
 *      the default behavior: the best async technique.
 * - No specified option means use the best async technique.
 */
function normalizeContextManager (opts, logger) {
  // Treat the empty string, e.g. `ELASTIC_APM_CONTEXT_MANAGER=`, as if it had
  // not been specified.
  if (opts.contextManager === '') {
    delete opts.contextManager
  }

  if ('contextManager' in opts && !(opts.contextManager in ALLOWED_CONTEXT_MANAGER)) {
    logger.warn('Invalid "contextManager" config value %j, falling back to default behavior',
      opts.contextManager)
    delete opts.contextManager
  }

  if ('asyncHooks' in opts) {
    if ('contextManager' in opts) {
      logger.warn({ asyncHooks: opts.asyncHooks, contextManager: opts.contextManager },
        'both `asyncHooks` and `contextManager` config options were specified: the `asyncHooks` value will be ignored')
      delete opts.asyncHooks
    } else if (opts.asyncHooks === false) {
      logger.warn('the `asyncHooks` config option is deprecated; instead of `asyncHooks: false` option, use `contextManager: "patch"`')
      opts.contextManager = 'patch'
      delete opts.asyncHooks
    } else if (opts.asyncHooks === true) {
      logger.warn('the `asyncHooks` config option is deprecated; `asyncHooks: true` is the default behavior')
      delete opts.asyncHooks
    } else {
      delete opts.asyncHooks // Some bogus value.
    }
  }
}

// transactionSampleRate is specified to be:
// - in the range [0,1]
// - rounded to 4 decimal places of precision (e.g. 0.0001, 0.5678, 0.9999)
// - with the special case that a value in the range (0, 0.0001] should be
//   rounded to 0.0001 -- to avoid a small value being rounded to zero.
//
// https://github.com/elastic/apm/blob/main/specs/agents/tracing-sampling.md
function normalizeTransactionSampleRate (opts, logger) {
  if ('transactionSampleRate' in opts) {
    // The value was already run through `Number(...)` in `normalizeNumbers`.
    const rate = opts.transactionSampleRate
    if (isNaN(rate) || rate < 0 || rate > 1) {
      opts.transactionSampleRate = DEFAULTS.transactionSampleRate
      logger.warn('Invalid "transactionSampleRate" config value %s, falling back to default %s',
        rate, opts.transactionSampleRate)
    } else if (rate > 0 && rate < 0.0001) {
      opts.transactionSampleRate = 0.0001
    } else {
      opts.transactionSampleRate = Math.round(rate * 10000) / 10000
    }
  }
}

function normalizeSanitizeFieldNames (opts) {
  if (opts.sanitizeFieldNames) {
    const wildcard = new WildcardMatcher()
    for (const ptn of opts.sanitizeFieldNames) {
      const re = wildcard.compile(ptn)
      opts.sanitizeFieldNamesRegExp.push(re)
    }
  }
}

function normalizeCloudProvider (opts, logger) {
  if ('cloudProvider' in opts) {
    const allowedValues = ['auto', 'gcp', 'azure', 'aws', 'none']
    if (allowedValues.indexOf(opts.cloudProvider) === -1) {
      logger.warn('Invalid "cloudProvider" config value %s, falling back to default %s',
        opts.cloudProvider, DEFAULTS.cloudProvider)
      opts.cloudProvider = DEFAULTS.cloudProvider
    }
  }
}

function normalizeIgnoreOptions (opts) {
  if (opts.transactionIgnoreUrls) {
    // We can't guarantee that opts will be a Config so set a
    // default value. This is to work around CENTRAL_CONFIG tests
    // that call this method with a plain object `{}`
    if (!opts.transactionIgnoreUrlRegExp) {
      opts.transactionIgnoreUrlRegExp = []
    }
    const wildcard = new WildcardMatcher()
    for (const ptn of opts.transactionIgnoreUrls) {
      const re = wildcard.compile(ptn)
      opts.transactionIgnoreUrlRegExp.push(re)
    }
  }

  if (opts.ignoreUrls) {
    for (const ptn of opts.ignoreUrls) {
      if (typeof ptn === 'string') opts.ignoreUrlStr.push(ptn)
      else opts.ignoreUrlRegExp.push(ptn)
    }
    delete opts.ignoreUrls
  }

  if (opts.ignoreUserAgents) {
    for (const ptn of opts.ignoreUserAgents) {
      if (typeof ptn === 'string') opts.ignoreUserAgentStr.push(ptn)
      else opts.ignoreUserAgentRegExp.push(ptn)
    }
    delete opts.ignoreUserAgents
  }

  if (opts.ignoreMessageQueues) {
    if (!opts.ignoreMessageQueuesRegExp) {
      opts.ignoreMessageQueuesRegExp = []
    }
    const wildcard = new WildcardMatcher()
    for (const ptn of opts.ignoreMessageQueues) {
      const re = wildcard.compile(ptn)
      opts.ignoreMessageQueuesRegExp.push(re)
    }
  }
}

function normalizeElasticsearchCaptureBodyUrls (opts) {
  if (opts.elasticsearchCaptureBodyUrls) {
    const wildcard = new WildcardMatcher()
    for (const ptn of opts.elasticsearchCaptureBodyUrls) {
      const re = wildcard.compile(ptn)
      opts.elasticsearchCaptureBodyUrlsRegExp.push(re)
    }
  }
}

function normalizeNumbers (opts) {
  for (const key of NUM_OPTS) {
    if (key in opts) opts[key] = Number(opts[key])
  }

  for (const key of MINUS_ONE_EQUAL_INFINITY) {
    if (opts[key] === -1) opts[key] = Infinity
  }
}

function normalizeBytes (opts) {
  for (const key of BYTES_OPTS) {
    if (key in opts) opts[key] = bytes(String(opts[key]))
  }
}

function normalizeDurationOptions (opts, logger) {
  for (const optSpec of DURATION_OPTS) {
    const key = optSpec.name
    if (key in opts) {
      const val = secondsFromDuration(opts[key], optSpec.defaultUnit,
        optSpec.allowedUnits, optSpec.allowNegative)
      if (val === null) {
        if (key in DEFAULTS) {
          const def = DEFAULTS[key]
          logger.warn('invalid duration value "%s" for "%s" config option: using default "%s"',
            opts[key], key, def)
          opts[key] = secondsFromDuration(def, optSpec.defaultUnit,
            optSpec.allowedUnits, optSpec.allowNegative)
        } else {
          logger.warn('invalid duration value "%s" for "%s" config option: ignoring this option',
            opts[key], key)
          delete opts[key]
        }
      } else {
        opts[key] = val
      }
    }
  }
}

// Array config vars are either already an array of strings, or a
// comma-separated string (whitespace is trimmed):
//    'foo, bar' => ['foo', 'bar']
function normalizeArrays (opts) {
  for (const key of ARRAY_OPTS) {
    if (key in opts && typeof opts[key] === 'string') {
      opts[key] = opts[key].split(',').map(v => v.trim())
    }
  }
}

// KeyValuePairs config vars are either an object or a comma-separated string
// of key=value pairs (whitespace around the "key=value" strings is trimmed):
//    {'foo': 'bar', 'eggs': 'spam'} => [['foo', 'bar'], ['eggs', 'spam']]
//    foo=bar, eggs=spam             => [['foo', 'bar'], ['eggs', 'spam']]
function normalizeKeyValuePairs (opts) {
  for (const key of KEY_VALUE_OPTS) {
    if (key in opts) {
      if (typeof opts[key] === 'object' && !Array.isArray(opts[key])) {
        opts[key] = Object.entries(opts[key])
        return
      }

      if (!Array.isArray(opts[key]) && typeof opts[key] === 'string') {
        opts[key] = opts[key].split(',').map(v => v.trim())
      }

      if (Array.isArray(opts[key])) {
        // Note: Currently this assumes no '=' in the value. Also this does not
        // trim whitespace.
        opts[key] = opts[key].map(
          value => typeof value === 'string' ? value.split('=') : value)
      }
    }
  }
}

function normalizeBools (opts, logger) {
  for (const key of BOOL_OPTS) {
    if (key in opts) opts[key] = strictBool(logger, key, opts[key])
  }
}

function truncateOptions (opts) {
  if (opts.serviceVersion) opts.serviceVersion = truncate(String(opts.serviceVersion), INTAKE_STRING_MAX_SIZE)
  if (opts.hostname) opts.hostname = truncate(String(opts.hostname), INTAKE_STRING_MAX_SIZE)
}

// Translate a string byte size, e.g. '10kb', into an integer number of bytes.
function bytes (input) {
  const matches = input.match(/^(\d+)(b|kb|mb|gb)$/i)
  if (!matches) return Number(input)

  const suffix = matches[2].toLowerCase()
  let value = Number(matches[1])

  if (!suffix || suffix === 'b') {
    return value
  }

  value *= 1024
  if (suffix === 'kb') {
    return value
  }

  value *= 1024
  if (suffix === 'mb') {
    return value
  }

  value *= 1024
  if (suffix === 'gb') {
    return value
  }
}

// Convert a given duration config option into a number of seconds.
// If the given duration is invalid, this returns `null`.
// Units are *case-sensitive*.
//
// @param {String|Number} duration - Typically a string of the form `<num><unit>`,
//    for example `30s`, `-1ms`, `2m`. The `defaultUnit` is used if a unit is
//    not part of the string, or if duration is a number. If given as a string,
//    decimal ('1.5s') and exponential-notation ('1e-3s') values are not allowed.
// @param {String} defaultUnit
// @param {Array} allowedUnits - An array of the allowed unit strings. This
//    array may include any number of `us`, `ms`, `s`, and `m`.
// @param {Boolean} allowNegative - Whether a negative number is allowed.
//
// Examples:
//    secondsFromDuration('30s', 's', ['ms', 's', 'm'], false) // => 30
//    secondsFromDuration('-1s', 's', ['ms', 's', 'm'], false) // => null
//    secondsFromDuration('-1ms', 's', ['ms', 's', 'm'], true) // => -0.001
//    secondsFromDuration(500, 'ms', ['us', 'ms', 's', 'm'], false) // => 0.5
//
function secondsFromDuration (duration, defaultUnit, allowedUnits, allowNegative) {
  let val
  let unit
  if (typeof duration === 'string') {
    let match
    if (allowNegative) {
      match = /^(-?\d+)(\w+)?$/.exec(duration)
    } else {
      match = /^(\d+)(\w+)?$/.exec(duration)
    }
    if (!match) {
      return null
    }
    val = Number(match[1])
    if (isNaN(val) || !Number.isFinite(val)) {
      return null
    }
    unit = match[2] || defaultUnit
    if (!allowedUnits.includes(unit)) {
      return null
    }
  } else if (typeof duration === 'number') {
    if (isNaN(duration)) {
      return null
    } else if (duration < 0 && !allowNegative) {
      return null
    }
    val = duration
    unit = defaultUnit
  } else {
    return null
  }

  // Scale to seconds.
  switch (unit) {
    case 'us':
      val /= 1e6
      break
    case 'ms':
      val /= 1e3
      break
    case 's':
      break
    case 'm':
      val *= 60
      break
    default:
      throw new Error(`unknown unit "${unit}" from "${duration}"`)
  }

  return val
}

function strictBool (logger, key, value) {
  if (typeof value === 'boolean') {
    return value
  }
  // This will return undefined for unknown inputs, resulting in them being skipped.
  switch (value) {
    case 'false': return false
    case 'true': return true
    default: {
      logger.warn('unrecognized boolean value "%s" for "%s"', value, key)
    }
  }
}

function maybePairsToObject (pairs) {
  return pairs ? pairsToObject(pairs) : undefined
}

function pairsToObject (pairs) {
  return pairs.reduce((object, [key, value]) => {
    object[key] = value
    return object
  }, {})
}

function loadConfigFile (configFile) {
  const confPath = path.resolve(configFile || process.env.ELASTIC_APM_CONFIG_FILE || 'elastic-apm-node.js')

  if (fs.existsSync(confPath)) {
    try {
      return require(confPath)
    } catch (err) {
      console.error('Elastic APM initialization error: Can\'t read config file %s', confPath)
      console.error(err.stack)
    }
  }

  return null
}

function loadServerCaCertFile (opts) {
  if (opts.serverCaCertFile) {
    try {
      return fs.readFileSync(opts.serverCaCertFile)
    } catch (err) {
      opts.logger.error('Elastic APM initialization error: Can\'t read server CA cert file %s (%s)', opts.serverCaCertFile, err.message)
    }
  }
}

// Return the User-Agent string the agent will use for its comms to APM Server.
//
// Per https://github.com/elastic/apm/blob/main/specs/agents/transport.md#user-agent
// the pattern is roughly this:
//    $repoName/$version ($serviceName $serviceVersion)
//
// The format of User-Agent is governed by https://datatracker.ietf.org/doc/html/rfc7231.
//    User-Agent = product *( RWS ( product / comment ) )
// We do not expect `$repoName` and `$version` to have surprise/invalid values.
// From `validateServiceName` above, we know that `$serviceName` is null or a
// string limited to `/^[a-zA-Z0-9 _-]+$/`. However, `$serviceVersion` is
// provided by the user and could have invalid characters.
//
// `comment` is defined by
// https://datatracker.ietf.org/doc/html/rfc7230#section-3.2.6 as:
//    comment        = "(" *( ctext / quoted-pair / comment ) ")"
//    obs-text       = %x80-FF
//    ctext          = HTAB / SP / %x21-27 / %x2A-5B / %x5D-7E / obs-text
//    quoted-pair    = "\" ( HTAB / SP / VCHAR / obs-text )
//
// `commentBadChar` below *approximates* these rules, and is used to replace
// invalid characters with '_' in the generated User-Agent string. This
// replacement isn't part of the APM spec.
function userAgentFromConf (conf) {
  let userAgent = `apm-agent-nodejs/${version}`

  // This regex *approximately* matches the allowed syntax for a "comment".
  // It does not handle "quoted-pair" or a "comment" in a comment.
  const commentBadChar = /[^\t \x21-\x27\x2a-\x5b\x5d-\x7e\x80-\xff]/g
  const commentParts = []
  if (conf.serviceName) {
    commentParts.push(conf.serviceName)
  }
  if (conf.serviceVersion) {
    commentParts.push(conf.serviceVersion.replace(commentBadChar, '_'))
  }
  if (commentParts.length > 0) {
    userAgent += ` (${commentParts.join(' ')})`
  }

  return userAgent
}

function getBaseClientConfig (conf, agent) {
  let clientLogger = null
  if (!logging.isLoggerCustom(agent.logger)) {
    // https://www.elastic.co/guide/en/ecs/current/ecs-event.html#field-event-module
    clientLogger = agent.logger.child({ 'event.module': 'apmclient' })
  }

  const clientConfig = {
    agentName: 'nodejs',
    agentVersion: version,
    serviceName: conf.serviceName,
    serviceNodeName: conf.serviceNodeName,
    serviceVersion: conf.serviceVersion,
    frameworkName: conf.frameworkName,
    frameworkVersion: conf.frameworkVersion,
    globalLabels: maybePairsToObject(conf.globalLabels),
    hostname: conf.hostname,
    environment: conf.environment,

    // Sanitize conf
    truncateKeywordsAt: INTAKE_STRING_MAX_SIZE,
    truncateLongFieldsAt: conf.longFieldMaxLength,
    // truncateErrorMessagesAt: see below

    // HTTP conf
    secretToken: conf.secretToken,
    apiKey: conf.apiKey,
    userAgent: userAgentFromConf(conf),
    serverUrl: conf.serverUrl,
    serverCaCert: loadServerCaCertFile(conf),
    rejectUnauthorized: conf.verifyServerCert,
    serverTimeout: conf.serverTimeout * 1000,

    // APM Agent Configuration via Kibana:
    centralConfig: conf.centralConfig,

    // Streaming conf
    size: conf.apiRequestSize,
    time: conf.apiRequestTime * 1000,
    maxQueueSize: conf.maxQueueSize,

    // Debugging/testing options
    logger: clientLogger,
    payloadLogFile: conf.payloadLogFile,
    apmServerVersion: conf.apmServerVersion,

    // Container conf
    containerId: conf.containerId,
    kubernetesNodeName: conf.kubernetesNodeName,
    kubernetesNamespace: conf.kubernetesNamespace,
    kubernetesPodName: conf.kubernetesPodName,
    kubernetesPodUID: conf.kubernetesPodUID
  }

  // Metadata handling.
  if (isLambdaExecutionEnvironment()) {
    // Tell the Client to wait for a subsequent `.setExtraMetadata()` call
    // before allowing intake requests. This will be called by `apm.lambda()`
    // on first Lambda function invocation.
    clientConfig.expectExtraMetadata = true
  } else if (conf.cloudProvider !== 'none') {
    clientConfig.cloudMetadataFetcher = new CloudMetadata(conf.cloudProvider, conf.logger, conf.serviceName)
  }

  if (conf.errorMessageMaxLength !== undefined) {
    // As of v10 of the http client, truncation of error messages will default
    // to `truncateLongFieldsAt` if `truncateErrorMessagesAt` is not specified.
    clientConfig.truncateErrorMessagesAt = conf.errorMessageMaxLength
  }

  return clientConfig
}

// Exports.
module.exports = {
  configLogger,
  initialConfig,
  createConfig,
  INTAKE_STRING_MAX_SIZE,
  CAPTURE_ERROR_LOG_STACK_TRACES_NEVER,
  CAPTURE_ERROR_LOG_STACK_TRACES_MESSAGES,
  CAPTURE_ERROR_LOG_STACK_TRACES_ALWAYS,
  CONTEXT_MANAGER_PATCH,
  CONTEXT_MANAGER_ASYNCHOOKS,
  CONTEXT_MANAGER_ASYNCLOCALSTORAGE,
  TRACE_CONTINUATION_STRATEGY_CONTINUE,
  TRACE_CONTINUATION_STRATEGY_RESTART,
  TRACE_CONTINUATION_STRATEGY_RESTART_EXTERNAL,

  // The following are exported for tests.
  DEFAULTS,
  DURATION_OPTS,
  secondsFromDuration,
  userAgentFromConf,
  ENV_TABLE
}
