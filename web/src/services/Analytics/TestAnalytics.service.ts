import AnalyticsService, {Categories} from './Analytics.service';

enum Actions {
  RunTest = 'run-test-button-click',
  TestRunClick = 'test-run-click',
}

type TTestAnalytics = {
  onRunTest(testId: string): void;
  onTestRunClick(testRunId: string): void;
};

const {event} = AnalyticsService(Categories.Test);

const TestAnalyticsService = (): TTestAnalytics => {
  const onRunTest = (testId: string) => {
    event(Actions.RunTest, testId);
  };

  const onTestRunClick = (testRunId: string) => {
    event(Actions.TestRunClick, testRunId);
  };

  return {
    onRunTest,
    onTestRunClick,
  };
};

export default TestAnalyticsService();
