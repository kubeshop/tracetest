import {Button, Form} from 'antd';
import {useCallback, useState} from 'react';
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
          <Button
            data-cy="edit-testsuite-submit"
            loading={isEditLoading}
            disabled={!isFormValid || !stateIsFinished}
            type="primary"
            onClick={() => form.submit()}
          >
            Save & Run
          </Button>
        </S.ButtonsContainer>
      </S.FormContainer>
    </S.Wrapper>
  );
};

export default EditTestSuite;
