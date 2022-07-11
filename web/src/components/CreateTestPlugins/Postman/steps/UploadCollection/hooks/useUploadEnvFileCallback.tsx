/* eslint-disable no-console */
import {Form, FormInstance} from 'antd';
import {UploadFile} from 'antd/es/upload/interface';
import {VariableScope} from 'postman-collection';
import {Dispatch, SetStateAction, useCallback} from 'react';
import {IRequestDetailsValues} from '../UploadCollection';
import {updateForm} from './useSelectTestCallback';
import {State} from './useUploadCollectionCallback';

export function useUploadEnvFileCallback(
  state: State,
  form: FormInstance<IRequestDetailsValues>,
  setState: Dispatch<SetStateAction<State>>,
  setTransientUrl: Dispatch<SetStateAction<string>>
): (file?: UploadFile) => void {
  const test = Form.useWatch('collectionTest');
  return useCallback(
    (file?: UploadFile) => {
      const readFile = new FileReader();
      readFile.onload = async (e: ProgressEvent<FileReader>) => {
        const contents = e?.target?.result;
        if (contents && typeof contents === 'string') {
          try {
            const variables = new VariableScope(JSON.parse(contents))?.values?.map(d => d) || [];
            const nextState = {requests: state.requests, variables};
            await setState(nextState);
            if (test) {
              await updateForm(nextState, test, form, setTransientUrl);
            }
          } catch (r) {
            console.info('erro');
            console.error(r);
          }
        }
      };

      readFile.readAsText(file as any);
    },
    [form, setState, setTransientUrl, state, test]
  );
}
