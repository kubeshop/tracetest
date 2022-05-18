import { Categories } from '../../constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  SpanAssertionCLick = 'test-assertion-table-span-assertion-click',
}

type TTraceAssertionTableAnalytics = {
  onSpanAssertionClick(assertionSpanId: string): void;
};

const TraceAssertionTableAnalyticsService = (): TTraceAssertionTableAnalytics => {
  const onSpanAssertionClick = (assertionSpanId: string) => {
    AnalyticsService.event(Categories.TestResults, Actions.SpanAssertionCLick, assertionSpanId);
  };

  return {
    onSpanAssertionClick,
  };
};

export default TraceAssertionTableAnalyticsService();
