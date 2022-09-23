export const BASE_URL = 'http://localhost:11633';
export const SENTRY_DNS = 'https://8411cbb3b7d84c879f711f0e642a28e3@o1229268.ingest.sentry.io/6375361';

export const SENTRY_ALLOWED_URLS = [/.*?localhost:3000/, /.*?tracetest.io/];

export const DOCUMENTATION_URL = 'https://kubeshop.github.io/tracetest/';
export const GITHUB_URL = 'https://github.com/kubeshop/tracetest';
export const GITHUB_ISSUES_URL = 'https://github.com/kubeshop/tracetest/issues/new/choose';
export const DISCORD_URL = 'https://discord.gg/6zupCZFQbe';

export const TRACE_SEMANTIC_CONVENTIONS_URL =
  'https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/trace/semantic_conventions';
export const RESOURCE_SEMANTIC_CONVENTIONS_URL =
  'https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/resource/semantic_conventions';
export const TRACE_DOCUMENTATION_URL =
  'https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/api.md';

export const ADVANCE_SELECTORS_DOCUMENTATION_URL = 'https://kubeshop.github.io/tracetest/advanced-selectors/';

export enum HTTP_METHOD {
  GET = 'GET',
  PUT = 'PUT',
  POST = 'POST',
  PATCH = 'PATCH',
  DELETE = 'DELETE',
  COPY = 'COPY',
  HEAD = 'HEAD',
  OPTIONS = 'OPTIONS',
  LINK = 'LINK',
  UNLINK = 'UNLINK',
  PURGE = 'PURGE',
  LOCK = 'LOCK',
  UNLOCK = 'UNLOCK',
  PROPFIND = 'PROPFIND',
  VIEW = 'VIEW',
}

export const durationRegExp = /(\d+)(ns|μs|ms|s|m|h)/;

export enum RouterSearchFields {
  SelectedAssertion = 'selectedAssertion',
  SelectedSpan = 'selectedSpan',
}

export enum SupportedPlugins {
  REST = 'REST',
  Messaging = 'Messaging',
  GRPC = 'GRPC',
  Postman = 'Postman',
  OpenAPI = 'OpenAPI',
}
