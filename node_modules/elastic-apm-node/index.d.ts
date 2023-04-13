/*
 * Copyright Elasticsearch B.V. and other contributors where applicable.
 * Licensed under the BSD 2-Clause License; you may not use this file except in
 * compliance with the BSD 2-Clause License.
 */

/// <reference types="node" />

// Note: We avoid import of any external `@types/...` to avoid TypeScript users
// needing to manually install them. The only exception is the prerequisite to
// `npm install -D @types/node`.
import type { IncomingMessage, ServerResponse } from 'http';
import { Connect } from './types/connect';
import { AwsLambda } from './types/aws-lambda';

declare namespace apm {
  // Agent API
  // https://www.elastic.co/guide/en/apm/agent/nodejs/current/agent-api.html
  export interface Agent {
    // Configuration
    start (options?: AgentConfigOptions): Agent;
    isStarted (): boolean;
    getServiceName (): string | undefined;
    setFramework (options: {
      name?: string;
      version?: string;
      overwrite?: boolean;
    }): void;
    addPatch (modules: string | Array<string>, handler: string | PatchHandler): void;
    removePatch (modules: string | Array<string>, handler: string | PatchHandler): void;
    clearPatches (modules: string | Array<string>): void;

    // Data collection hooks
    middleware: { connect (): Connect.ErrorHandleFunction };
    lambda (handler: AwsLambda.Handler): AwsLambda.Handler;
    lambda (type: string, handler: AwsLambda.Handler): AwsLambda.Handler;
    handleUncaughtExceptions (
      fn?: (err: Error) => void
    ): void;

    // Errors
    captureError (
      err: Error | string | ParameterizedMessageObject,
      callback?: CaptureErrorCallback
    ): void;
    captureError (
      err: Error | string | ParameterizedMessageObject,
      options?: CaptureErrorOptions,
      callback?: CaptureErrorCallback
    ): void;

    // Distributed Tracing
    currentTraceparent: string | null;
    currentTraceIds: {
      'trace.id'?: string;
      'transaction.id'?: string;
      'span.id'?: string;
    }

    // Transactions
    startTransaction(
      name?: string | null,
      options?: TransactionOptions
    ): Transaction | null;
    startTransaction(
      name: string | null,
      type: string | null,
      options?: TransactionOptions
    ): Transaction | null;
    /**
     * @deprecated Transaction 'subtype' is not used.
     */
    startTransaction(
      name: string | null,
      type: string | null,
      subtype: string | null,
      options?: TransactionOptions
    ): Transaction | null;
    /**
     * @deprecated Transaction 'subtype' and 'action' are not used.
     */
    startTransaction(
      name: string | null,
      type: string | null,
      subtype: string | null,
      action: string | null,
      options?: TransactionOptions
    ): Transaction | null;
    setTransactionName (name: string): void;
    endTransaction (result?: string | number, endTime?: number): void;
    currentTransaction: Transaction | null;

    // Spans
    startSpan(
      name?: string | null,
      options?: SpanOptions
    ): Span | null;
    startSpan(
      name: string | null,
      type: string | null,
      options?: SpanOptions
    ): Span | null;
    startSpan(
      name: string | null,
      type: string | null,
      subtype: string | null,
      options?: SpanOptions
    ): Span | null;
    startSpan(
      name: string | null,
      type: string | null,
      subtype: string | null,
      action: string | null,
      options?: SpanOptions
    ): Span | null;
    currentSpan: Span | null;

    // Context
    setLabel (name: string, value: LabelValue, stringify?: boolean): boolean;
    addLabels (labels: Labels, stringify?: boolean): boolean;
    setUserContext (user: UserObject): void;
    setCustomContext (custom: object): void;

    // Transport
    addFilter (fn: FilterFn): void;
    addErrorFilter (fn: FilterFn): void;
    addSpanFilter (fn: FilterFn): void;
    addTransactionFilter (fn: FilterFn): void;
    addMetadataFilter (fn: FilterFn): void;
    flush (callback?: Function): void;
    destroy (): void;

    // Utils
    logger: Logger;

    // Custom metrics
    registerMetric(name: string, callback: Function): void;
    registerMetric(name: string, labels: Labels, callback: Function): void;

    setTransactionOutcome(outcome: Outcome): void;
    setSpanOutcome(outcome: Outcome): void;
  }

  type Outcome = 'unknown' | 'success' | 'failure';

  // Transaction API
  // https://www.elastic.co/guide/en/apm/agent/nodejs/current/transaction-api.html
  export interface Transaction {
    // The following properties and methods are currently not documented as their API isn't considered official:
    // - timestamp, ended, id, traceId, parentId, sampled, duration()
    // - setUserContext(), setCustomContext(), toJSON(), setDefaultName(), setDefaultNameFromRequest()

    name: string;
    type: string | null;
    /**
     * @deprecated Transaction 'subtype' is not used.
     */
    subtype: string | null;
    /**
     * @deprecated Transaction 'action' is not used.
     */
    action: string | null;
    traceparent: string;
    outcome: Outcome;
    result: string | number;
    ids: {
      'trace.id': string;
      'transaction.id': string;
    }

    setType (type?: string | null, subtype?: string | null, action?: string | null): void;
    setLabel (name: string, value: LabelValue, stringify?: boolean): boolean;
    addLabels (labels: Labels, stringify?: boolean): boolean;
    setOutcome(outcome: Outcome): void;

    startSpan(
      name?: string | null,
      options?: SpanOptions
    ): Span | null;
    startSpan(
      name: string | null,
      type: string | null,
      options?: SpanOptions
    ): Span | null;
    startSpan(
      name: string | null,
      type: string | null,
      subtype: string | null,
      options?: SpanOptions
    ): Span | null;
    startSpan(
      name: string | null,
      type: string | null,
      subtype: string | null,
      action: string | null,
      options?: SpanOptions
    ): Span | null;
    ensureParentId (): string;
    end (result?: string | number | null, endTime?: number): void;
  }

  // Span API
  // https://www.elastic.co/guide/en/apm/agent/nodejs/current/span-api.html
  export interface Span {
    // The following properties and methods are currently not documented as their API isn't considered official:
    // - timestamp, ended, id, traceId, parentId, sampled, duration()
    // - customStackTrace(), setDbContext()

    transaction: Transaction;
    name: string;
    type: string | null;
    subtype: string | null;
    action: string | null;
    traceparent: string;
    outcome: Outcome;
    ids: {
      'trace.id': string;
      'span.id': string;
    }

    setType (type?: string | null, subtype?: string | null, action?: string | null): void;
    setLabel (name: string, value: LabelValue, stringify?: boolean): boolean;
    addLabels (labels: Labels, stringify?: boolean): boolean;
    setOutcome(outcome: Outcome): void;
    setServiceTarget(type?: string | null, name?: string | null): void;
    end (endTime?: number): void;
  }

  // https://www.elastic.co/guide/en/apm/agent/nodejs/current/configuration.html
  export interface AgentConfigOptions {
    abortedErrorThreshold?: string; // Also support `number`, but as we're removing this functionality soon, there's no need to advertise it
    active?: boolean;
    addPatch?: KeyValueConfig;
    apiKey?: string;
    apiRequestSize?: string; // Also support `number`, but as we're removing this functionality soon, there's no need to advertise it
    apiRequestTime?: string; // Also support `number`, but as we're removing this functionality soon, there's no need to advertise it
    asyncHooks?: boolean;
    breakdownMetrics?: boolean;
    captureBody?: CaptureBody;
    captureErrorLogStackTraces?: CaptureErrorLogStackTraces;
    captureExceptions?: boolean;
    captureHeaders?: boolean;
    /**
     * @deprecated Use `spanStackTraceMinDuration`.
     */
    captureSpanStackTraces?: boolean;
    cloudProvider?: string;
    configFile?: string;
    containerId?: string;
    contextPropagationOnly?: boolean;
    disableInstrumentations?: string | string[];
    disableSend?: boolean;
    elasticsearchCaptureBodyUrls?: Array<string>;
    environment?: string;
    errorMessageMaxLength?: string; // DEPRECATED: use `longFieldMaxLength`.
    errorOnAbortedRequests?: boolean;
    exitSpanMinDuration?: string;
    filterHttpHeaders?: boolean;
    frameworkName?: string;
    frameworkVersion?: string;
    globalLabels?: KeyValueConfig;
    hostname?: string;
    ignoreMessageQueues?: Array<string>;
    ignoreUrls?: Array<string | RegExp>;
    ignoreUserAgents?: Array<string | RegExp>;
    instrument?: boolean;
    instrumentIncomingHTTPRequests?: boolean;
    kubernetesNamespace?: string;
    kubernetesNodeName?: string;
    kubernetesPodName?: string;
    kubernetesPodUID?: string;
    logLevel?: LogLevel;
    logUncaughtExceptions?: boolean;
    logger?: Logger; // Notably this Logger interface matches the Pino Logger.
    longFieldMaxLength?: number;
    maxQueueSize?: number;
    metricsInterval?: string; // Also support `number`, but as we're removing this functionality soon, there's no need to advertise it
    metricsLimit?: number;
    payloadLogFile?: string;
    centralConfig?: boolean;
    sanitizeFieldNames?: Array<string>;
    secretToken?: string;
    serverCaCertFile?: string;
    serverTimeout?: string; // Also support `number`, but as we're removing this functionality soon, there's no need to advertise it
    serverUrl?: string;
    serviceName?: string;
    serviceNodeName?: string;
    serviceVersion?: string;
    sourceLinesErrorAppFrames?: number;
    sourceLinesErrorLibraryFrames?: number;
    sourceLinesSpanAppFrames?: number;
    sourceLinesSpanLibraryFrames?: number;
    spanCompressionEnabled?: boolean;
    spanCompressionExactMatchMaxDuration?: string;
    spanCompressionSameKindMaxDuration?: string;
    /**
     * @deprecated Use `spanStackTraceMinDuration`.
     */
    spanFramesMinDuration?: string;
    spanStackTraceMinDuration?: string;
    stackTraceLimit?: number;
    traceContinuationStrategy?: TraceContinuationStrategy;
    transactionIgnoreUrls?: Array<string>;
    transactionMaxSpans?: number;
    transactionSampleRate?: number;
    useElasticTraceparentHeader?: boolean;
    usePathAsTransactionName?: boolean;
    verifyServerCert?: boolean;
  }

  interface CaptureErrorOptions {
    request?: IncomingMessage;
    response?: ServerResponse;
    timestamp?: number;
    handled?: boolean;
    user?: UserObject;
    labels?: Labels;
    tags?: Labels;
    custom?: object;
    message?: string;
    captureAttributes?: boolean;
    skipOutcome?: boolean;
    /**
     * A Transaction or Span instance to make the parent of this error. If not
     * given (undefined), then the current span or transaction will be used. If
     * `null` is given, then no span or transaction will be used.
     */
    parent?: Transaction | Span | null;
  }

  interface Labels {
    [key: string]: LabelValue;
  }

  interface UserObject {
    id?: string | number;
    username?: string;
    email?: string;
  }

  interface ParameterizedMessageObject {
    message: string;
    params: Array<any>;
  }

  interface Logger {
    // Defining overloaded methods rather than a separate `interface LogFn`
    // as @types/pino does, because the IDE completion shows these as *methods*
    // rather than as properties, which is slightly nicer.
    fatal (msg: string, ...args: any[]): void;
    fatal (obj: {}, msg?: string, ...args: any[]): void;
    error (msg: string, ...args: any[]): void;
    error (obj: {}, msg?: string, ...args: any[]): void;
    warn (msg: string, ...args: any[]): void;
    warn (obj: {}, msg?: string, ...args: any[]): void;
    info (msg: string, ...args: any[]): void;
    info (obj: {}, msg?: string, ...args: any[]): void;
    debug (msg: string, ...args: any[]): void;
    debug (obj: {}, msg?: string, ...args: any[]): void;
    trace (msg: string, ...args: any[]): void;
    trace (obj: {}, msg?: string, ...args: any[]): void;
    // Allow a passed in Logger that has other properties, as a Pino logger
    // does. Discussion:
    // https://github.com/elastic/apm-agent-nodejs/pull/926/files#r266239656
    [propName: string]: any;
  }

  // Link and `links` are intended to be compatible with OTel's
  // equivalent APIs in "opentelemetry-js-api/src/trace/link.ts". Currently
  // span link attributes are not supported.
  export interface Link {
    /** A W3C trace-context 'traceparent' string, Transaction, or Span. */
    context: Transaction | Span | string; // This is a SpanContext in OTel.
  }

  export interface TransactionOptions {
    startTime?: number;
    // `childOf` is a W3C trace-context 'traceparent' string. Passing a
    // Transaction or Span is deprecated.
    childOf?: Transaction | Span | string;
    tracestate?: string; // A W3C trace-context 'tracestate' string.
    links?: Link[];
  }

  export interface SpanOptions {
    startTime?: number;
    childOf?: Transaction | Span | string;
    exitSpan?: boolean;
    links?: Link[];
  }

  type CaptureBody = 'off' | 'errors' | 'transactions' | 'all';
  type CaptureErrorLogStackTraces = 'never' | 'messages' | 'always';
  type LogLevel = 'trace' | 'debug' | 'info' | 'warn' | 'warning' | 'error' | 'fatal' | 'critical' | 'off';
  type TraceContinuationStrategy = 'continue' | 'restart' | 'restart_external';

  type CaptureErrorCallback = (err: Error | null, id: string) => void;
  type FilterFn = (payload: Payload) => Payload | boolean | void;
  type LabelValue = string | number | boolean | null | undefined;
  type KeyValueConfig = string | Labels | Array<Array<LabelValue>>

  type Payload = { [propName: string]: any }

  type PatchHandler = (exports: any, agent: Agent, options: PatchOptions) => any;

  interface PatchOptions {
    version: string | undefined;
    enabled: boolean;
  }
}

declare const apm: apm.Agent;
export = apm;
