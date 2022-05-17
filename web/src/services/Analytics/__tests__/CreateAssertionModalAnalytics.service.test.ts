import CreateAssertionModalAnalyticsService, {Actions} from '../CreateAssertionModalAnalytics.service';
import AnalyticsService from '../Analytics.service';
import {Categories, Labels} from '../../../constants/Analytics.constants';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('CreateAssertionModalAnalyticsService', () => {
  it('should trigger the onCreateAssertionFormSubmit event', () => {
    CreateAssertionModalAnalyticsService.onCreateAssertionFormSubmit('testId');

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.Assertion,
      Actions.CreateAssertionFormSubmit,
      'testId'
    );
  });

  it('should trigger the onEditAssertionFormSubmit event', () => {
    CreateAssertionModalAnalyticsService.onEditAssertionFormSubmit('assertionId');

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.Assertion,
      Actions.EditAssertionFormSubmit,
      'assertionId'
    );
  });

  it('should trigger the onSelectorChange event', () => {
    CreateAssertionModalAnalyticsService.onSelectorChange('selector');

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Assertion, Actions.SelectorChange, 'selector');
  });

  it('should trigger the onChecksChange event', () => {
    CreateAssertionModalAnalyticsService.onChecksChange('checks');

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Assertion, Actions.ChecksChange, 'checks');
  });

  it('should trigger the onAddCheck event', () => {
    CreateAssertionModalAnalyticsService.onAddCheck();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Assertion, Actions.AddCheck, Labels.Button);
  });

  it('should trigger the onRemoveCheck event', () => {
    CreateAssertionModalAnalyticsService.onRemoveCheck();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Assertion, Actions.RemoveCheck, Labels.Button);
  });
});
