import Demo from 'models/Demo.model';
import Polling from 'models/Polling.model';
import {ResourceType, SupportedDemosFormField, TDraftDemo, TDraftResource, TDraftSpec} from 'types/Settings.types';

const SettingService = () => ({
  getEnabledDemos(demos: Demo[]): Demo[] {
    return demos.filter(demo => demo.enabled);
  },
  getDefaultPollingProfile(pollingProfiles: Polling[]): Polling | undefined {
    return pollingProfiles.find(pollingProfile => pollingProfile.default);
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
