export enum TestState {
  CREATED = 'CREATED',
  EXECUTING = 'EXECUTING',
  AWAITING_TRACE = 'AWAITING_TRACE',
  AWAITING_TEST_RESULTS = 'AWAITING_TEST_RESULTS',
  TRIGGER_FAILED = 'TRIGGER_FAILED',
  TRACE_FAILED = 'TRACE_FAILED',
  ASSERTION_FAILED = 'ASSERTION_FAILED',
  ANALYZING_TRACE = 'ANALYZING_TRACE',
  ANALYZING_ERROR = 'ANALYZING_ERROR',
  FAILED = 'FAILED',
  STOPPED = 'STOPPED',
  FINISHED = 'FINISHED',
  WAITING = 'WAITING',
  SKIPPED = 'SKIPPED',
}

export const TestStateMap: Record<
  TestState,
  {
    status: 'success' | 'processing' | 'error' | 'default' | 'warning';
    label: string;
    percent?: number;
    showInfo?: boolean;
  }
> = {
  [TestState.CREATED]: {
    status: 'default',
    label: 'Created',
  },
  [TestState.EXECUTING]: {
    status: 'processing',
    label: 'Running',
    percent: 25,
  },
  [TestState.AWAITING_TRACE]: {
    status: 'warning',
    label: 'Awaiting trace',
    percent: 50,
    showInfo: true,
  },
  [TestState.AWAITING_TEST_RESULTS]: {
    status: 'success',
    label: 'Awaiting test results',
    percent: 75,
  },
  [TestState.FAILED]: {
    status: 'error',
    label: 'Failed',
  },
  [TestState.STOPPED]: {
    status: 'default',
    label: 'Stopped by user',
  },
  [TestState.TRIGGER_FAILED]: {
    status: 'error',
    label: 'Trigger Failed',
  },
  [TestState.TRACE_FAILED]: {
    status: 'error',
    label: 'Trace Failed',
  },
  [TestState.ASSERTION_FAILED]: {
    status: 'error',
    label: 'Test Specs Failed',
  },
  [TestState.FINISHED]: {
    status: 'success',
    label: 'Finished',
  },
  [TestState.WAITING]: {
    status: 'default',
    label: 'Waiting',
  },
  [TestState.SKIPPED]: {
    status: 'warning',
    label: 'Skipped',
  },
  [TestState.ANALYZING_TRACE]: {
    status: 'default',
    label: 'Linting Trace',
  },
  [TestState.ANALYZING_ERROR]: {
    status: 'error',
    label: 'Lint Error',
  },
};

export enum RunDetailModes {
  TRIGGER = 'trigger',
  TRACE = 'trace',
  TEST = 'test',
  AUTOMATE = 'automate',
}
