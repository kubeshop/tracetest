import {Categories, Labels} from '../../constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  CreateTestFormSubmit = 'create-test-form-submit',
}

type TCreateTestAnalytics = {
  onCreateTestFormSubmit(): void;
};

const CreateTestAnalyticsService = (): TCreateTestAnalytics => {
  const onCreateTestFormSubmit = () => {
    AnalyticsService.event(Categories.Home, Actions.CreateTestFormSubmit, Labels.Form);
  };

  return {
    onCreateTestFormSubmit,
  };
};

export default CreateTestAnalyticsService();
