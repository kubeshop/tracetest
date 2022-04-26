import useAnalytics, {Categories} from '../Analytics/useAnalytics';

enum Actions {
  ClickSpan = 'click-span-node',
}

type TTraceDiagramAnalytics = {
  onClickSpan(spanId: string): void;
};

const useTraceDiagramAnalytics = (): TTraceDiagramAnalytics => {
  const {event} = useAnalytics(Categories.Trace);

  const onClickSpan = (spanId: string) => {
    event(Actions.ClickSpan, spanId);
  };

  return {
    onClickSpan,
  };
};

export default useTraceDiagramAnalytics;
