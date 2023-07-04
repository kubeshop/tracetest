export const SENTRY_DNS = 'https://8411cbb3b7d84c879f711f0e642a28e3@o1229268.ingest.sentry.io/6375361';
export const SENTRY_ALLOWED_URLS = [/.*?localhost:3000/, /.*?tracetest.io/];

export const DOCUMENT_TITLE = 'Tracetest';

export const DOCUMENTATION_URL = 'https://docs.tracetest.io';
export const GITHUB_URL = 'https://github.com/kubeshop/tracetest';
export const GITHUB_ISSUES_URL = 'https://github.com/kubeshop/tracetest/issues/new/choose';
export const DISCORD_URL = 'https://discord.gg/6zupCZFQbe';
export const OCTOLIINT_ISSUE_URL = 'https://github.com/kubeshop/tracetest/issues/2615';
export const CLI_RUNNING_TESTS_URL = 'https://docs.tracetest.io/cli/running-tests';

export const INGESTOR_ENDPOINT_URL = 'https://docs.tracetest.io/configuration/ingestor-endpoint';

export const TRACE_SEMANTIC_CONVENTIONS_URL =
  'https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/trace/semantic_conventions';
export const RESOURCE_SEMANTIC_CONVENTIONS_URL =
  'https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/resource/semantic_conventions';
export const TRACE_DOCUMENTATION_URL =
  'https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/api.md';

export const ADD_TEST_URL = 'https://docs.tracetest.io/web-ui/creating-tests';
export const ADD_TEST_OUTPUTS_DOCUMENTATION_URL = 'https://docs.tracetest.io/web-ui/creating-test-outputs';
export const ANALYZER_DOCUMENTATION_URL = 'https://docs.tracetest.io/concepts/tracetest-analyzer-concepts';
export const EXPRESSIONS_DOCUMENTATION_URL = 'https://docs.tracetest.io/concepts/expressions';
export const ENVIRONMENTS_DOCUMENTATION_URL = 'https://docs.tracetest.io/concepts/environments';

export const SELECTOR_LANGUAGE_CHEAT_SHEET_URL = `${process.env.PUBLIC_URL}/SL_cheat_sheet.pdf`;

export const POKESHOP_GITHUB = 'https://github.com/kubeshop/pokeshop';
export const OTEL_DEMO_GITHUB = 'https://github.com/open-telemetry/opentelemetry-demo';

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
  CURL = 'CURL',
  Messaging = 'Messaging',
  GRPC = 'GRPC',
  Postman = 'Postman',
  OpenAPI = 'OpenAPI',
  TraceID = 'TraceID',
}
