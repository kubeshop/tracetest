import {Categories, Labels} from '../../constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  EditAssertionFormSubmit = 'edit-assertion-form-submit',
  CreateAssertionFormSubmit = 'create-assertion-form-submit',
  SelectorChange = 'create-assertion-modal-selector-change',
  ChecksChange = 'create-assertion-modal-assertion-checks-change',
  AddCheck = 'create-assertion-modal-add-check',
  RemoveCheck = 'create-assertion-modal-remove-check',
}

type TCreateAssertionModalAnalytics = {
  onEditAssertionFormSubmit(assertionId: string): void;
  onCreateAssertionFormSubmit(testId: string): void;
  onSelectorChange(selector: string): void;
  onChecksChange(checks: string): void;
  onAddCheck(): void;
  onRemoveCheck(): void;
};

const CreateAssertionModalAnalyticsService = (): TCreateAssertionModalAnalytics => {
  const onCreateAssertionFormSubmit = (testId: string) => {
    AnalyticsService.event(Categories.Assertion, Actions.CreateAssertionFormSubmit, testId);
  };

  const onEditAssertionFormSubmit = (assertionId: string) => {
    AnalyticsService.event(Categories.Assertion, Actions.EditAssertionFormSubmit, assertionId);
  };

  const onSelectorChange = (selector: string) => {
    AnalyticsService.event(Categories.Assertion, Actions.SelectorChange, selector);
  };

  const onChecksChange = (checks: string) => {
    AnalyticsService.event(Categories.Assertion, Actions.ChecksChange, checks);
  };

  const onAddCheck = () => {
    AnalyticsService.event(Categories.Assertion, Actions.AddCheck, Labels.Button);
  };

  const onRemoveCheck = () => {
    AnalyticsService.event(Categories.Assertion, Actions.RemoveCheck, Labels.Button);
  };

  return {
    onCreateAssertionFormSubmit,
    onEditAssertionFormSubmit,
    onSelectorChange,
    onChecksChange,
    onAddCheck,
    onRemoveCheck,
  };
};

export default CreateAssertionModalAnalyticsService();
