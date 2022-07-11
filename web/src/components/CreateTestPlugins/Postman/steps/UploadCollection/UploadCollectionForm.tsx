import {Col, Form, FormInstance, Input, Row} from 'antd';
import useValidate from 'components/CreateTestPlugins/Postman/steps/UploadCollection/hooks/useValidate';
import {IRequestDetailsValues} from 'components/CreateTestPlugins/Postman/steps/UploadCollection/UploadCollection';
import React, {Dispatch, SetStateAction, useState} from 'react';
import RequestDetailsAuthInput from '../../../Rest/steps/RequestDetails/RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsHeadersInput from '../../../Rest/steps/RequestDetails/RequestDetailsHeadersInput';
import RequestDetailsUrlInput from '../../../Rest/steps/RequestDetails/RequestDetailsUrlInput';
import {CollectionFileField} from './fields/CollectionFileField';
import {EnvFileField} from './fields/EnvFileField';
import {useSelectTestCallback} from './hooks/useSelectTestCallback';
import {State} from './hooks/useUploadCollectionCallback';
import {SelectTestFromCollection} from './fields/SelectTestFromCollection';

export const FORM_ID = 'upload-collection-test';

interface IProps {
  form: FormInstance<IRequestDetailsValues>;
  setTransientUrl: Dispatch<SetStateAction<string>>;
  onSubmit(values: IRequestDetailsValues): void;
  onValidation(isValid: boolean): void;
}

const UploadCollectionForm = ({form, onSubmit, onValidation, setTransientUrl}: IProps) => {
  const handleOnValuesChange = useValidate(onValidation, setTransientUrl);
  const [state, setState] = useState<State>({requests: [], variables: []});

  return (
    <Form
      autoComplete="off"
      form={form}
      layout="vertical"
      name={FORM_ID}
      initialValues={{url: ''}}
      onFinish={onSubmit}
      onValuesChange={handleOnValuesChange}
    >
      <div style={{display: 'grid'}}>
        <CollectionFileField setState={setState} />
        <EnvFileField state={state} form={form} setState={setState} setTransientUrl={setTransientUrl} />
        <SelectTestFromCollection
          requests={state.requests}
          onChange={useSelectTestCallback(state, form, setTransientUrl)}
        />
        <Row gutter={12}>
          <Col span={12}>
            <RequestDetailsUrlInput />
          </Col>
          <Col span={12}>
            <Form.Item className="input-body" data-cy="body" label="Request body" name="body" style={{marginBottom: 0}}>
              <Input.TextArea placeholder="Enter request body text" />
            </Form.Item>
          </Col>
        </Row>
        <Row gutter={12}>
          <Col span={12}>
            <RequestDetailsHeadersInput />
          </Col>
          <Col span={12}>
            <RequestDetailsAuthInput form={form} />
          </Col>
        </Row>
      </div>
    </Form>
  );
};

export default UploadCollectionForm;
