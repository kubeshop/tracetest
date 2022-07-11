import {
  TRawTest,
  TTest,
  TDraftTest,
  // IBasicValues,
} from 'types/Test.types';
import {TriggerTypes} from 'constants/Test.constants';
import TestDefinitionService from './TestDefinition.service';
import Validator from '../utils/Validator';
import RpcService from './Triggers/Rpc.service';
import HttpService from './Triggers/Http.service';

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
  [TriggerTypes.grpc]: RpcService,
  [TriggerTypes.http]: HttpService,
} as const;

const TestService = () => ({
  async getRequest(type: TriggerTypes, draft: TDraftTest, original?: TTest): Promise<TRawTest> {
    const {name, description} = draft;
    const triggerService = TriggerServiceMap[type];
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

  async validateDraft(type: TriggerTypes, draft: TDraftTest, isBasicDetails = false): Promise<boolean> {
    const triggerService = TriggerServiceMap[type];
    const isTriggerValid = await triggerService.validateDraft(draft);

    return (isBasicDetails && basicDetailsValidation(draft)) || (isTriggerValid && authValidation(draft));
  },
});

export default TestService();
