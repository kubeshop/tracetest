import {Form, Modal} from 'antd';
import {Dispatch, SetStateAction} from 'react';
import styled from 'styled-components';
import {EnvironmentForm} from './EnvironmentForm';
import {EnvironmentState} from './EnvironmentState';
import {IEnvironment} from './IEnvironment';

const CustomModal = styled(Modal)`
  && {
    .ant-modal-footer {
      display: none;
    }
  }
`;

interface IProps {
  state: EnvironmentState;
  setState: Dispatch<SetStateAction<EnvironmentState>>;
}

export const EnvsModal: React.FC<IProps> = ({state, setState}) => {
  const [form] = Form.useForm<IEnvironment>();
  const onCancel = () => {
    setState(st => ({...st, dialog: false, environment: undefined}));
    form.setFieldsValue({name: '', description: '', variables: []});
  };
  return (
    <CustomModal
      cancelText="Cancel"
      onCancel={onCancel}
      title="Create Environment"
      visible={state.dialog}
      data-cy="delete-confirmation-modal"
      footer={[]}
    >
      <EnvironmentForm onCancel={onCancel} form={form} state={state} />
    </CustomModal>
  );
};
