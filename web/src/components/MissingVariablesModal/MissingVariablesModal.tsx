import {useEffect} from 'react';
import {Form} from 'antd';
import VariablesService from 'services/Variables.service';
import {TDraftVariables, TTestVariablesMap} from 'types/Variables.types';
import {TEnvironmentValue} from 'types/Environment.types';
import * as S from './MissingVariablesModal.styled';
import MissingVariablesModalFooter from './MissingVariablesModalFooter';
import MissingVariablesModalForm from './MissingVariablesModalForm';
import useValidateVariablesDraft from './hooks/useValidateVariablesDraft';

interface IProps {
  isOpen: boolean;
  name: string;
  onClose(): void;
  onSubmit(values: TEnvironmentValue[]): void;
  variables: TTestVariablesMap;
}

const MissingVariablesModal = ({isOpen, onClose, onSubmit, variables, name}: IProps) => {
  const [form] = Form.useForm<TDraftVariables>();
  const {isValid, onValidate} = useValidateVariablesDraft();

  useEffect(() => {
    if (isOpen) {
      const draft = VariablesService.getDraftVariables(variables);
      onValidate({}, draft);
      form.setFieldsValue(draft);
    } else form.resetFields();
  }, [isOpen]);

  return (
    <S.Modal
      footer={<MissingVariablesModalFooter isValid={isValid} onCancel={onClose} onSave={() => form.submit()} />}
      onCancel={onClose}
      title={<S.Title>{name} - Undefined Variables</S.Title>}
      visible={isOpen}
      width={520}
    >
      <Form<TDraftVariables>
        autoComplete="off"
        form={form}
        layout="vertical"
        name="testOutput"
        onFinish={values => onSubmit(VariablesService.getFlatVariablesFromDraft(values))}
        onValuesChange={onValidate}
      >
        <S.Description>
          The following variables are referenced in this test but are not defined. Please provide a value to use for
          each of these missing variables.
        </S.Description>
        <MissingVariablesModalForm variables={variables} />
      </Form>
    </S.Modal>
  );
};

export default MissingVariablesModal;
