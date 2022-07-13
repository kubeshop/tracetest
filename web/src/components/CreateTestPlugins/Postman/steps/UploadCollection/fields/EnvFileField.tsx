import {Form, FormInstance} from 'antd';
import React, {Dispatch, SetStateAction} from 'react';
import RequestDetailsFileInput from '../../../../Rpc/steps/RequestDetails/RequestDetailsFileInput';
import {useUploadEnvFileCallback} from '../hooks/useUploadEnvFileCallback';
import {IRequestDetailsValues} from '../UploadCollection';

interface IProps {
  form: FormInstance<IRequestDetailsValues>;
  setTransientUrl: Dispatch<SetStateAction<string>>;
}

export const EnvFileField = ({form, setTransientUrl}: IProps) => {
  const collectionFile = Form.useWatch('collectionFile');
  return (
    <Form.Item data-cy="envFile" name="envFile" label="Upload environment file (optional)">
      <RequestDetailsFileInput
        disabled={!collectionFile}
        accept=".json"
        onChange={useUploadEnvFileCallback(form, setTransientUrl)}
      />
    </Form.Item>
  );
};
