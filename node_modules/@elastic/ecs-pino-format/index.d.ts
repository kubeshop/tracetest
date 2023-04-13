import type { LoggerOptions } from "pino";

interface Config {
  /**
   * Whether to convert a logged `err` field to ECS error fields.
   * Default true, to match Pino's default of having an `err` serializer.
   */
  convertErr?: boolean;
  /**
   * Whether to convert logged `req` and `res` HTTP request and response fields
   * to ECS HTTP, User agent, and URL fields. Default false.
   */
  convertReqRes?: boolean;
  /**
   * Whether to automatically integrate with
   * Elastic APM (https://github.com/elastic/apm-agent-nodejs). If a started
   * APM agent is detected, then log records will include the following
   * fields:
   *
   * - "service.name" - the configured serviceName in the agent
   * - "event.dataset" - set to "$serviceName.log" for correlation in Kibana
   * - "trace.id", "transaction.id", and "span.id" - if there is a current
   *   active trace when the log call is made
   *
   * Default true.
   */
  apmIntegration?: boolean;
}

declare function createEcsPinoOptions(config?: Config): LoggerOptions;

export = createEcsPinoOptions;
