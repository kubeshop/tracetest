import {FormInstance} from 'antd';
import {Dispatch, SetStateAction, useCallback} from 'react';
import {RecursivePartial} from 'utils/Common';
import {HTTP_METHOD} from '../../../../../../constants/Common.constants';
import {IRequestDetailsValues} from '../UploadCollection';
import {substituteVariable} from './substituteVariable';
import {transformAuthSettings} from './transformAuthSettings';
import {State} from './useUploadCollectionCallback';

export function valuesFromRequest(
  {requests, variables}: State,
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
  state: State,
  identifier: string,
  form: FormInstance<IRequestDetailsValues>,
  setTransientUrl: (value: ((prevState: string) => string) | string) => void
): void {
  const input = valuesFromRequest(state, identifier);
  if (input) {
    form.setFieldsValue(input);
    if (input?.url) {
      setTransientUrl(input?.url);
    }
  }
}

export function useSelectTestCallback(
  state: State,
  form: FormInstance<IRequestDetailsValues>,
  setTransientUrl: Dispatch<SetStateAction<string>>
) {
  return useCallback(
    (identifier: string) => {
      updateForm(state, identifier, form, setTransientUrl);
    },
    [state, form, setTransientUrl]
  );
}
