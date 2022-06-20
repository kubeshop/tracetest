import {Form, FormInstance, Input} from 'antd';

import {Steps} from 'components/GuidedTour/homeStepList';
import {HTTP_METHOD} from 'constants/Common.constants';
import React from 'react';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {RequestAuth} from '../../types/Common.types';
import {CreateTestFormAuthInput} from './CreateTestFormAuthInput/CreateTestFormAuthInput';
import {CreateTestFormDemoHelper} from './CreateTestFormDemoHelper';
import {CreateTestFormHeadersInput} from './CreateTestFormHeadersInput';
import {CreateTestFormUrlInput} from './CreateTestFormUrlInput';
import * as S from './CreateTestModal.styled';
import {useCreateTestFormHandleChangeCallback} from './useCreateTestFormHandleChangeCallback';

export const FORM_ID = 'create-test';

export interface ICreateTestValues {
  body: string;
  auth: RequestAuth;
  headers: {
    key: string;
    value: string;
  }[];
  method: HTTP_METHOD;
  name: string;
  url: string;
}

interface IProps {
  form: FormInstance<ICreateTestValues>;
  onSelectDemo(value: string): void;
  onSubmit(values: ICreateTestValues): Promise<void>;
  onValidation(isValid: boolean): void;
  selectedDemo: string;
}

const CreateTestForm = ({form, onSelectDemo, onSubmit, onValidation, selectedDemo}: IProps) => {
  const handleOnValuesChange = useCreateTestFormHandleChangeCallback(onValidation);
  return (
    <Form<any>
      autoComplete="off"
      data-cy="create-test-modal"
      form={form}
      layout="vertical"
      name={FORM_ID}
      onFinish={onSubmit}
      onValuesChange={handleOnValuesChange}
    >
      <S.GlobalStyle />
      <CreateTestFormDemoHelper
        selectedDemo={selectedDemo}
        form={form}
        onSelectDemo={onSelectDemo}
        onValidation={onValidation}
      />
      <CreateTestFormUrlInput />
      <Form.Item
        className="input-name"
        data-cy="name"
        data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Name)}
        label="Name"
        name="name"
        rules={[{required: true, message: 'Please enter a test name'}]}
      >
        <Input placeholder="Enter test name" />
      </Form.Item>
      <CreateTestFormHeadersInput />
      <Form.Item
        className="input-body"
        data-cy="body"
        data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Body)}
        label="Request body"
        name="body"
        style={{marginBottom: 0}}
      >
        <Input.TextArea placeholder="Enter request body text" />
      </Form.Item>
      <CreateTestFormAuthInput form={form} />
    </Form>
  );
};

export default CreateTestForm;
