import {Categories, Labels} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  CreateTestFormSubmit = 'create-test-form-submit',
  DemoTestClick = 'demo-test-click',
}

const CreateTestAnalyticsService = () => {
  const onCreateTestFormSubmit = () => {
    AnalyticsService.event(Categories.Home, Actions.CreateTestFormSubmit, Labels.Form);
  };

  const onDemoTestClick = () => {
    AnalyticsService.event(Categories.Home, Actions.DemoTestClick, Labels.Button);
  };

  return {
    onCreateTestFormSubmit,
    onDemoTestClick,
  };
};

export default CreateTestAnalyticsService();
