import {Categories, Labels} from 'constants/Analytics.constants';
import HomeAnalyticsService, {Actions} from '../HomeAnalytics.service';
import AnalyticsService from '../Analytics.service';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('HomeAnalyticsService', () => {
  it('should trigger the onCreateTestClick event', () => {
    HomeAnalyticsService.onCreateTestClick();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Home, Actions.CreateTestClick, Labels.Button);
  });

  it('should trigger the onGuidedTourClick event', () => {
    HomeAnalyticsService.onGuidedTourClick();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Home, Actions.GuidedTourClick, Labels.Button);
  });

  it('should trigger the onTestClick event', () => {
    HomeAnalyticsService.onTestClick('testId');

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Home, Actions.TestClick, 'testId');
  });
});
