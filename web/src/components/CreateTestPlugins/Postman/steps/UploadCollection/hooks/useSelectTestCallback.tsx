import {FormInstance} from 'antd';
import {VariableDefinition} from 'postman-collection';
import {useCallback} from 'react';
import {RecursivePartial} from 'utils/Common';
import {HTTP_METHOD} from '../../../../../../constants/Common.constants';
import {IRequestDetailsValues} from '../UploadCollection';
import {RequestDefinitionExtended} from './getRequestsFromCollection';
import {substituteVariable} from './substituteVariable';
import {transformAuthSettings} from './transformAuthSettings';

export function valuesFromRequest(
  requests: RequestDefinitionExtended[],
  variables: VariableDefinition[],
  identifier: string
): undefined | RecursivePartial<IRequestDetailsValues> {
  const request = requests.find(({id}) => identifier === id);
  if (request) {
    const url = request?.url.toString();
    return {
      url: substituteVariable(variables, url),
      method: request?.method as HTTP_METHOD,
      headers:
        request?.headers?.all()?.map(({key, value}) => ({
          value: substituteVariable(variables, value),
          key: substituteVariable(variables, key),
        })) || [],
      body: request?.body?.raw || '',
      auth: transformAuthSettings(request, variables),
    };
  }
  return undefined;
}

export function updateForm(
  requests: RequestDefinitionExtended[],
  variables: VariableDefinition[],
  identifier: string,
  form: FormInstance<IRequestDetailsValues>,
  setTransientUrl: (value: ((prevState: string) => string) | string) => void
): void {
  const input = valuesFromRequest(requests, variables, identifier);
  if (input) {
    form.setFieldsValue(input);
    if (input?.url) {
      setTransientUrl(input?.url);
    }
  }
}

export function useSelectTestCallback(
  form: FormInstance<IRequestDetailsValues>,
  setTransientUrl: React.Dispatch<React.SetStateAction<string>>,
  requests: RequestDefinitionExtended[],
  variables: VariableDefinition[]
) {
  return useCallback(
    (identifier: string) => {
      updateForm(requests, variables, identifier, form, setTransientUrl);
    },
    [form, setTransientUrl, requests, variables]
  );
}
