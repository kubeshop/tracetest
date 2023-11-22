import {Model} from 'types/Common.types';

export type TRawRunMetadata = Record<string, string>;

export enum KnownSources {
  WEB = 'web',
  CLI = 'cli',
  API = 'api',
  K6 = 'k6',
  UNKNOWN = 'unknown',
}

type TKnownSources = KnownSources | string;

type RunMetadata = Model<
  Record<string, string>,
  {
    name: string;
    buildNumber: string;
    branch: string;
    url: string;
    source: TKnownSources;
  }
>;

function RunMetadata({
  name = '',
  buildNumber = '',
  branch = '',
  url = '',
  source = '',
  ...rest
}: TRawRunMetadata): RunMetadata {
  return {
    name,
    buildNumber,
    branch,
    url,
    source,
    ...rest,
  };
}

export default RunMetadata;
