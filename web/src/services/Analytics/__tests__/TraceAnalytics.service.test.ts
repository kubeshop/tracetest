import {VisualizationType} from 'components/RunDetailTrace/RunDetailTrace';
import {Categories, Labels} from 'constants/Analytics.constants';
import TraceAnalyticsService, {Actions} from '../TestRunAnalytics.service';
import AnalyticsService from '../Analytics.service';

jest.mock('../Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe('TraceAnalyticsService', () => {
  it('should trigger the onChangeTab event', () => {
    TraceAnalyticsService.onChangeTab('request');

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.TestRun, `${Actions.ChangeTab}-request`, Labels.Tab);
  });

  it('should trigger the onAddAssertionButtonClick event', () => {
    TraceAnalyticsService.onAddAssertionButtonClick();

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.TestRun,
      Actions.AddAssertionButtonClick,
      Labels.Button
    );
  });

  it('should trigger the onTimelineSpanClick event', () => {
    const spanId = '1234';
    TraceAnalyticsService.onTimelineSpanClick(spanId);

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.TestRun, Actions.TimelineSpanClick, spanId);
  });

  it('should trigger the onAttributeCopy event', () => {
    TraceAnalyticsService.onAttributeCopy();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.SpanDetail, Actions.AttributeCopy, Labels.Button);
  });

  it('should trigger the onSwitchDiagramView event', () => {
    TraceAnalyticsService.onSwitchDiagramView(VisualizationType.Dag);

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.TestRun,
      `${Actions.SwitchDiagramView}-${VisualizationType.Dag}`,
      Labels.Button
    );
  });

  it('should trigger the onAttributeCheckClick event', () => {
    TraceAnalyticsService.onAttributeCheckClick();

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.SpanDetail,
      Actions.AttributeCheckClick,
      Labels.Button
    );
  });

  it('should trigger the onRevertAllClick event', () => {
    TraceAnalyticsService.onRevertAllClick();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.TestRun, Actions.RevertAllClick, Labels.Button);
  });

  it('should trigger the onPublishClick event', () => {
    TraceAnalyticsService.onPublishClick();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.TestRun, Actions.PublishClick, Labels.Button);
  });
});
