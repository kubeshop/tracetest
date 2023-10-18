import {Button, Form} from 'antd';
import {useCallback, useState} from 'react';
import AllowButton, {Operation} from 'components/AllowButton';
import CreateButton from 'components/CreateButton';
import {TDraftTestSuite} from 'types/TestSuite.types';
import {useTestSuite} from 'providers/TestSuite';
import useValidateTestSuiteDraft from 'hooks/useValidateTestSuiteDraft';
import {TestState} from 'constants/TestRun.constants';
import TestSuiteRun from 'models/TestSuiteRun.model';
import * as S from './EditTestSuite.styled';
import EditTestSuiteForm from '../EditTestSuiteForm';

interface IProps {
  testSuite: TDraftTestSuite;
  testSuiteRun: TestSuiteRun;
}

const EditTestSuite = ({testSuite, testSuiteRun}: IProps) => {
  const [form] = Form.useForm<TDraftTestSuite>();
  const {onEdit, isEditLoading} = useTestSuite();
  const [isFormValid, setIsFormValid] = useState(true);
  const stateIsFinished = ([TestState.FINISHED, TestState.FAILED] as string[]).includes(testSuiteRun.state);

  const onChange = useValidateTestSuiteDraft({setIsValid: setIsFormValid});

  const handleOnSubmit = useCallback(
    async (values: TDraftTestSuite) => {
      onEdit(values);
    },
    [onEdit]
  );

  return (
    <S.Wrapper data-cy="edit-testsuite-form">
      <S.FormContainer>
        <S.Title>Edit Test Suite</S.Title>
        <EditTestSuiteForm form={form} testSuite={testSuite} onSubmit={handleOnSubmit} onValidation={onChange} />
        <S.ButtonsContainer>
          <Button data-cy="edit-testsuite-reset" onClick={() => form.resetFields()}>
            Reset
          </Button>
          <AllowButton
            operation={Operation.Edit}
            ButtonComponent={CreateButton}
            data-cy="edit-testsuite-submit"
            disabled={!isFormValid || !stateIsFinished}
            loading={isEditLoading}
            onClick={() => form.submit()}
            type="primary"
          >
            Save & Run
          </AllowButton>
        </S.ButtonsContainer>
      </S.FormContainer>
    </S.Wrapper>
  );
};

export default EditTestSuite;
