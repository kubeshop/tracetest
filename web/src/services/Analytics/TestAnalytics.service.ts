import {Categories} from '../../constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  RunTest = 'run-test-button-click',
  TestRunClick = 'test-run-result-click',
}

type TTestAnalytics = {
  onRunTest(testId: string): void;
  onTestRunClick(testRunId: string): void;
};

const TestAnalyticsService = (): TTestAnalytics => {
  const onRunTest = (testId: string) => {
    AnalyticsService.event(Categories.Test, Actions.RunTest, testId);
  };

  const onTestRunClick = (testRunId: string) => {
    AnalyticsService.event(Categories.Test, Actions.TestRunClick, testRunId);
  };

  return {
    onRunTest,
    onTestRunClick,
  };
};

export default TestAnalyticsService();
