import TestAnalyticsService, {Actions} from '../TestAnalytics.service';
import AnalyticsService from '../Analytics.service';
import {Categories} from '../../../constants/Analytics.constants';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('TestAnalyticsService', () => {
  it('should trigger the onRunTest event', () => {
    TestAnalyticsService.onRunTest('testId');

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Test, Actions.RunTest, 'testId');
  });

  it('should trigger the onTestRunClick event', () => {
    TestAnalyticsService.onTestRunClick('testRunId');

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Test, Actions.TestRunClick, 'testRunId');
  });
});
