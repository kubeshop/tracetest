import TraceDiagramAnalyticsService, {Actions} from '../TraceDiagramAnalytics.service';
import AnalyticsService from '../Analytics.service';
import {Categories} from '../../../constants/Analytics.constants';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('TraceDiagramAnalyticsService', () => {
  it('should trigger the onClickSpan event', () => {
    TraceDiagramAnalyticsService.onClickSpan('spanId');

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.Trace,
      Actions.ClickSpan,
      'spanId'
    );
  });
});
