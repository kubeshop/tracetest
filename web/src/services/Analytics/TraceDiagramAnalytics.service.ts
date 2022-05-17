import {Categories} from '../../constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  ClickSpan = 'span-node-click',
}

type TTraceDiagramAnalytics = {
  onClickSpan(spanId: string): void;
};

const TraceDiagramAnalyticsService = (): TTraceDiagramAnalytics => {
  const onClickSpan = (spanId: string) => {
    AnalyticsService.event(Categories.Trace, Actions.ClickSpan, spanId);
  };

  return {
    onClickSpan,
  };
};

export default TraceDiagramAnalyticsService();
