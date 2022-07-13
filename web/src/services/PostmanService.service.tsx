import {FormInstance} from 'antd';
import {Collection, Item, ItemGroup, Request, RequestAuthDefinition, VariableDefinition} from 'postman-collection';
import {IUploadCollectionValues} from '../components/CreateTestPlugins/Postman/steps/UploadCollection/UploadCollection';
import {HTTP_METHOD} from '../constants/Common.constants';
import {TRequestAuth} from '../types/Test.types';
import {RecursivePartial} from '../utils/Common';

export interface RequestDefinitionExtended extends Request {
  id: string;
  name: string;
}

type AuthType = 'apiKey' | 'basic' | 'bearer';

const Postman = () => ({
  valuesFromRequest(
    requests: RequestDefinitionExtended[],
    variables: VariableDefinition[],
    identifier: string
  ): undefined | RecursivePartial<IUploadCollectionValues> {
    const request = requests.find(({id}) => identifier === id);
    if (request) {
      const url = request?.url.toString();
      return {
        url: this.substituteVariable(variables, url),
        method: request?.method as HTTP_METHOD,
        headers:
          request?.headers?.all()?.map(({key, value}) => ({
            value: this.substituteVariable(variables, value),
            key: this.substituteVariable(variables, key),
          })) || [],
        body: request?.body?.raw || '',
        auth: this.transformAuthSettings(request, variables),
      };
    }
    return undefined;
  },
  updateForm(
    requests: RequestDefinitionExtended[],
    variables: VariableDefinition[],
    identifier: string,
    form: FormInstance<IUploadCollectionValues>,
    setTransientUrl: (value: ((prevState: string) => string) | string) => void
  ): void {
    const input = this.valuesFromRequest(requests, variables, identifier);
    if (input) {
      form.setFieldsValue(input);
      if (input?.url) {
        setTransientUrl(input?.url);
      }
    }
  },
  flattenCollectionItemsIntoRequestDefinition(structure: Item | ItemGroup<Item>): RequestDefinitionExtended[] {
    const itemRequest: RequestDefinitionExtended[] =
      'request' in structure ? [{...structure.request, name: structure.name, id: structure.id} as Request] : [];
    const groupRequest: RequestDefinitionExtended[] =
      'items' in structure
        ? structure.items
            .all()
            .map(a => this.flattenCollectionItemsIntoRequestDefinition(a))
            .flat()
        : [];
    return [...itemRequest, ...groupRequest];
  },
  getRequestsFromCollection(collection: Collection): RequestDefinitionExtended[] {
    try {
      return collection.items.all().flatMap(test => this.flattenCollectionItemsIntoRequestDefinition(test));
    } catch (e) {
      return [];
    }
  },
  substituteVariable(variables: VariableDefinition[], value: string | undefined) {
    const regExpMatchArray = (value || '').match(/\{{([^}]+)}}/);
    return regExpMatchArray
      ? value
          ?.replaceAll(
            regExpMatchArray?.[0],
            variables.find(variable => variable.key === regExpMatchArray?.[1])?.value || value
          )
          .replaceAll(' ', '')
      : value;
  },

  translateType(type: NonNullable<RequestAuthDefinition['type']>): AuthType | undefined {
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
  },
  transformAuthSettings(request: RequestDefinitionExtended, variables: VariableDefinition[]): TRequestAuth {
    if (request?.auth) {
      if (['apikey', 'basic', 'bearer'].includes(request?.auth.type)) {
        const authParameters = request?.auth?.parameters().all();
        if (request.auth.type === 'apikey') {
          return {
            type: this.translateType(request.auth.type),
            apiKey: {
              in: authParameters?.find(({key}) => key === 'in')?.value || 'header',
              key: this.substituteVariable(variables, authParameters?.find(({key}) => key === 'key')?.value),
              value: this.substituteVariable(variables, authParameters?.find(({key}) => key === 'value')?.value),
            },
          };
        }
        if (request.auth.type === 'basic') {
          return {
            type: this.translateType(request.auth.type),
            basic: {
              username: this.substituteVariable(
                variables,
                authParameters?.find(({key}) => key === 'username')?.value || ''
              ),
              password: this.substituteVariable(variables, authParameters?.find(({key}) => key === 'password')?.value),
            },
          };
        }
        if (request.auth.type === 'bearer') {
          return {
            type: this.translateType(request.auth.type),
            bearer: {
              token: this.substituteVariable(variables, authParameters?.find(({key}) => key === 'token')?.value || ''),
            },
          };
        }
      }
    }
    return undefined;
  },
});

export default Postman();
