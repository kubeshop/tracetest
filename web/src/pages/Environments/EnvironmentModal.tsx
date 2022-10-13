import {Form} from 'antd';
import EnvironmentForm from 'components/EnvironmentForm';
import {Dispatch, SetStateAction} from 'react';
import {CustomModal} from './EnvironmentModal.styled';
import {TEnvironment} from '../../types/Environment.types';

interface IProps {
  setIsFormOpen: Dispatch<SetStateAction<boolean>>;
  isFormOpen: boolean;
  setEnvironment: Dispatch<SetStateAction<TEnvironment | undefined>>;
  environment?: TEnvironment;
}

export const EnvironmentModal = ({setEnvironment, setIsFormOpen, isFormOpen, environment}: IProps) => {
  const [form] = Form.useForm<TEnvironment>();
  const onCancel = () => {
    form.setFieldsValue({variables: []});
    setEnvironment(undefined);
    setIsFormOpen(false);
    form.resetFields(['name', 'description', 'variables', 'id']);
  };
  return (
    <CustomModal
      cancelText="Cancel"
      onCancel={onCancel}
      title="Create Environment"
      visible={isFormOpen}
      data-cy="delete-confirmation-modal"
      footer={[]}
    >
      <EnvironmentForm onCancel={onCancel} form={form} environment={environment} />
    </CustomModal>
  );
};
