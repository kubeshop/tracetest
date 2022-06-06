import {SupportedDiagrams} from 'components/Diagram/Diagram';
import {Categories, Labels} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  ChangeTab = 'change-tab-click',
  AddAssertionButtonClick = 'add-assertion-button-click',
  TimelineSpanClick = 'timeline-span-click',
  AttributeCheckClick = 'attribute-check-click',
  SwitchDiagramView = 'switch-diagram-view-click',
  AttributeCopy = 'attribute-copy-click',
  RevertAllClick = 'revert-all-click',
  PublishClick = 'publish-click',
}

const TraceAnalyticsService = () => {
  const onChangeTab = (tabName: string) => {
    AnalyticsService.event(Categories.Trace, `${Actions.ChangeTab}-${tabName}`, Labels.Tab);
  };

  const onAddAssertionButtonClick = () => {
    AnalyticsService.event(Categories.Trace, Actions.AddAssertionButtonClick, Labels.Button);
  };

  const onTimelineSpanClick = () => {
    AnalyticsService.event(Categories.Trace, Actions.TimelineSpanClick, Labels.Button);
  };

  const onAttributeCopy = () => {
    AnalyticsService.event(Categories.SpanDetail, Actions.AttributeCopy, Labels.Button);
  };

  const onSwitchDiagramView = (diagramType: SupportedDiagrams) => {
    AnalyticsService.event(Categories.Trace, `${Actions.SwitchDiagramView}-${diagramType}`, Labels.Button);
  };

  const onAttributeCheckClick = () => {
    AnalyticsService.event(Categories.SpanDetail, Actions.AttributeCheckClick, Labels.Button);
  };

  const onRevertAllClick = () => {
    AnalyticsService.event(Categories.Trace, Actions.RevertAllClick, Labels.Button);
  };

  const onPublishClick = () => {
    AnalyticsService.event(Categories.Trace, Actions.PublishClick, Labels.Button);
  };

  return {
    onChangeTab,
    onAddAssertionButtonClick,
    onTimelineSpanClick,
    onAttributeCopy,
    onSwitchDiagramView,
    onAttributeCheckClick,
    onRevertAllClick,
    onPublishClick,
  };
};

export default TraceAnalyticsService();
