import {Categories} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  ClickSpan = 'span-node-click',
}

type TTraceDiagramAnalytics = {
  onClickSpan(spanId: string): void;
};

const TraceDiagramAnalyticsService = (): TTraceDiagramAnalytics => ({
  onClickSpan(spanId) {
    AnalyticsService.event(Categories.TestRun, Actions.ClickSpan, spanId);
  },
});

export default TraceDiagramAnalyticsService();
