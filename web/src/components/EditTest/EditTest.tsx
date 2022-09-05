import {Button, Form} from 'antd';
import {useCallback} from 'react';
import {TDraftTest, TTest} from 'types/Test.types';
import EditTestForm from 'components/EditTestForm';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import {useTest} from 'providers/Test/Test.provider';
import * as S from './EditTest.styled';

interface IProps {
  test: TTest;
}

const EditTest = ({test}: IProps) => {
  const {onEdit, isEditLoading} = useTest();
  const plugin = TriggerTypeToPlugin[test.trigger.type];

  const {isValid, onValidate} = useValidateTestDraft({pluginName: plugin.name, isDefaultValid: true});

  const [form] = Form.useForm<TDraftTest>();

  const handleOnSubmit = useCallback(
    async (values: TDraftTest) => {
      onEdit(values);
    },
    [onEdit]
  );

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
            disabled={!isValid}
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
