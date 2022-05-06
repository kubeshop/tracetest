export const BASE_URL = 'http://localhost:8080';
export const SENTRY_DNS = 'https://8411cbb3b7d84c879f711f0e642a28e3@o1229268.ingest.sentry.io/6375361';

export const SENTRY_ALLOWED_URLS = [/.*?localhost:3000/, /.*?tracetest.io/];

export const DOCUMENTATION_URL = 'https://kubeshop.github.io/tracetest/';
export const GITHUB_URL = 'https://github.com/kubeshop/tracetest';
export const GITHUB_ISSUES_URL = 'https://github.com/kubeshop/tracetest/issues/new/choose';
export const DISCORD_URL = 'https://discord.gg/6zupCZFQbe';

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
