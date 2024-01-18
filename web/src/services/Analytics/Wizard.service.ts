import {Categories} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

enum Actions {
  StepSelect = 'wizard-step-select',
  AgentTypeSelect = 'wizard-agent-type-select',
  StepComplete = 'wizard-step-complete',
  TracingBackendTypeSelect = 'wizard-tracing-backend-type-select',
  TriggerTypeSelect = 'wizard-trigger-type-select',
}

type TWizardAnalytics = {
  onStepSelect(step: string): void;
  onAgentTypeSelect(agentType: string): void;
  onStepComplete(step: string): void;
  onTracingBackendTypeSelect(tracingBackendType: string): void;
  onTriggerTypeSelect(triggerType: string): void;
};

const WizardAnalytics = (): TWizardAnalytics => ({
  onStepSelect(step: string) {
    AnalyticsService.event(Categories.Wizard, Actions.StepSelect, step);
  },
  onAgentTypeSelect(agentType: string) {
    AnalyticsService.event(Categories.Wizard, Actions.AgentTypeSelect, agentType);
  },
  onStepComplete(step: string) {
    AnalyticsService.event(Categories.Wizard, Actions.StepComplete, step);
  },
  onTracingBackendTypeSelect(tracingBackendType: string) {
    AnalyticsService.event(Categories.Wizard, Actions.TracingBackendTypeSelect, tracingBackendType);
  },
  onTriggerTypeSelect(triggerType: string) {
    AnalyticsService.event(Categories.Wizard, Actions.TriggerTypeSelect, triggerType);
  },
});

export default WizardAnalytics();
