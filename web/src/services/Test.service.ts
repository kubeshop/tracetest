import {TRawTest, TTest, TDraftTest} from 'types/Test.types';
import { SupportedPlugins } from 'constants/Plugins.constants';
import { IPlugin } from 'types/Plugins.types';
import TestDefinitionService from './TestDefinition.service';
import Validator from '../utils/Validator';
import RpcService from './Triggers/Rpc.service';
import HttpService from './Triggers/Http.service';
import PostmanService from './Triggers/Postman.service';

const authValidation = ({auth}: TDraftTest): boolean => {
  switch (auth?.type) {
    case 'apiKey':
      return Validator.required(auth?.apiKey?.key) && Validator.required(auth?.apiKey?.value);
    case 'basic':
      return Validator.required(auth?.basic?.username) && Validator.required(auth?.basic?.password);
    case 'bearer':
      return Validator.required(auth?.bearer?.token);
    default:
      return true;
  }
};

const basicDetailsValidation = ({name, description}: TDraftTest): boolean => {
  return Validator.required(name) && Validator.required(description);
};

const TriggerServiceMap = {
  [SupportedPlugins.RPC]: RpcService,
  [SupportedPlugins.REST]: HttpService,
  [SupportedPlugins.Messaging]: HttpService,
  [SupportedPlugins.OpenAPI]: HttpService,
  [SupportedPlugins.Postman]: PostmanService,
} as const;

const TestService = () => ({
  async getRequest({type, name: pluginName}: IPlugin, draft: TDraftTest, original?: TTest): Promise<TRawTest> {
    const {name, description} = draft;
    const triggerService = TriggerServiceMap[pluginName];
    const request = await triggerService.getRequest(draft);

    return {
      name,
      description,
      serviceUnderTest: {
        triggerType: type,
        triggerSettings: {
          [type]: request,
        },
      },
      ...(original
        ? {
            definition: {
              definitions: original.definition.definitionList.map(def => TestDefinitionService.toRaw(def)),
            },
          }
        : {}),
    };
  },

  async validateDraft(pluginName: SupportedPlugins, draft: TDraftTest, isBasicDetails = false): Promise<boolean> {
    const triggerService = TriggerServiceMap[pluginName];
    const isTriggerValid = await triggerService.validateDraft(draft);

    return (isBasicDetails && basicDetailsValidation(draft)) || (isTriggerValid && authValidation(draft));
  },
});

export default TestService();
