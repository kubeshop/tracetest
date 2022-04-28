import AnalyticsService, {Categories} from './Analytics.service';

enum Actions {
  ClickSpan = 'click-span-node',
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
