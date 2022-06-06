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

const TestAnalyticsService = () => {
  const onRunTest = () => {
    AnalyticsService.event(Categories.Test, Actions.RunTest, Labels.Button);
  };

  const onTestRunClick = () => {
    AnalyticsService.event(Categories.Test, Actions.TestRunClick, Labels.Button);
  };

  const onTestCardCollapse = () => {
    AnalyticsService.event(Categories.Home, Actions.TestCardCollapse, Labels.Button);
  };

  const onDeleteTest = () => {
    AnalyticsService.event(Categories.Home, Actions.DeleteTest, Labels.Button);
  };

  const onDeleteTestRun = () => {
    AnalyticsService.event(Categories.Test, Actions.DeleteTestRun, Labels.Button);
  };

  const onDisplayTestInfo = () => {
    AnalyticsService.event(Categories.Trace, Actions.DisplayTestInfo, Labels.Button);
  };

  return {
    onRunTest,
    onTestRunClick,
    onTestCardCollapse,
    onDeleteTest,
    onDeleteTestRun,
    onDisplayTestInfo,
  };
};

export default TestAnalyticsService();
