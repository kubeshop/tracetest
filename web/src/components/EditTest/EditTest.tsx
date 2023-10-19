import {Button, Form} from 'antd';
import AllowButton, {Operation} from 'components/AllowButton';
import CreateButton from 'components/CreateButton';
import EditTestForm from 'components/EditTestForm';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {isRunStateFinished} from 'models/TestRun.model';
import {useTest} from 'providers/Test/Test.provider';
import {useCallback, useState} from 'react';
import {TDraftTest} from 'types/Test.types';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import Test from 'models/Test.model';
import * as S from './EditTest.styled';

interface IProps {
  test: Test;
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
  const stateIsFinished = isRunStateFinished(run.state);

  return (
    <S.Wrapper data-cy="edit-test-form">
      <S.FormContainer>
        <S.Title>Edit Test</S.Title>
        <EditTestForm form={form} test={test} onSubmit={handleOnSubmit} onValidation={onValidate} />
        <S.ButtonsContainer>
          <Button data-cy="edit-test-reset" onClick={() => form.resetFields()}>
            Reset
          </Button>
          <AllowButton
            operation={Operation.Edit}
            ButtonComponent={CreateButton}
            data-cy="edit-test-submit"
            disabled={!isValid || !stateIsFinished}
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

export default EditTest;
