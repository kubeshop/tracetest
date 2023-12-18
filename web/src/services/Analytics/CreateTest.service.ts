import {Categories, Labels} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

enum Actions {
  FormSubmit = 'create-test-form-submit',
  DemoClick = 'create-test-demo-click',
  TriggerTypeSelect = 'create-test-trigger-type-select',
  ImportTypeSelect = 'create-test-import-type-select',
}

type TCreateTestAnalytics = {
  onFormSubmit(): void;
  onDemoClick(): void;
  onTriggerSelect(triggerType: string): void;
  onImportSelect(importType: string): void;
};

const CreateTestAnalytics = (): TCreateTestAnalytics => ({
  onFormSubmit() {
    AnalyticsService.event(Categories.Home, Actions.FormSubmit, Labels.Form);
  },
  onDemoClick() {
    AnalyticsService.event(Categories.Home, Actions.DemoClick, Labels.Button);
  },
  onTriggerSelect(triggerType: string) {
    AnalyticsService.event(Categories.Home, Actions.TriggerTypeSelect, triggerType);
  },
  onImportSelect(importType: string) {
    AnalyticsService.event(Categories.Home, Actions.ImportTypeSelect, importType);
  },
});

export default CreateTestAnalytics();
