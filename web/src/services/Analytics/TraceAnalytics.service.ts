import {Categories, Labels} from '../../constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  ChangeTab = 'change-tab',
  AddAssertionButtonClick = 'add-assertion-button-click',
  TimelineSpanClick = 'timeline-span-click',
}

type TTraceAnalytics = {
  onChangeTab(tabName: string): void;
  onAddAssertionButtonClick(): void;
  onTimelineSpanClick(spanId: string): void;
};

const TraceAnalyticsService = (): TTraceAnalytics => {
  const onChangeTab = (tabName: string) => {
    AnalyticsService.event(Categories.Trace, `${Actions.ChangeTab}-${tabName}`, Labels.Tab);
  };

  const onAddAssertionButtonClick = () => {
    AnalyticsService.event(Categories.Trace, Actions.AddAssertionButtonClick, Labels.Button);
  };

  const onTimelineSpanClick = (spanId: string) => {
    AnalyticsService.event(Categories.Trace, Actions.TimelineSpanClick, spanId);
  };

  return {
    onChangeTab,
    onAddAssertionButtonClick,
    onTimelineSpanClick,
  };
};

export default TraceAnalyticsService();
