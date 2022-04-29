import AnalyticsService, { Categories } from "./Analytics.service";

enum Actions {
  SpanAssertionCLick = 'test-assertion-table-span-assertion-click',
}

type TTraceAssertionTableAnalytics = {
  onSpanAssertionClick(assertionSpanId: string): void;
};

const {event} = AnalyticsService(Categories.TestResults);

const TraceAssertionTableAnalyticsService = (): TTraceAssertionTableAnalytics => {
  const onSpanAssertionClick = (assertionSpanId: string) => {
    event(Actions.SpanAssertionCLick, assertionSpanId);
  };

  return {
    onSpanAssertionClick,
  };
};

export default TraceAssertionTableAnalyticsService();
