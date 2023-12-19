import {Categories, Labels} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

enum Actions {
  PluginClick = 'trace-analyzer-plugin-click',
  SpanNameClick = 'trace-analyzer-span-name-click',
  DocsClick = 'trace-analyzer-docs-click',
  SpanErrorsClick = 'trace-analyzer-span-errors-click',
}

type TTraceAnalyzerAnalytics = {
  onPluginClick(): void;
  onSpanNameClick(): void;
  onDocsClick(): void;
  onSpanErrorsClick(): void;
};

const TraceAnalyzerAnalytics = (): TTraceAnalyzerAnalytics => ({
  onPluginClick: () => {
    AnalyticsService.event(Categories.TraceAnalyzer, Actions.PluginClick, Labels.Button);
  },
  onSpanNameClick: () => {
    AnalyticsService.event(Categories.TraceAnalyzer, Actions.SpanNameClick, Labels.Button);
  },
  onDocsClick: () => {
    AnalyticsService.event(Categories.TraceAnalyzer, Actions.DocsClick, Labels.Link);
  },
  onSpanErrorsClick: () => {
    AnalyticsService.event(Categories.TraceAnalyzer, Actions.SpanErrorsClick, Labels.Button);
  },
});

export default TraceAnalyzerAnalytics();
