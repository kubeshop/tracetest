import {Form} from 'antd';
import {useCallback} from 'react';
import * as Step from 'components/TestSuitePlugin/Step.styled';
import {useCreateTestSuite} from 'providers/CreateTestSuite';
import {TDraftTestSuite} from 'types/TestSuite.types';
import {ComponentNames} from 'constants/Plugins.constants';
import TestsSelectionForm from './TestsSelectionForm';

const TestsSelection = () => {
  const [form] = Form.useForm<TDraftTestSuite>();
  const {onNext} = useCreateTestSuite();

  const handleSubmit = useCallback(
    ({steps = []}: TDraftTestSuite) => {
      onNext({steps});
    },
    [onNext]
  );

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title $withSubtitle>Tests</Step.Title>
        <Step.Subtitle>Define the tests and their order of execution for this test suite</Step.Subtitle>
        <Form<TDraftTestSuite>
          id={ComponentNames.TestsSelection}
          autoComplete="off"
          data-cy="create-testsuite-modal"
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          onValuesChange={() => null}
        >
          <TestsSelectionForm />
        </Form>
      </Step.FormContainer>
    </Step.Step>
  );
};

export default TestsSelection;
