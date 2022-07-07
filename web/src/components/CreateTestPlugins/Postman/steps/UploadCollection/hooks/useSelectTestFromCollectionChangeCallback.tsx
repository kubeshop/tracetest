import {FormInstance} from 'antd';
import {useCallback} from 'react';
import {RecursivePartial} from 'utils/Common';
import {HTTP_METHOD} from '../../../../../../constants/Common.constants';
import {IRequestDetailsValues} from '../UploadCollection';
import {transformAuthSettings} from './transformAuthSettings';
import {State} from './useUploadCollectionCallback';

export function useSelectTestFromCollectionChangeCallback(state: State, form: FormInstance<IRequestDetailsValues>) {
  return useCallback(
    (identifier: string) => {
      const request = state.requests.find(({id}) => identifier === id);
      if (request) {
        const input: RecursivePartial<IRequestDetailsValues> = {
          url: request?.url.toString(),
          method: request?.method as HTTP_METHOD,
          headers: request?.headers?.all()?.map(({key, value}) => ({value, key})) || [],
          body: request?.body?.raw || '',
          auth: transformAuthSettings(request),
        };
        form.setFieldsValue(input);
      }
    },
    [state, form]
  );
}
