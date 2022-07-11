import {RequestAuthDefinition, VariableDefinition} from 'postman-collection';
import {TRequestAuth} from '../../../../../../types/Test.types';
import {RequestDefinitionExtended} from './getRequestsFromCollection';
import {substituteVariable} from './substituteVariable';

type AuthType = 'apiKey' | 'basic' | 'bearer';

function translateType(type: NonNullable<RequestAuthDefinition['type']>): AuthType | undefined {
  switch (type) {
    case 'apikey':
      return 'apiKey';
    case 'basic':
      return 'basic';
    case 'bearer':
      return 'bearer';
    default:
      return undefined;
  }
}

export function transformAuthSettings(
  request: RequestDefinitionExtended,
  variables: VariableDefinition[]
): TRequestAuth {
  if (request?.auth) {
    if (['apikey', 'basic', 'bearer'].includes(request?.auth.type)) {
      const authParameters = request?.auth?.parameters().all();
      if (request.auth.type === 'apikey') {
        return {
          type: translateType(request.auth.type),
          apiKey: {
            in: authParameters?.find(({key}) => key === 'in')?.value || 'header',
            key: substituteVariable(variables, authParameters?.find(({key}) => key === 'key')?.value),
            value: substituteVariable(variables, authParameters?.find(({key}) => key === 'value')?.value),
          },
        };
      }
      if (request.auth.type === 'basic') {
        return {
          type: translateType(request.auth.type),
          basic: {
            username: substituteVariable(variables, authParameters?.find(({key}) => key === 'username')?.value || ''),
            password: substituteVariable(variables, authParameters?.find(({key}) => key === 'password')?.value),
          },
        };
      }
      if (request.auth.type === 'bearer') {
        return {
          type: translateType(request.auth.type),
          bearer: {
            token: substituteVariable(variables, authParameters?.find(({key}) => key === 'token')?.value || ''),
          },
        };
      }
    }
  }
  return undefined;
}
