import {Form, FormInstance, Input, Select} from 'antd';
import useValidate from 'components/CreateTestPlugins/Postman/steps/UploadCollection/hooks/useValidate';
import {IRequestDetailsValues} from 'components/CreateTestPlugins/Postman/steps/UploadCollection/UploadCollection';
import * as S from 'components/CreateTestPlugins/Rpc/steps/RequestDetails/RequestDetails.styled';
import RequestDetailsFileInput from 'components/CreateTestPlugins/Rpc/steps/RequestDetails/RequestDetailsFileInput';
import React, {useState} from 'react';
import RequestDetailsAuthInput from '../../../Rest/steps/RequestDetails/RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsHeadersInput from '../../../Rest/steps/RequestDetails/RequestDetailsHeadersInput';
import RequestDetailsUrlInput from '../../../Rest/steps/RequestDetails/RequestDetailsUrlInput';
import {useSelectTestFromCollectionChangeCallback} from './hooks/useSelectTestFromCollectionChangeCallback';
import {State, useUploadCollectionCallback} from './hooks/useUploadCollectionCallback';
import {useUploadEnvFileCallback} from './hooks/useUploadEnvFileCallback';

export const FORM_ID = 'upload-collection-test';

interface IProps {
  form: FormInstance<IRequestDetailsValues>;
  onSubmit(values: IRequestDetailsValues): void;
  onValidation(isValid: boolean): void;
}

const UploadCollectionForm = ({form, onSubmit, onValidation}: IProps) => {
  const handleOnValuesChange = useValidate(onValidation);
  const [state, setState] = useState<State>({
    requests: [],
  });
  return (
    <Form
      autoComplete="off"
      form={form}
      layout="vertical"
      name={FORM_ID}
      onFinish={onSubmit}
      onValuesChange={handleOnValuesChange}
    >
      <S.GlobalStyle />
      <div style={{display: 'grid'}}>
        <Form.Item data-cy="collectionFile" name="file" label="Upload Postman Collection">
          <RequestDetailsFileInput accept=".json" onChange={useUploadCollectionCallback(setState)} />
        </Form.Item>
        <Form.Item data-cy="envFile" name="envFile" label="Upload environment file (optional)">
          <RequestDetailsFileInput accept=".json" onChange={useUploadEnvFileCallback()} />
        </Form.Item>
        <Form.Item
          style={{width: 490}}
          data-cy="collectionTest"
          name="collectionTest"
          label="Select test from Postman Collection"
          initialValue={null}
        >
          <Select<string>
            data-cy="collection-test-select"
            placeholder="Select test from uploaded collection"
            onChange={useSelectTestFromCollectionChangeCallback(state, form)}
          >
            {state.requests.map(({id, name}) => (
              <Select.Option data-cy={`postman-test-${id}`} key={id} value={id}>
                {name}
              </Select.Option>
            ))}
          </Select>
        </Form.Item>
        <div style={{display: 'flex', paddingTop: 32}}>
          <span style={{flexBasis: '50%', paddingRight: 8}}>
            <RequestDetailsUrlInput />
          </span>
          <span style={{flexBasis: '50%', paddingLeft: 8}}>
            <RequestDetailsAuthInput form={form} />
          </span>
        </div>
        <div style={{display: 'flex'}}>
          <span style={{flexBasis: '50%', paddingRight: 8}}>
            <RequestDetailsHeadersInput />
          </span>
          <span style={{flexBasis: '50%', paddingLeft: 8}}>
            <Form.Item className="input-body" data-cy="body" label="Request body" name="body" style={{marginBottom: 0}}>
              <Input.TextArea placeholder="Enter request body text" />
            </Form.Item>
          </span>
        </div>
      </div>
    </Form>
  );
};

export default UploadCollectionForm;
