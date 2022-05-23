import TraceAnalyticsService, {Actions} from '../TraceAnalytics.service';
import AnalyticsService from '../Analytics.service';
import {Categories, Labels} from '../../../constants/Analytics.constants';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('TraceAnalyticsService', () => {
  it('should trigger the onChangeTab event', () => {
    TraceAnalyticsService.onChangeTab('request');

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Trace, `${Actions.ChangeTab}-request`, Labels.Tab);
  });

  it('should trigger the onAddAssertionButtonClick event', () => {
    TraceAnalyticsService.onAddAssertionButtonClick();

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.Trace,
      Actions.AddAssertionButtonClick,
      Labels.Button
    );
  });

  it('should trigger the onTimelineSpanClick event', () => {
    TraceAnalyticsService.onTimelineSpanClick('spanId');

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Trace, Actions.TimelineSpanClick, 'spanId');
  });
});
