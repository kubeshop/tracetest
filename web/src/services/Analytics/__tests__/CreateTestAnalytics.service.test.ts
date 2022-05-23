import CreateTestAnalyticsService, {Actions} from '../CreateTestAnalytics.service';
import AnalyticsService from '../Analytics.service';
import {Categories, Labels} from '../../../constants/Analytics.constants';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('CreateTestAnalyticsService', () => {
  it('should trigger the onCreateAssertionFormSubmit event', () => {
    CreateTestAnalyticsService.onCreateTestFormSubmit();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Home, Actions.CreateTestFormSubmit, Labels.Form);
  });
});
