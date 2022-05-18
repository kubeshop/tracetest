import TraceAssertionTableAnalyticsService, {Actions} from '../TraceAssertionTableAnalytics.service';
import AnalyticsService from '../Analytics.service';
import {Categories} from '../../../constants/Analytics.constants';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('TraceAssertionTableAnalyticsService', () => {
  it('should trigger the onSpanAssertionClick event', () => {
    TraceAssertionTableAnalyticsService.onSpanAssertionClick('assertionSpanId');

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.TestResults,
      Actions.SpanAssertionCLick,
      'assertionSpanId'
    );
  });
});
