import Demo from 'models/Demo.model';
import Config from 'models/Config.model';
import {
  ResourceType,
  SupportedDemosFormField,
  SupportedDemosFormFieldMap,
  TDraftConfig,
  TDraftDemo,
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

  // forms
  getDemoFormInitialValues(demos: Demo[]): TDraftDemo {
    return Object.values(SupportedDemosFormField).reduce((draft, demoName) => {
      const demoType = SupportedDemosFormFieldMap[demoName];
      const enabledDemo = demos.find(demo => demo.type === demoType);

      return {
        ...draft,
        [demoName]: enabledDemo || {
          type: demoType,
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
