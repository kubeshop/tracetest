import AssertionResultsMock from '../../models/__mocks__/AssertionResults.mock';
import TraceService from '../Trace.service';

describe('TraceService', () => {
  describe('getTestResultCount', () => {
    it('should return zeroed result', () => {
      const testResultCount = TraceService.getTestResultCount();
      expect(testResultCount).toEqual({totalFailedCount: 0, totalPassedCount: 0});
    });

    it('should return full result', () => {
      const assertionResults = AssertionResultsMock.model();
      const testResultCount = TraceService.getTestResultCount(assertionResults);
      expect(testResultCount.totalFailedCount).toBeGreaterThan(0);
      expect(testResultCount.totalPassedCount).toBeGreaterThan(0);
    });
  });
});
