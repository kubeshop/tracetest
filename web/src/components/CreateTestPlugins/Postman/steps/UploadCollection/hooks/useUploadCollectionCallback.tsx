import {FormInstance} from 'antd';
import {RcFile} from 'antd/lib/upload';
import {Collection} from 'postman-collection';
import {useCallback} from 'react';
import {IRequestDetailsValues} from '../UploadCollection';
import {getRequestsFromCollection} from './getRequestsFromCollection';

export function useUploadCollectionCallback(form: FormInstance<IRequestDetailsValues>): (file?: RcFile) => void {
  return useCallback(
    async (file?: RcFile) => {
      const contents = await file?.text();
      if (contents && typeof contents === 'string') {
        try {
          const collection = new Collection(JSON.parse(contents));
          form.setFieldsValue({
            variables: collection.variables.all(),
            requests: getRequestsFromCollection(collection),
          });
        } catch (r) {
          // eslint-disable-next-line no-console
          console.error('error');
        }
      }
    },
    [form]
  );
}
