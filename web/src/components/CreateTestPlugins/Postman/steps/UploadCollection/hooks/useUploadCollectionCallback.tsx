import {FormInstance} from 'antd';
import {RcFile} from 'antd/lib/upload';
import {Collection} from 'postman-collection';
import {useCallback} from 'react';
import PostmanService from '../../../../../../services/PostmanService.service';
import {IUploadCollectionValues} from '../UploadCollection';

export function useUploadCollectionCallback(form: FormInstance<IUploadCollectionValues>): (file?: RcFile) => void {
  return useCallback(
    async (file?: RcFile) => {
      try {
        const contents = await file?.text();
        if (contents) {
          const collection = new Collection(JSON.parse(contents));
          form.setFieldsValue({
            variables: collection.variables.all(),
            requests: PostmanService.getRequestsFromCollection(collection),
          });
        }
      } catch (r) {
        // eslint-disable-next-line no-console
        console.error('error');
      }
    },
    [form]
  );
}
