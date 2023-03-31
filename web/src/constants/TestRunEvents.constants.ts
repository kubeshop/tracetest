export enum TestRunStage {
  Trigger = 'trigger',
  Trace = 'trace',
  Test = 'test',
}

export enum PollingInfoType {
  Periodic = 'periodic',
}

export enum OutputInfoLogLevel {
  Warning = 'warning',
  Error = 'error',
}

export enum LogLevel {
  Error = 'Error',
  Info = 'Info',
  Start = 'Start',
  Success = 'Success',
  Warning = 'Warning',
}

export enum TraceEventType {
  DATA_STORE_CONNECTION_INFO = 'DATA_STORE_CONNECTION_INFO',
  POLLING_ITERATION_INFO = 'POLLING_ITERATION_INFO',
}
