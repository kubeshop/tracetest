import {Form} from 'antd';
import {useCallback, useEffect} from 'react';

import * as Step from 'components/CreateTestPlugins/Step.styled';
import {ComponentNames} from 'constants/Plugins.constants';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {ITraceIDValues} from 'types/Test.types';
import ValueForm from './VariableNameForm';

const VariableName = () => {
  const {onNext, draftTest, onIsFormValid} = useCreateTest();
  const [form] = Form.useForm<ITraceIDValues>();

  const handleSubmit = useCallback(
    (values: ITraceIDValues) => {
      onNext(values);
    },
    [onNext]
  );

  const onRefreshData = useCallback(async () => {
    const {id} = draftTest as ITraceIDValues;
    form.setFieldsValue({id});

    try {
      form.validateFields();
      onIsFormValid(true);
    } catch (err) {
      onIsFormValid(false);
    }
  }, [draftTest, form, onIsFormValid]);

  useEffect(() => {
    onRefreshData();
  }, [onRefreshData]);

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Enter the variable name to use for the Trace Id</Step.Title>
        <Step.Subtitle>Please enter a variable name for the trace id or accept the default &apos;traceId&apos; name</Step.Subtitle>
        <Form<ITraceIDValues>
          id={ComponentNames.TraceIdVariableName}
          autoComplete="off"
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <ValueForm />
        </Form>
      </Step.FormContainer>
    </Step.Step>
  );
};

export default VariableName;
