export enum TestState {
  CREATED = 'CREATED',
  EXECUTING = 'EXECUTING',
  AWAITING_TRACE = 'AWAITING_TRACE',
  AWAITING_TEST_RESULTS = 'AWAITING_TEST_RESULTS',
  FAILED = 'FAILED',
  FINISHED = 'FINISHED',
}

export const TestStateMap: Record<
  TestState,
  {status: 'success' | 'processing' | 'error' | 'default' | 'warning'; label: string; percent?: number}
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
  [TestState.FINISHED]: {
    status: 'success',
    label: 'Finished',
  },
};
