import AnalyticsService, {Categories} from './Analytics.service';

enum Actions {
  ClickSpan = 'span-node-click',
}

type TTraceDiagramAnalytics = {
  onClickSpan(spanId: string): void;
};

const {event} = AnalyticsService(Categories.Trace);

const TraceDiagramAnalyticsService = (): TTraceDiagramAnalytics => {
  const onClickSpan = (spanId: string) => {
    event(Actions.ClickSpan, spanId);
  };

  return {
    onClickSpan,
  };
};

export default TraceDiagramAnalyticsService();
