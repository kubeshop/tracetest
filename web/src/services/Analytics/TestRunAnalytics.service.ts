import {VisualizationType} from 'components/RunDetailTrace/RunDetailTrace';
import {Categories, Labels} from 'constants/Analytics.constants';
import {RunDetailModes} from 'constants/TestRun.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  ChangeMode = 'change-mode-click',
  ChangeTab = 'change-tab-click',
  TriggerResponseTabChange = 'trigger-response-tab-change',
  TriggerResponseHeaderCopy = 'trigger-response-header-copy',
  TriggerEditSubmit = 'trigger-edit-submit',
  LoadTestDefinition = 'load-test-definition',
  LoadJUnitReport = 'load-junit-report',
  AddAssertionButtonClick = 'add-assertion-button-click',
  TimelineSpanClick = 'timeline-span-click',
  AttributeCheckClick = 'attribute-check-click',
  SwitchDiagramView = 'switch-diagram-view-click',
  AttributeCopy = 'attribute-copy-click',
  RevertAllClick = 'revert-all-click',
  PublishClick = 'publish-click',
  AttributeDrawerOpen = 'attribute-drawer-toggle-click',
}

const TestRunAnalyticsService = () => ({
  onChangeMode(mode: RunDetailModes) {
    AnalyticsService.event(Categories.TestRun, Actions.ChangeMode, mode);
  },
  onChangeTab(tabName: string) {
    AnalyticsService.event(Categories.TestRun, `${Actions.ChangeTab}-${tabName}`, Labels.Tab);
  },
  onTriggerResponseTabChange(tabName: string) {
    AnalyticsService.event(Categories.TestRun, Actions.TriggerResponseTabChange, tabName);
  },
  onTriggerResponseHeaderCopy() {
    AnalyticsService.event(Categories.TestRun, Actions.TriggerResponseHeaderCopy, Labels.Button);
  },
  onTriggerEditSubmit() {
    AnalyticsService.event(Categories.TestRun, Actions.TriggerEditSubmit, Labels.Form);
  },
  onLoadJUnitReport() {
    AnalyticsService.event(Categories.TestRun, Actions.LoadJUnitReport, Labels.Button);
  },
  onLoadTestDefinition() {
    AnalyticsService.event(Categories.TestRun, Actions.LoadTestDefinition, Labels.Button);
  },
  onAttributeDrawerOpen() {
    AnalyticsService.event(Categories.TestRun, Actions.AttributeDrawerOpen, Labels.Button);
  },
  onAddAssertionButtonClick() {
    AnalyticsService.event(Categories.TestRun, Actions.AddAssertionButtonClick, Labels.Button);
  },
  onTimelineSpanClick(spanId: string) {
    AnalyticsService.event(Categories.TestRun, Actions.TimelineSpanClick, spanId);
  },
  onAttributeCopy() {
    AnalyticsService.event(Categories.SpanDetail, Actions.AttributeCopy, Labels.Button);
  },
  onSwitchDiagramView(diagramType: VisualizationType) {
    AnalyticsService.event(Categories.TestRun, `${Actions.SwitchDiagramView}-${diagramType}`, Labels.Button);
  },
  onAttributeCheckClick() {
    AnalyticsService.event(Categories.SpanDetail, Actions.AttributeCheckClick, Labels.Button);
  },
  onRevertAllClick() {
    AnalyticsService.event(Categories.TestRun, Actions.RevertAllClick, Labels.Button);
  },
  onPublishClick() {
    AnalyticsService.event(Categories.TestRun, Actions.PublishClick, Labels.Button);
  },
});

export default TestRunAnalyticsService();
