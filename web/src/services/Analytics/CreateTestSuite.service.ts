import {Categories, Labels} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

enum Actions {
  FormSubmit = 'create-test-suite-form-submit',
  NextClick = 'create-test-suite-next-click',
  PrevClick = 'create-test-suite-prev-click',
}

type TCreateTestSuiteAnalytics = {
  onFormSubmit(): void;
  onNextClick(stepName: string): void;
  onPrevClick(stepName: string): void;
};

const CreateTestSuiteAnalytics = (): TCreateTestSuiteAnalytics => ({
  onFormSubmit() {
    AnalyticsService.event(Categories.Home, Actions.FormSubmit, Labels.Form);
  },
  onNextClick(stepName) {
    AnalyticsService.event(Categories.Home, Actions.NextClick, stepName);
  },
  onPrevClick(stepName) {
    AnalyticsService.event(Categories.Home, Actions.PrevClick, stepName);
  },
});

export default CreateTestSuiteAnalytics();
