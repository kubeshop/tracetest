import {TAssertionResults} from '../types/Assertion.types';

const TraceService = () => ({
  getTestResultCount({resultList}: TAssertionResults = {resultList: [], allPassed: false}) {
    const [totalPassedCount, totalFailedCount] = resultList.reduce<[number, number]>(
      ([innerTotalPassedCount, innerTotalFailedCount], {resultList: testResultList}) => {
        const [passed, failed] = testResultList.reduce<[number, number]>(
          ([passedResultCount, failedResultCount], {spanResults}) => {
            const passedCount = spanResults.filter(({passed: hasPassed}) => hasPassed).length;
            const failedCount = spanResults.filter(({passed: hasPassed}) => !hasPassed).length;

            return [passedResultCount + passedCount, failedResultCount + failedCount];
          },
          [0, 0]
        );

        return [innerTotalPassedCount + passed, innerTotalFailedCount + failed];
      },
      [0, 0]
    );

    return {totalFailedCount, totalPassedCount};
  },
});

export default TraceService();
