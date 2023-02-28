import {Form} from 'antd';
import {useEffect, useState} from 'react';

import EnvironmentForm from 'components/EnvironmentForm';
import Environment from 'models/Environment.model';
import * as S from './EnvironmentModal.styled';
import EnvironmentModalFooter from './EnvironmentModalFooter';

interface IProps {
  environment?: Environment;
  isOpen: boolean;
  isLoading: boolean;
  onClose(): void;
  onSubmit(environment: Environment): void;
}

export const DEFAULT_VALUES = [{key: '', value: ''}];

const EnvironmentModal = ({environment, isOpen, onClose, onSubmit, isLoading}: IProps) => {
  const [form] = Form.useForm<Environment>();
  const [isFormValid, setIsFormValid] = useState(false);
  const isEditing = Boolean(environment);

  useEffect(() => {
    if (environment && isOpen) form.setFieldsValue(environment);
    if (!isOpen || !environment) {
      form.resetFields();
      form.setFieldsValue({values: [{key: '', value: ''}]});
    }
  }, [environment, form, isOpen]);

  const handleOnValidate = (changedValues: any, {name, description, values}: Environment) => {
    setIsFormValid(Boolean(name) && Boolean(description) && Boolean(values.length));
  };

  const handleOnSubmit = async (values: Environment) => {
    onSubmit(values);
    onClose();
  };

  return (
    <S.Modal
      cancelText="Cancel"
      footer={
        <EnvironmentModalFooter
          isEditing={isEditing}
          isLoading={isLoading}
          isValid={isFormValid}
          onCancel={onClose}
          onSave={() => form.submit()}
        />
      }
      onCancel={onClose}
      title={<S.Title>{isEditing ? 'Edit Environment' : 'Create Environment'}</S.Title>}
      visible={isOpen}
    >
      <EnvironmentForm
        form={form}
        initialValues={environment}
        onSubmit={handleOnSubmit}
        onValidate={handleOnValidate}
      />
    </S.Modal>
  );
};

export default EnvironmentModal;
