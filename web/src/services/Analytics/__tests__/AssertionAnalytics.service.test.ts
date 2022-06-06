import {Categories, Labels} from 'constants/Analytics.constants';
import AssertionAnalyticsService, {Actions} from '../AssertionAnalytics.service';
import AnalyticsService from '../Analytics.service';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('AssertionAnalyticsService', () => {
  it('should trigger the onAssertionEdit event', () => {
    AssertionAnalyticsService.onAssertionEdit();

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.TestResults,
      Actions.EditAssertionButtonClick,
      Labels.Button
    );
  });

  it('should trigger the onAssertionClick event', () => {
    AssertionAnalyticsService.onAssertionClick();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.TestResults, Actions.AssertionClick, Labels.Button);
  });

  it('should trigger the onAssertionDelete event', () => {
    AssertionAnalyticsService.onAssertionDelete();

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.TestResults,
      Actions.AssertionDeleteButtonClick,
      Labels.Button
    );
  });

  it('should trigger the onRevertAssertion event', () => {
    AssertionAnalyticsService.onRevertAssertion();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.TestResults, Actions.AssertionRevert, Labels.Button);
  });
});
