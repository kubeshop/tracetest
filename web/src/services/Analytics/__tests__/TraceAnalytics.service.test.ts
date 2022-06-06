import {Categories, Labels} from 'constants/Analytics.constants';
import TraceAnalyticsService, {Actions} from '../TraceAnalytics.service';
import AnalyticsService from '../Analytics.service';
import {SupportedDiagrams} from '../../../components/Diagram/Diagram';

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
    TraceAnalyticsService.onTimelineSpanClick();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Trace, Actions.TimelineSpanClick, Labels.Button);
  });

  it('should trigger the onAttributeCopy event', () => {
    TraceAnalyticsService.onAttributeCopy();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Trace, Actions.TimelineSpanClick, Labels.Button);
  });

  it('should trigger the onSwitchDiagramView event', () => {
    TraceAnalyticsService.onSwitchDiagramView(SupportedDiagrams.DAG);

    expect(AnalyticsService.event).toHaveBeenCalledWith(
      Categories.Trace,
      `${Actions.SwitchDiagramView}-${SupportedDiagrams.DAG}`,
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

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Trace, Actions.RevertAllClick, Labels.Button);
  });

  it('should trigger the onPublishClick event', () => {
    TraceAnalyticsService.onPublishClick();

    expect(AnalyticsService.event).toHaveBeenCalledWith(Categories.Trace, Actions.PublishClick, Labels.Button);
  });
});
