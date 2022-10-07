import {Form} from 'antd';
import EnvironmentForm from 'components/EnvironmentForm';
import {Dispatch, SetStateAction} from 'react';
import {CustomModal} from './EnvironmentModal.styled';
import {IEnvironment} from './IEnvironment';

interface IProps {
  setIsFormOpen: Dispatch<SetStateAction<boolean>>;
  isFormOpen: boolean;
  setEnvironment: Dispatch<SetStateAction<IEnvironment | undefined>>;
  environment?: IEnvironment;
}

export const EnvironmentModal = ({setEnvironment, setIsFormOpen, isFormOpen, environment}: IProps) => {
  const [form] = Form.useForm<IEnvironment>();
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
