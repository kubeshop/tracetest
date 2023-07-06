import {SupportedPlugins} from 'constants/Common.constants';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import {TriggerTypes} from 'constants/Test.constants';
import {toRawTestOutputs} from 'models/TestOutput.model';
import {IPlugin} from 'types/Plugins.types';
import {TDraftTest} from 'types/Test.types';
import Validator from 'utils/Validator';
import Test, {TRawTestResource} from 'models/Test.model';
import TestDefinitionService from './TestDefinition.service';
import GrpcService from './Triggers/Grpc.service';
import HttpService from './Triggers/Http.service';
import PostmanService from './Triggers/Postman.service';
import CurlService from './Triggers/Curl.service';
import TraceIDService from './Triggers/TraceID.service';

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

const basicDetailsValidation = ({name}: TDraftTest): boolean => {
  return Validator.required(name);
};

const TriggerServiceMap = {
  [SupportedPlugins.GRPC]: GrpcService,
  [SupportedPlugins.REST]: HttpService,
  [SupportedPlugins.Messaging]: HttpService,
  [SupportedPlugins.OpenAPI]: HttpService,
  [SupportedPlugins.Postman]: PostmanService,
  [SupportedPlugins.CURL]: CurlService,
  [SupportedPlugins.TraceID]: TraceIDService,
} as const;

const TriggerServiceByTypeMap = {
  [TriggerTypes.grpc]: GrpcService,
  [TriggerTypes.http]: HttpService,
  [TriggerTypes.traceid]: TraceIDService,
} as const;

const TestService = () => ({
  async getRequest({type, name: pluginName}: IPlugin, draft: TDraftTest, original?: Test): Promise<TRawTestResource> {
    const {name, description} = draft;
    const triggerService = TriggerServiceMap[pluginName];
    const request = await triggerService.getRequest(draft);

    const trigger = {
      type,
      triggerType: type,
      [type]: request,
    };

    return {
      type: 'Test',
      spec: {
        name,
        description,
        trigger,
        ...(original
          ? {
              outputs: toRawTestOutputs(original.outputs ?? []),
              specs: original.definition.specs.map(def => TestDefinitionService.toRaw(def)),
            }
          : {}),
      },
    };
  },

  async validateDraft(pluginName: SupportedPlugins, draft: TDraftTest, isBasicDetails = false): Promise<boolean> {
    const triggerService = TriggerServiceMap[pluginName];
    const isTriggerValid = await triggerService.validateDraft(draft);

    return (isBasicDetails && basicDetailsValidation(draft)) || (isTriggerValid && authValidation(draft));
  },

  getInitialValues({trigger: {request, type}, name, description}: Test) {
    const triggerService = TriggerServiceByTypeMap[type];

    return {
      name,
      description,
      type,
      ...triggerService.getInitialValues!(request),
    };
  },

  getUpdatedRawTest(test: Test, partialTest: Partial<Test>): Promise<TRawTestResource> {
    const plugin = TriggerTypeToPlugin[test?.trigger?.type || TriggerTypes.http];
    const testTriggerData = this.getInitialValues(test);
    const updatedTest = {...test, ...partialTest};
    return this.getRequest(plugin, testTriggerData, updatedTest);
  },
});

export default TestService();
