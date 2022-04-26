import useAnalytics, {Categories, Labels} from '../Analytics/useAnalytics';

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

const useTraceAnalytics = (): TTraceAnalytics => {
  const {event} = useAnalytics(Categories.Trace);

  const onChangeTab = (tabName: string) => {
    event(Actions.ChangeTab, tabName);
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

export default useTraceAnalytics;
