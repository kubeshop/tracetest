import Demo from 'models/Demo.model';
import Config from 'models/Config.model';
import {
  ResourceType,
  SupportedDemos,
  SupportedDemosFormField,
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
    const supportedDemos = Object.values(SupportedDemos);
    const supportedDemoFields = Object.values(SupportedDemosFormField);

    let draft = {};

    for (let i = 0; i < supportedDemos.length; i+=1) {
      const demoType = supportedDemos[i];
      const demoName = supportedDemoFields[i];

      const enabledDemo = demos.find(demo => demo.type === demoType);

      draft = {
        ...draft,
        [demoName]: enabledDemo || {
          type: demoName,
          enabled: false,
          name: demoName,
        },
      };
    }

    return draft as TDraftDemo;
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
