export interface IKeyValue {
  value: string;
  key: string;
}

export const DEFAULT_HEADERS: IKeyValue[] = [{key: 'Content-Type', value: 'application/json'}];

export enum TriggerTypes {
  http = 'http',
  grpc = 'grpc',
  traceid = 'traceid',
  kafka = 'kafka',
  cypress = 'cypress',
  playwright = 'playwright',
}

export enum ImportTypes {
  postman = 'postman',
  curl = 'curl',
  definition = 'definition',
}

export enum SortBy {
  Created = 'created',
  LastRun = 'last_run',
  Name = 'name',
}

export enum SortDirection {
  Asc = 'asc',
  Desc = 'desc',
}

export const sortOptions = [
  {
    label: 'Last Run',
    value: 'last_run',
    params: {sortDirection: SortDirection.Desc, sortBy: SortBy.LastRun},
  },
  {
    label: 'Recently Created',
    value: 'created',
    params: {sortDirection: SortDirection.Desc, sortBy: SortBy.Created},
  },
  {
    label: 'Name, A to Z',
    value: 'name_asc',
    params: {sortDirection: SortDirection.Asc, sortBy: SortBy.Name},
  },
  {
    label: 'Name, Z to A',
    value: 'name_desc',
    params: {sortDirection: SortDirection.Desc, sortBy: SortBy.Name},
  },
] as const;

export enum TracetestApiTags {
  VARIABLE_SET = 'variableSet',
  TESTSUITE = 'testSuite',
  TESTSUITE_RUN = 'testSuiteRun',
  TEST = 'test',
  TEST_DEFINITION = 'testDefinition',
  TEST_RUN = 'testRun',
  SPAN = 'span',
  EXPRESSION = 'expression',
  RESOURCE = 'resource',
  DATA_STORE = 'dataStore',
  WIZARD = 'wizard',
  SETTING = 'setting',
}

export const TracetestApiTagsList = Object.values(TracetestApiTags);
