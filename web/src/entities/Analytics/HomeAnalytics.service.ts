import AnalyticsService, { Categories, Labels } from "./Analytics.service";

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

const {event} = AnalyticsService(Categories.Home);

const HomeAnalyticsService = (): THomeAnalytics => {  
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

export default HomeAnalyticsService();
