import useAnalytics, {Categories} from '../Analytics/useAnalytics';

enum Actions {
  SpanAssertionCLick = 'test-assertion-table-span-assertion-click',
}

type TTraceAssertionTableAnalytics = {
  onSpanAssertionClick(assertionSpanId: string): void;
};

const useTraceAssertionTableAnalytics = (): TTraceAssertionTableAnalytics => {
  const {event} = useAnalytics(Categories.TestResults);

  const onSpanAssertionClick = (assertionSpanId: string) => {
    event(Actions.SpanAssertionCLick, assertionSpanId);
  };

  return {
    onSpanAssertionClick,
  };
};

export default useTraceAssertionTableAnalytics;
