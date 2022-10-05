import {Categories} from 'constants/Analytics.constants';
import TraceDiagramAnalyticsService, {Actions} from '../TraceDiagramAnalytics.service';
import AnalyticsService from '../Analytics.service';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('TraceDiagramAnalyticsService', () => {
  it('should trigger the onClickSpan event', () => {
    TraceDiagramAnalyticsService.onClickSpan('spanId');

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.TestRun, Actions.ClickSpan, 'spanId');
  });
});
