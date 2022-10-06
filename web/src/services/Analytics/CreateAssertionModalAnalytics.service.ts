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

const CreateAssertionModalAnalyticsService = () => ({
  onCreateAssertionFormSubmit() {
    AnalyticsService.event(Categories.Assertion, Actions.CreateAssertionFormSubmit, Labels.Form);
  },
  onEditAssertionFormSubmit() {
    AnalyticsService.event(Categories.Assertion, Actions.EditAssertionFormSubmit, Labels.Form);
  },
  onSelectorChange() {
    AnalyticsService.event(Categories.Assertion, Actions.SelectorChange, Labels.Input);
  },
  onChecksChange() {
    AnalyticsService.event(Categories.Assertion, Actions.ChecksChange, Labels.Input);
  },
  onAddCheck() {
    AnalyticsService.event(Categories.Assertion, Actions.AddCheck, Labels.Button);
  },
  onRemoveCheck() {
    AnalyticsService.event(Categories.Assertion, Actions.RemoveCheck, Labels.Button);
  },
  onAssertionFormOpen() {
    AnalyticsService.event(Categories.Assertion, Actions.OpenForm, Labels.Button);
  },
  onConfirmationModalOpen() {
    AnalyticsService.event(Categories.Assertion, Actions.ConfirmationModalOpen, Labels.Button);
  },
});

export default CreateAssertionModalAnalyticsService();
