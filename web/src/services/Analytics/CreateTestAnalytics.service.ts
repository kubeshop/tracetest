import {Categories, Labels} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  CreateTestFormSubmit = 'create-test-form-submit',
  DemoTestClick = 'demo-test-click',
  CreateTestPluginSelected = 'create-test-plugin-selected',
  CreateTestNextClick = 'create-test-next-click',
  CreateTestPrevClick = 'create-test-prev-click',
}

const CreateTestAnalyticsService = () => ({
  onCreateTestFormSubmit() {
    AnalyticsService.event(Categories.Home, Actions.CreateTestFormSubmit, Labels.Form);
  },
  onDemoTestClick() {
    AnalyticsService.event(Categories.Home, Actions.DemoTestClick, Labels.Button);
  },
  onPluginSelected(pluginName: string) {
    AnalyticsService.event(Categories.Home, Actions.CreateTestPluginSelected, pluginName);
  },
  onNextClick(stepName: string) {
    AnalyticsService.event(Categories.Home, Actions.CreateTestNextClick, stepName);
  },
  onPrevClick(stepName: string) {
    AnalyticsService.event(Categories.Home, Actions.CreateTestPrevClick, stepName);
  },
});

export default CreateTestAnalyticsService();
