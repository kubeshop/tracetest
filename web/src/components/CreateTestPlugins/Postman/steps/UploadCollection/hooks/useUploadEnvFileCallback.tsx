/* eslint-disable no-console */
import {Form, FormInstance} from 'antd';
import {UploadFile} from 'antd/es/upload/interface';
import {VariableScope} from 'postman-collection';
import {Dispatch, SetStateAction, useCallback} from 'react';
import {IRequestDetailsValues} from '../UploadCollection';
import {updateForm} from './useSelectTestCallback';

export function useUploadEnvFileCallback(
  form: FormInstance<IRequestDetailsValues>,
  setTransientUrl: Dispatch<SetStateAction<string>>
): (file?: UploadFile) => void {
  const collectionTest = Form.useWatch('collectionTest');
  const requests = Form.useWatch('requests');
  return useCallback(
    (file?: UploadFile) => {
      const readFile = new FileReader();
      readFile.onload = async (e: ProgressEvent<FileReader>) => {
        const contents = e?.target?.result;
        if (contents && typeof contents === 'string') {
          try {
            const variables = new VariableScope(JSON.parse(contents))?.values?.map(d => d) || [];
            form.setFieldsValue({requests, variables});
            if (collectionTest) {
              await updateForm(requests, variables, collectionTest, form, setTransientUrl);
            }
          } catch (r) {
            console.info('erro');
            console.error(r);
          }
        }
      };

      readFile.readAsText(file as any);
    },
    [form, requests, collectionTest, setTransientUrl]
  );
}
