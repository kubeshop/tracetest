import {Categories, Labels} from 'constants/Analytics.constants';
import CreateAssertionModalAnalyticsService, {Actions} from '../CreateAssertionModalAnalytics.service';
import AnalyticsService from '../Analytics.service';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('CreateAssertionModalAnalyticsService', () => {
  it('should trigger the onCreateAssertionFormSubmit event', () => {
    CreateAssertionModalAnalyticsService.onCreateAssertionFormSubmit();

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.Assertion,
      Actions.CreateAssertionFormSubmit,
      Labels.Form
    );
  });

  it('should trigger the onEditAssertionFormSubmit event', () => {
    CreateAssertionModalAnalyticsService.onEditAssertionFormSubmit();

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.Assertion,
      Actions.EditAssertionFormSubmit,
      Labels.Form
    );
  });

  it('should trigger the onSelectorChange event', () => {
    CreateAssertionModalAnalyticsService.onSelectorChange();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Assertion, Actions.SelectorChange, Labels.Input);
  });

  it('should trigger the onChecksChange event', () => {
    CreateAssertionModalAnalyticsService.onChecksChange();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Assertion, Actions.ChecksChange, Labels.Input);
  });

  it('should trigger the onAddCheck event', () => {
    CreateAssertionModalAnalyticsService.onAddCheck();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Assertion, Actions.AddCheck, Labels.Button);
  });

  it('should trigger the onRemoveCheck event', () => {
    CreateAssertionModalAnalyticsService.onRemoveCheck();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Assertion, Actions.RemoveCheck, Labels.Button);
  });

  it('should trigger the onAssertionFormOpen event', () => {
    CreateAssertionModalAnalyticsService.onAssertionFormOpen();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Assertion, Actions.OpenForm, Labels.Button);
  });

  it('should trigger the onConfirmationModalOpen event', () => {
    CreateAssertionModalAnalyticsService.onConfirmationModalOpen();

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.Assertion,
      Actions.ConfirmationModalOpen,
      Labels.Button
    );
  });
});
