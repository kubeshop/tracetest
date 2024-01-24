export const SENTRY_DNS = 'https://8411cbb3b7d84c879f711f0e642a28e3@o1229268.ingest.sentry.io/6375361';
export const SENTRY_ALLOWED_URLS = [/.*?localhost:3000/, /.*?tracetest.io/];

export const DOCUMENT_TITLE = 'Tracetest';

export const DOCUMENTATION_URL = 'https://docs.tracetest.io';
export const GITHUB_URL = 'https://github.com/kubeshop/tracetest';
export const GITHUB_ISSUES_URL = 'https://github.com/kubeshop/tracetest/issues/new/choose';
export const COMMUNITY_SLACK_URL = 'https://dub.sh/tracetest-community';
export const OCTOLIINT_ISSUE_URL = 'https://github.com/kubeshop/tracetest/issues/2615';
export const CLI_RUNNING_TESTS_URL = 'https://docs.tracetest.io/cli/running-tests';
export const CLI_RUNNING_TEST_SUITES_URL = 'https://docs.tracetest.io/cli/running-test-suites';

export const INGESTOR_ENDPOINT_URL =
  'https://docs.tracetest.io/configuration/opentelemetry-collector-configuration-file-reference';

export const TRACE_SEMANTIC_CONVENTIONS_URL =
  'https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/trace/semantic_conventions';
export const RESOURCE_SEMANTIC_CONVENTIONS_URL =
  'https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/resource/semantic_conventions';
export const TRACE_DOCUMENTATION_URL =
  'https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/api.md';

export const OPENING_TRACETEST_URL = 'https://docs.tracetest.io/getting-started/open/';
export const ADD_TEST_URL = 'https://docs.tracetest.io/web-ui/creating-tests';
export const ADD_TEST_SUITE_URL = 'https://docs.tracetest.io/web-ui/creating-test-suites';
export const ADD_TEST_OUTPUTS_DOCUMENTATION_URL = 'https://docs.tracetest.io/web-ui/creating-test-outputs';
export const ANALYZER_DOCUMENTATION_URL = 'https://docs.tracetest.io/analyzer/concepts';
export const TEST_RUNNER_DOCUMENTATION_URL = 'https://docs.tracetest.io/configuration/test-runner';
export const ANALYZER_RULES_DOCUMENTATION_URL = 'https://docs.tracetest.io/analyzer/rules';
export const EXPRESSIONS_DOCUMENTATION_URL = 'https://docs.tracetest.io/concepts/expressions';
export const VARIABLE_SET_DOCUMENTATION_URL = 'https://docs.tracetest.io/concepts/variable-sets';
export const GITHUB_ACTION_URL =
  'https://github.com/kubeshop/tracetest-github-action/tree/main/examples/tracetest-cli-with-tracetest-core';
export const CLI_DOCS_URL = 'https://docs.tracetest.io/cli/cli-installation-reference';
export const TRACE_POLLING_DOCUMENTATION_URL = 'https://docs.tracetest.io/configuration/trace-polling/';
export const DEMO_DOCUMENTATION_URL = 'https://docs.tracetest.io/configuration/demo/';

export const SELECTOR_LANGUAGE_CHEAT_SHEET_URL = `${process.env.PUBLIC_URL}/SL_cheat_sheet.pdf`;

export const POKESHOP_GITHUB = 'https://github.com/kubeshop/pokeshop';
export const OTEL_DEMO_GITHUB = 'https://github.com/open-telemetry/opentelemetry-demo';

export const AGENT_DOCS_URL = 'https://docs.tracetest.io/concepts/agent';

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
  GRPC = 'GRPC',
  Kafka = 'Kafka',
  TraceID = 'TraceID',
  Cypress = 'Cypress',
  Playwright = 'Playwright',
}
