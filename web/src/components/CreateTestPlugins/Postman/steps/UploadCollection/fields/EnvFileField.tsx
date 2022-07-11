import {Form, FormInstance} from 'antd';
import React, {Dispatch, SetStateAction} from 'react';
import RequestDetailsFileInput from '../../../../Rpc/steps/RequestDetails/RequestDetailsFileInput';
import {State} from '../hooks/useUploadCollectionCallback';
import {useUploadEnvFileCallback} from '../hooks/useUploadEnvFileCallback';
import {IRequestDetailsValues} from '../UploadCollection';

interface IProps {
  setState: Dispatch<SetStateAction<State>>;
  form: FormInstance<IRequestDetailsValues>;
  state: State;
  setTransientUrl: Dispatch<SetStateAction<string>>;
}

export const EnvFileField = ({state, form, setState, setTransientUrl}: IProps) => {
  const collectionFile = Form.useWatch('collectionFile');
  return (
    <Form.Item data-cy="envFile" name="envFile" label="Upload environment file (optional)">
      <RequestDetailsFileInput
        disabled={!collectionFile}
        accept=".json"
        onChange={useUploadEnvFileCallback(state, form, setState, setTransientUrl)}
      />
    </Form.Item>
  );
};
