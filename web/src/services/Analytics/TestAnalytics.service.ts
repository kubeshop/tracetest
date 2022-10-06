import {Categories, Labels} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  RunTest = 'run-test-button-click',
  TestRunClick = 'test-run-result-click',
  TestCardCollapse = 'test-card-collapse',
  DeleteTest = 'delete-test-button-click',
  DeleteTestRun = 'delete-test-run-button-click',
  DisplayTestInfo = 'display-test-info-button-hover',
}

const TestAnalyticsService = () => ({
  onRunTest() {
    AnalyticsService.event(Categories.Test, Actions.RunTest, Labels.Button);
  },
  onTestRunClick() {
    AnalyticsService.event(Categories.Test, Actions.TestRunClick, Labels.Button);
  },
  onTestCardCollapse() {
    AnalyticsService.event(Categories.Home, Actions.TestCardCollapse, Labels.Button);
  },
  onDeleteTest() {
    AnalyticsService.event(Categories.Home, Actions.DeleteTest, Labels.Button);
  },
  onDeleteTestRun() {
    AnalyticsService.event(Categories.Test, Actions.DeleteTestRun, Labels.Button);
  },
  onDisplayTestInfo() {
    AnalyticsService.event(Categories.TestRun, Actions.DisplayTestInfo, Labels.Button);
  },
});

export default TestAnalyticsService();
