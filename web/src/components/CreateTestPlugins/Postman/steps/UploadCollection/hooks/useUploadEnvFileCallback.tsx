import {Form, FormInstance} from 'antd';
import {RcFile} from 'antd/lib/upload';
import {VariableScope} from 'postman-collection';
import {Dispatch, SetStateAction, useCallback} from 'react';
import PostmanService from 'services/PostmanService.service';
import {IRequestDetailsValues} from '../UploadCollection';

export function useUploadEnvFileCallback(
  form: FormInstance<IRequestDetailsValues>,
  setTransientUrl: Dispatch<SetStateAction<string>>
): (file?: RcFile) => Promise<void> {
  const collectionTest = Form.useWatch('collectionTest');
  const requests = Form.useWatch('requests');
  return useCallback(
    async (file?: RcFile) => {
      try {
        const contents = await file?.text();
        if (contents) {
          const variables = new VariableScope(JSON.parse(contents))?.values?.map(d => d) || [];
          form.setFieldsValue({requests, variables});
          if (collectionTest) {
            await PostmanService.updateForm(requests, variables, collectionTest, form, setTransientUrl);
          }
        }
      } catch (r) {
        // eslint-disable-next-line no-console
        console.error(r);
      }
    },
    [form, requests, collectionTest, setTransientUrl]
  );
}
