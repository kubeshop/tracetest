import {Button, Form} from 'antd';
import EditTestForm from 'components/EditTestForm';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {useTest} from 'providers/Test/Test.provider';
import {useCallback, useState} from 'react';
import {TDraftTest, TTest} from 'types/Test.types';
import {TestState} from 'constants/TestRun.constants';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import * as S from './EditTest.styled';

interface IProps {
  test: TTest;
}

const EditTest = ({test}: IProps) => {
  const {onEdit, isEditLoading} = useTest();
  const plugin = TriggerTypeToPlugin[test.trigger.type];
  const [isValid, setIsValid] = useState(true);

  const onValidate = useValidateTestDraft({pluginName: plugin.name, setIsValid});
  const [form] = Form.useForm<TDraftTest>();

  const handleOnSubmit = useCallback(
    async (values: TDraftTest) => {
      TestRunAnalyticsService.onTriggerEditSubmit();
      onEdit(values);
    },
    [onEdit]
  );

  const {run} = useTestRun();
  const stateIsFinished = ([TestState.FINISHED, TestState.FAILED] as string[]).includes(run.state);

  return (
    <S.Wrapper data-cy="edit-test-form">
      <S.FormContainer>
        <S.Title>Edit Test</S.Title>
        <EditTestForm form={form} test={test} onSubmit={handleOnSubmit} onValidation={onValidate} />
        <S.ButtonsContainer>
          <Button data-cy="edit-test-reset" onClick={() => form.resetFields()}>
            Reset
          </Button>
          <Button
            data-cy="edit-test-submit"
            loading={isEditLoading}
            disabled={!isValid || !stateIsFinished}
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

export default EditTest;
