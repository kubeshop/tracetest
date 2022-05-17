import AssertionTableAnalyticsService, {Actions} from '../AssertionTableAnalytics.service';
import AnalyticsService from '../Analytics.service';
import {Categories} from '../../../constants/Analytics.constants';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('AssertionTableAnalyticsService', () => {
  it('should trigger the onEditAssertionButtonClick event', () => {
    AssertionTableAnalyticsService.onEditAssertionButtonClick('assertionId');

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.SpanDetail,
      Actions.EditAssertionButtonClick,
      'assertionId'
    );
  });
});
