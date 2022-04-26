import useAnalytics, {Categories, Labels} from '../../components/Analytics/useAnalytics';

enum Actions {
  CreateTestClick = 'create-test-button-click',
  GuidedTourClick = 'guided-tour-click',
  TestClick = 'test-click',
}

type THomeAnalytics = {
  onCreateTestClick(): void;
  onGuidedTourClick(): void;
  onTestClick(testId: string): void;
};

const useHomeAnalytics = (): THomeAnalytics => {
  const {event} = useAnalytics(Categories.Home);

  const onCreateTestClick = () => {
    event(Actions.CreateTestClick, Labels.Button);
  };

  const onGuidedTourClick = () => {
    event(Actions.GuidedTourClick, Labels.Button);
  };

  const onTestClick = (testId: string) => {
    event(Actions.GuidedTourClick, testId);
  };

  return {
    onCreateTestClick,
    onGuidedTourClick,
    onTestClick,
  };
};

export default useHomeAnalytics;
