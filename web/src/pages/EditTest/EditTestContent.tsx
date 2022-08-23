import {useNavigate} from 'react-router-dom';
import {useTour} from '@reactour/tour';
import {Button, Form} from 'antd';
import CreateTestHeader from 'components/CreateTestHeader';
import {useCallback} from 'react';
import {useEditTestMutation, useRunTestMutation} from 'redux/apis/TraceTest.api';
import {TDraftTest, TTest} from 'types/Test.types';
import EditTestForm from 'components/EditTestForm';
import TestService from 'services/Test.service';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import * as S from './EditTest.styled';

interface IProps {
  test: TTest;
}

const EditTestContent = ({test}: IProps) => {
  const navigate = useNavigate();

  const {setIsOpen} = useTour();
  const [editTest, {isLoading: isLoadingCreateTest}] = useEditTestMutation();
  const [runTest, {isLoading: isLoadingRunTest}] = useRunTestMutation();
  const plugin = TriggerTypeToPlugin[test.trigger.type];

  const {isValid, onValidate} = useValidateTestDraft({pluginName: plugin.name, isDefaultValid: true});

  const isLoading = isLoadingCreateTest || isLoadingRunTest;
  const [form] = Form.useForm<TDraftTest>();

  const handleOnSubmit = useCallback(
    async (values: TDraftTest) => {
      const rawTest = await TestService.getRequest(plugin, values, test);

      await editTest({
        test: rawTest,
        testId: test.id,
      }).unwrap();

      const run = await runTest({testId: test.id}).unwrap();
      setIsOpen(false);

      navigate(`/test/${test.id}/run/${run.id}`);
    },
    [editTest, navigate, plugin, runTest, setIsOpen, test]
  );

  return (
    <S.Wrapper data-cy="edit-test-form">
      <CreateTestHeader onBack={() => navigate('/')} title="Edit Test" />
      <S.FormContainer>
        <EditTestForm form={form} test={test} onSubmit={handleOnSubmit} onValidation={onValidate} />
        <S.ButtonsContainer>
          <Button
            data-cy="edit-test-submit"
            loading={isLoading}
            disabled={!isValid}
            type="primary"
            onClick={() => form.submit()}
          >
            Save
          </Button>
        </S.ButtonsContainer>
      </S.FormContainer>
    </S.Wrapper>
  );
};

export default EditTestContent;
