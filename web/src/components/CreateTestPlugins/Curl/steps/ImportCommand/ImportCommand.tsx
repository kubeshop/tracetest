import {Form} from 'antd';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {ComponentNames} from 'constants/Plugins.constants';
import {useCallback, useEffect} from 'react';
import {ICurlValues} from 'types/Test.types';
import CurlService from 'services/Triggers/Curl.service';
import ImportCommandForm from './ImportCommandForm';

const ImportCommand = () => {
  const {onNext, draftTest, onIsFormValid} = useCreateTest();
  const [form] = Form.useForm<ICurlValues>();

  const handleSubmit = useCallback(
    ({command}: ICurlValues) => {
      const draft = CurlService.getRequestFromCommand(command);

      onNext(draft);
    },
    [onNext]
  );

  useEffect(() => {
    const {command} = draftTest as ICurlValues;
    form.setFieldsValue({command});
    onIsFormValid(CurlService.getIsValidCommand(command));
  }, [draftTest, form, onIsFormValid]);

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>CURL command</Step.Title>
        <Form<ICurlValues>
          id={ComponentNames.ImportCommand}
          autoComplete="off"
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          onValuesChange={(_: any, {command}: ICurlValues) => {
            onIsFormValid(CurlService.getIsValidCommand(command));
          }}
        >
          <ImportCommandForm />
        </Form>
      </Step.FormContainer>
    </Step.Step>
  );
};

export default ImportCommand;
