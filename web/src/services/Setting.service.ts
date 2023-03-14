import Demo from 'models/Demo.model';
import Polling from 'models/Polling.model';
import Config from 'models/Config.model';
import {
  ResourceType,
  SupportedDemosFormField,
  TDraftConfig,
  TDraftDemo,
  TDraftPollingProfiles,
  TDraftResource,
  TDraftSpec,
} from 'types/Settings.types';

const SettingService = () => ({
  getEnabledDemos(demos: Demo[]): Demo[] {
    return demos.filter(demo => demo.enabled);
  },
  getConfigInitialValues(config: Config): TDraftConfig {
    return (
      config || {
        name: 'current',
        analyticsEnabled: false,
      }
    );
  },

  getPollingProfileInitialValues(pollingProfiles: Polling[]): TDraftPollingProfiles {
    return (
      pollingProfiles.find(pollingProfile => pollingProfile.default) || {
        name: 'default',
        default: true,
        strategy: 'periodic',
      }
    );
  },

  // forms
  getDemoFormInitialValues(demos: Demo[]): TDraftDemo {
    return Object.values(SupportedDemosFormField).reduce((draft, demoName) => {
      const enabledDemo = demos.find(demo => demo.type === demoName);

      return {
        ...draft,
        [demoName]: enabledDemo || {
          type: demoName,
          enabled: false,
          name: demoName,
        },
      };
    }, {} as TDraftDemo);
  },

  getDemoFormValues(draft: TDraftDemo): TDraftResource[] {
    return Object.values(draft).map(demo => this.getDraftResource(ResourceType.DemoType, demo));
  },

  getDraftResource(resourceType: ResourceType, draft: TDraftSpec): TDraftResource {
    return {
      type: resourceType,
      spec: draft,
    };
  },
});

export default SettingService();
