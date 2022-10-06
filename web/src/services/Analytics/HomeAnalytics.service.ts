import {Categories, Labels} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  CreateTestClick = 'create-test-button-click',
  GuidedTourClick = 'guided-tour-click',
  TestClick = 'test-click',
}

type THomeAnalytics = {
  onCreateTestClick(): void;
  onGuidedTourClick(): void;
  onTestClick(testId: string): void;
};

const HomeAnalyticsService = (): THomeAnalytics => ({
  onCreateTestClick() {
    AnalyticsService.event(Categories.Home, Actions.CreateTestClick, Labels.Button);
  },
  onGuidedTourClick() {
    AnalyticsService.event(Categories.Home, Actions.GuidedTourClick, Labels.Button);
  },
  onTestClick(testId: string) {
    AnalyticsService.event(Categories.Home, Actions.TestClick, testId);
  },
});

export default HomeAnalyticsService();
