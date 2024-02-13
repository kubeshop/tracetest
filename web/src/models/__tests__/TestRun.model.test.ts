import TestRun from '../TestRun.model';
import TestRunResultMock from '../__mocks__/TestRun.mock';

describe('Test Run', () => {
  it('should generate a test run result object', () => {
    const rawTestRunResult = TestRunResultMock.raw();
    const testRunResult = TestRun(rawTestRunResult);

    expect(testRunResult.id).toEqual(rawTestRunResult.id);
    expect(testRunResult.totalAssertionCount).toEqual(0);
    expect(testRunResult.passedAssertionCount).toEqual(0);
    expect(testRunResult.failedAssertionCount).toEqual(0);
  });

  it('should handle a non finished result', () => {
    const rawTestRunResult = TestRunResultMock.raw({
      completedAt: undefined,
      trace: undefined,
    });

    const testRunResult = TestRun(rawTestRunResult);

    expect(testRunResult.executionTime).toEqual(0);
  });
});
