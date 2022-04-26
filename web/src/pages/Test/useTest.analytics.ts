import useAnalytics, {Categories} from '../../components/Analytics/useAnalytics';

enum Actions {
  RunTest = 'run-test-button-click',
  TestRunClick = 'test-run-click',
}

type TTestAnalytics = {
  onRunTest(testId: string): void;
  onTestRunClick(testRunId: string): void;
};

const useTestAnalytics = (): TTestAnalytics => {
  const {event} = useAnalytics(Categories.Test);

  const onRunTest = (testId: string) => {
    event(Actions.RunTest, testId);
  };

  const onTestRunClick = (testRunId: string) => {
    event(Actions.RunTest, testRunId);
  };

  return {
    onRunTest,
    onTestRunClick,
  };
};

export default useTestAnalytics;
