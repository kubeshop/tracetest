import AnalyticsService, {Categories, Labels} from './Analytics.service';

enum Actions {
  ChangeTab = 'change-tab',
  AddAssertionButtonClick = 'add-assertion-button-click',
  TimelineSpanClick = 'timeline-span-click',
}

type TTraceAnalytics = {
  onChangeTab(tabName: string): void;
  onAddAssertionButtonClick(): void;
  onTimelineSpanClick(spanId: string): void;
};

const {event} = AnalyticsService(Categories.Trace);

const TraceAnalyticsService = (): TTraceAnalytics => {
  const onChangeTab = (tabName: string) => {
    event(`${Actions.ChangeTab}-${tabName}`, Labels.Table);
  };

  const onAddAssertionButtonClick = () => {
    event(Actions.AddAssertionButtonClick, Labels.Button);
  };

  const onTimelineSpanClick = (spanId: string) => {
    event(Actions.TimelineSpanClick, spanId);
  };

  return {
    onChangeTab,
    onAddAssertionButtonClick,
    onTimelineSpanClick,
  };
};

export default TraceAnalyticsService();
