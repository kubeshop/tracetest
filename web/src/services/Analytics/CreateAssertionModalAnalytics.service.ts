import {Categories, Labels} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  EditAssertionFormSubmit = 'edit-assertion-form-submit',
  CreateAssertionFormSubmit = 'create-assertion-form-submit',
  SelectorChange = 'create-assertion-modal-selector-change',
  ChecksChange = 'create-assertion-modal-assertion-checks-change',
  AddCheck = 'create-assertion-modal-add-check',
  RemoveCheck = 'create-assertion-modal-remove-check',
  OpenForm = 'open-create-assertion-modal-form',
  ConfirmationModalOpen = 'open-create-assertion-modal-confirmation-modal',
}

const CreateAssertionModalAnalyticsService = () => {
  const onCreateAssertionFormSubmit = () => {
    AnalyticsService.event(Categories.Assertion, Actions.CreateAssertionFormSubmit, Labels.Form);
  };

  const onEditAssertionFormSubmit = () => {
    AnalyticsService.event(Categories.Assertion, Actions.EditAssertionFormSubmit, Labels.Form);
  };

  const onSelectorChange = () => {
    AnalyticsService.event(Categories.Assertion, Actions.SelectorChange, Labels.Input);
  };

  const onChecksChange = () => {
    AnalyticsService.event(Categories.Assertion, Actions.ChecksChange, Labels.Input);
  };

  const onAddCheck = () => {
    AnalyticsService.event(Categories.Assertion, Actions.AddCheck, Labels.Button);
  };

  const onRemoveCheck = () => {
    AnalyticsService.event(Categories.Assertion, Actions.RemoveCheck, Labels.Button);
  };

  const onAssertionFormOpen = () => {
    AnalyticsService.event(Categories.Assertion, Actions.OpenForm, Labels.Button);
  };

  const onConfirmationModalOpen = () => {
    AnalyticsService.event(Categories.Assertion, Actions.ConfirmationModalOpen, Labels.Button);
  };

  return {
    onCreateAssertionFormSubmit,
    onEditAssertionFormSubmit,
    onSelectorChange,
    onChecksChange,
    onAddCheck,
    onRemoveCheck,
    onAssertionFormOpen,
    onConfirmationModalOpen,
  };
};

export default CreateAssertionModalAnalyticsService();
