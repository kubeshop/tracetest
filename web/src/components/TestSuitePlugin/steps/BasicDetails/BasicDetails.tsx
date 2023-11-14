import {Form} from 'antd';
import {useCallback} from 'react';
import * as Step from 'components/TestPlugins/Step.styled';
import {useCreateTestSuite} from 'providers/CreateTestSuite';
import {TDraftTestSuite} from 'types/TestSuite.types';
import BasicDetailsForm from 'components/BasicDetailsForm/BasicDetailsForm';
import {ComponentNames} from 'constants/Plugins.constants';

const BasicDetails = () => {
  const [form] = Form.useForm<TDraftTestSuite>();
  const {onNext, onIsFormValid} = useCreateTestSuite();

  const handleSubmit = useCallback(
    ({name, description}: TDraftTestSuite) => {
      onNext({name, description});
    },
    [onNext]
  );

  const onValidate = useCallback(
    async (changedValues, draft: TDraftTestSuite) => {
      onIsFormValid(!!draft.name);
    },
    [onIsFormValid]
  );

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title $withSubtitle>Name your Test Suite</Step.Title>
        <Step.Subtitle>
          Create a test suite which will run several tests in sequence. You can set outputs into variables and use these
          variables later in the test suite in other tests.{' '}
        </Step.Subtitle>
        <Form<TDraftTestSuite>
          id={ComponentNames.BasicDetails}
          autoComplete="off"
          data-cy="create-test-modal"
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          onValuesChange={onValidate}
        >
          <BasicDetailsForm />
        </Form>
      </Step.FormContainer>
    </Step.Step>
  );
};

export default BasicDetails;
