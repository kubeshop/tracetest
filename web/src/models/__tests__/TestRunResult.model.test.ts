import faker from '@faker-js/faker';
import TestRunResult from '../TestRunResult.model';
import TestRunResultMock from '../__mocks__/TestRunResult.mock';

describe('Test Run Result', () => {
  it('should generate a test run result object', () => {
    const rawTestRunResult = TestRunResultMock.raw();
    const testRunResult = TestRunResult(rawTestRunResult);

    expect(testRunResult.resultId).toEqual(rawTestRunResult.resultId);
    expect(testRunResult.trace).not.toEqual(undefined);
    expect(testRunResult.totalAssertionCount).toEqual(0);
    expect(testRunResult.passedAssertionCount).toEqual(0);
    expect(testRunResult.failedAssertionCount).toEqual(0);
  });

  it('should generate a test run result with assertion count calculation', () => {
    const rawTestRunResult = TestRunResultMock.raw({
      assertionResult: faker.datatype.array(faker.datatype.number({min: 2, max: 5})).map(() => ({
        assertionId: faker.datatype.uuid(),
        spanAssertionResults: faker.datatype.array(faker.datatype.number({min: 2, max: 5})).map(() => ({
          spanAssertionId: faker.datatype.uuid(),
          spanId: faker.datatype.uuid(),
          passed: faker.datatype.boolean(),
          observedValue: faker.random.word(),
        })),
      })),
    });

    const testRunResult = TestRunResult(rawTestRunResult);

    expect(testRunResult.totalAssertionCount).not.toEqual(0);
    expect(testRunResult.passedAssertionCount).not.toEqual(0);
    expect(testRunResult.failedAssertionCount).not.toEqual(0);
  });

  it('should handle a non finished result', () => {
    const rawTestRunResult = TestRunResultMock.raw({
      completedAt: undefined,
      trace: undefined,
    });

    const testRunResult = TestRunResult(rawTestRunResult);

    expect(testRunResult.trace).toEqual(undefined);
    expect(testRunResult.executionTime).toEqual(0);
  });
});
