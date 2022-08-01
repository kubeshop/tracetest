import {Form} from 'antd';
import {IPostmanValues, TDraftTestForm} from 'types/Test.types';
import RequestDetailsFileInput from '../../../../Grpc/steps/RequestDetails/RequestDetailsFileInput';
import {useUploadEnvFileCallback} from '../hooks/useUploadEnvFileCallback';

interface IProps {
  form: TDraftTestForm<IPostmanValues>;
}

export const EnvFileField = ({form}: IProps) => {
  const collectionFile = Form.useWatch('collectionFile');
  return (
    <Form.Item data-cy="envFile" name="envFile" label="Upload environment file (optional)">
      <RequestDetailsFileInput
        disabled={!collectionFile}
        accept=".json"
        onChange={useUploadEnvFileCallback(form)}
      />
    </Form.Item>
  );
};
