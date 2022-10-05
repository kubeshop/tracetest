import {Categories, Labels} from 'constants/Analytics.constants';
import TestAnalyticsService, {Actions} from '../TestAnalytics.service';
import AnalyticsService from '../Analytics.service';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('TestAnalyticsService', () => {
  it('should trigger the onRunTest event', () => {
    TestAnalyticsService.onRunTest();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Test, Actions.RunTest, Labels.Button);
  });

  it('should trigger the onTestRunClick event', () => {
    TestAnalyticsService.onTestRunClick();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Test, Actions.TestRunClick, Labels.Button);
  });

  it('should trigger the onTestCardCollapse event', () => {
    TestAnalyticsService.onTestCardCollapse();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Home, Actions.TestCardCollapse, Labels.Button);
  });

  it('should trigger the onDeleteTest event', () => {
    TestAnalyticsService.onDeleteTest();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Home, Actions.DeleteTest, Labels.Button);
  });

  it('should trigger the onDeleteTestRun event', () => {
    TestAnalyticsService.onDeleteTestRun();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Test, Actions.DeleteTestRun, Labels.Button);
  });

  it('should trigger the onDisplayTestInfo event', () => {
    TestAnalyticsService.onDisplayTestInfo();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.TestRun, Actions.DisplayTestInfo, Labels.Button);
  });
});
