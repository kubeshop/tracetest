import {Form} from 'antd';
import {useEffect, useState} from 'react';

import EnvironmentForm from 'components/EnvironmentForm';
import {useCreateEnvironmentMutation, useUpdateEnvironmentMutation} from 'redux/apis/TraceTest.api';
import {TEnvironment} from 'types/Environment.types';
import * as S from './EnvironmentModal.styled';
import EnvironmentModalFooter from './EnvironmentModalFooter';

interface IProps {
  environment?: TEnvironment;
  isOpen: boolean;
  onClose(): void;
}

export const DEFAULT_VALUES = [{key: '', value: ''}];

export const EnvironmentModal = ({environment, isOpen, onClose}: IProps) => {
  const [form] = Form.useForm<TEnvironment>();
  const [createEnvironment, {isLoading}] = useCreateEnvironmentMutation();
  const [updateEnvironment, {isLoading: isLoadingUpdate}] = useUpdateEnvironmentMutation();
  const [isFormValid, setIsFormValid] = useState(false);
  const isEditing = Boolean(environment);

  useEffect(() => {
    if (environment && isOpen) form.setFieldsValue(environment);
    if (!isOpen || !environment) {
      form.resetFields();
      form.setFieldsValue({values: [{key: '', value: ''}]});
    }
  }, [environment, form, isOpen]);

  const handleOnValidate = (changedValues: any, {name, description, values}: TEnvironment) => {
    setIsFormValid(Boolean(name) && Boolean(description) && Boolean(values.length));
  };

  const handleOnSubmit = async (values: TEnvironment) => {
    if (isEditing) {
      await updateEnvironment({environment: values, environmentId: environment?.id ?? ''});
    } else {
      await createEnvironment(values);
    }
    onClose();
  };

  return (
    <S.Modal
      cancelText="Cancel"
      footer={
        <EnvironmentModalFooter
          isEditing={isEditing}
          isLoading={isLoading || isLoadingUpdate}
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
