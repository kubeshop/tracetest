import {Categories, Labels} from 'constants/Analytics.constants';
import CreateTestAnalyticsService, {Actions} from '../CreateTestAnalytics.service';
import AnalyticsService from '../Analytics.service';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('CreateTestAnalyticsService', () => {
  it('should trigger the onDemoTestClick event', () => {
    CreateTestAnalyticsService.onCreateTestFormSubmit();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Home, Actions.CreateTestFormSubmit, Labels.Form);
  });

  it('should trigger the onDemoTestClick event', () => {
    CreateTestAnalyticsService.onDemoTestClick();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Home, Actions.DemoTestClick, Labels.Button);
  });
});
