import {Collection, Item, ItemGroup, Request, RequestAuthDefinition, VariableDefinition} from 'postman-collection';
import {HTTP_METHOD} from 'constants/Common.constants';
import {IImportService, IPostmanValues, TDraftTestForm, TRequestAuth} from 'types/Test.types';
import Validator from 'utils/Validator';
import HttpService from '../Triggers/Http.service';
import {Plugins} from '../../constants/Plugins.constants';

export interface RequestDefinitionExtended extends Request {
  id: string;
  name: string;
}

type AuthType = 'apiKey' | 'basic' | 'bearer';

interface IPostmanTriggerService extends IImportService {
  valuesFromRequest(
    requests: RequestDefinitionExtended[],
    variables: VariableDefinition[],
    identifier: string
  ): undefined | Partial<IPostmanValues>;
  updateForm(
    requests: RequestDefinitionExtended[],
    variables: VariableDefinition[],
    identifier: string,
    form: TDraftTestForm<IPostmanValues>
  ): void;
  flattenCollectionItemsIntoRequestDefinition(structure: Item | ItemGroup<Item>): RequestDefinitionExtended[];
  getRequestsFromCollection(collection: Collection): RequestDefinitionExtended[];
  substituteVariable(variables: VariableDefinition[], value: string | undefined): any;
  translateType(type: NonNullable<RequestAuthDefinition['type']>): AuthType | undefined;
  transformAuthSettings(request: RequestDefinitionExtended, variables: VariableDefinition[]): TRequestAuth;
}

const Postman = (): IPostmanTriggerService => ({
  async getRequest(values) {
    const {collectionTest, variables, requests} = values as IPostmanValues;
    const draft = (await this.valuesFromRequest(requests, variables, collectionTest || '')) || {};

    return {
      draft,
      plugin: Plugins.REST,
    };
  },
  async validateDraft(values) {
    const {collectionTest, variables, requests} = values as IPostmanValues;

    if (!Validator.required(collectionTest)) return false;

    const draft = await this.valuesFromRequest(requests, variables, collectionTest || '');
    return !!draft && HttpService.validateDraft(draft);
  },
  valuesFromRequest(requests, variables, identifier) {
    const request = requests.find(({id}) => identifier === id);
    if (request) {
      const url = request?.url.toString() || '';
      const headers =
        request?.headers?.all()?.map(({key, value}) => ({
          value: this.substituteVariable(variables, value) || '',
          key: this.substituteVariable(variables, key) || '',
        })) || [];

      return {
        name: request.name,
        url: this.substituteVariable(variables, url),
        method: request?.method as HTTP_METHOD,
        headers,
        body: request?.body?.raw || '',
        auth: this.transformAuthSettings(request, variables),
      };
    }
    return undefined;
  },
  updateForm(requests, variables, identifier, form) {
    const input = this.valuesFromRequest(requests, variables, identifier);
    if (input) {
      form.setFieldsValue(input);
    }
  },
  flattenCollectionItemsIntoRequestDefinition(structure) {
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
  getRequestsFromCollection(collection) {
    try {
      return collection.items.all().flatMap(test => this.flattenCollectionItemsIntoRequestDefinition(test));
    } catch (e) {
      return [];
    }
  },
  substituteVariable(variables, value) {
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

  translateType(type) {
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
  transformAuthSettings(request, variables) {
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
