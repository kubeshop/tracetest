import {FormInstance} from 'antd';
import {Dispatch, SetStateAction, useCallback} from 'react';
import {RecursivePartial} from 'utils/Common';
import {HTTP_METHOD} from '../../../../../../constants/Common.constants';
import {IRequestDetailsValues} from '../UploadCollection';
import {transformAuthSettings} from './transformAuthSettings';
import {State} from './useUploadCollectionCallback';

export function useSelectTestCallback(
  state: State,
  form: FormInstance<IRequestDetailsValues>,
  setTransientUrl: Dispatch<SetStateAction<string>>
) {
  return useCallback(
    (identifier: string) => {
      const request = state.requests.find(({id}) => identifier === id);
      if (request) {
        let url = request?.url.toString();
        const input: RecursivePartial<IRequestDetailsValues> = {
          url,
          method: request?.method as HTTP_METHOD,
          headers: request?.headers?.all()?.map(({key, value}) => ({value, key})) || [],
          body: request?.body?.raw || '',
          auth: transformAuthSettings(request),
        };
        form.setFieldsValue(input);
        setTransientUrl(url);
      }
    },
    [state, form]
  );
}
