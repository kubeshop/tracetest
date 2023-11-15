import {Form} from 'antd';
import {RcFile} from 'antd/lib/upload';
import {VariableScope} from 'postman-collection';
import {useCallback} from 'react';
import PostmanService from 'services/Importers/Postman.service';
import {IPostmanValues, TDraftTestForm} from 'types/Test.types';

export function useUploadEnvFileCallback(form: TDraftTestForm<IPostmanValues>): (file?: RcFile) => Promise<void> {
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
            await PostmanService.updateForm(requests, variables, collectionTest, form);
          }
        }
      } catch (r) {
        // eslint-disable-next-line no-console
        console.error(r);
      }
    },
    [form, requests, collectionTest]
  );
}
