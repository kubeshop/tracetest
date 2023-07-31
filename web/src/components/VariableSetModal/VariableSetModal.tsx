import {Form} from 'antd';
import {useEffect, useState} from 'react';

import VariableSetForm from 'components/VariableSet';
import VariableSet from 'models/VariableSet.model';
import VariableSetService from 'services/VariableSet.service';
import * as S from './VariableSetModal.styled';
import VariableSetModalFooter from './VariableSetModalFooter';

interface IProps {
  variableSet?: VariableSet;
  isOpen: boolean;
  isLoading: boolean;
  onClose(): void;
  onSubmit(variableSet: VariableSet): void;
}

export const DEFAULT_VALUES = [{key: '', value: ''}];

const VariableSetModal = ({variableSet, isOpen, onClose, onSubmit, isLoading}: IProps) => {
  const [form] = Form.useForm<VariableSet>();
  const [isFormValid, setIsFormValid] = useState(false);
  const isEditing = Boolean(variableSet);

  useEffect(() => {
    if (variableSet && isOpen) form.setFieldsValue(variableSet);
    if (!isOpen || !variableSet) {
      form.resetFields();
      form.setFieldsValue({values: [{key: '', value: ''}]});
    }
  }, [variableSet, form, isOpen]);

  const handleOnValidate = (changedValues: any, draft: VariableSet) => {
    setIsFormValid(VariableSetService.validateDraft(draft));
  };

  const handleOnSubmit = async (values: VariableSet) => {
    onSubmit(values);
    onClose();
  };

  return (
    <S.Modal
      cancelText="Cancel"
      footer={
        <VariableSetModalFooter
          isEditing={isEditing}
          isLoading={isLoading}
          isValid={isFormValid}
          onCancel={onClose}
          onSave={() => form.submit()}
        />
      }
      onCancel={onClose}
      title={<S.Title>{isEditing ? 'Edit Variable Set' : 'Create Variable Set'}</S.Title>}
      visible={isOpen}
    >
      <VariableSetForm
        form={form}
        initialValues={variableSet}
        onSubmit={handleOnSubmit}
        onValidate={handleOnValidate}
      />
    </S.Modal>
  );
};

export default VariableSetModal;
