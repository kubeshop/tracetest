import {Form, Modal} from 'antd';
import {Dispatch, SetStateAction, useCallback, useState} from 'react';
import styled from 'styled-components';
import SearchInput from '../../components/SearchInput';
import {IEnvironment} from '../../redux/apis/TraceTest.api';
import {EnvironmentForm} from './EnvironmentForm';
import {EnvironmentState} from './EnvironmentState';
import EnvList from './EnvList';
import * as S from './Envs.styled';
import EnvsActions from './EnvsActions';

const CustomModal = styled(Modal)`
  && {
    .ant-modal-footer {
      display: none;
    }
  }
`;

const EnvsModal = ({
  state,
  setState,
}: {
  state: EnvironmentState;
  setState: Dispatch<SetStateAction<EnvironmentState>>;
}) => {
  const [form] = Form.useForm<IEnvironment>();

  const onCancel = () => {
    form.setFieldsValue({variables: []});
    setState(st => ({...st, environment: undefined, dialog: false}));
    form.resetFields(['name', 'description', 'variables', 'id']);
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

const EnvsContent: React.FC = () => {
  const [state, setState] = useState<EnvironmentState>({query: '', dialog: false, environment: undefined});
  const onSearch = useCallback((value: string) => setState(st => ({...st, query: value})), [setState]);
  const openDialog = useCallback((dialog: boolean) => setState(st => ({...st, dialog})), [setState]);

  return (
    <S.Wrapper>
      <S.TitleText>All Envs</S.TitleText>
      <S.PageHeader>
        <SearchInput onSearch={onSearch} placeholder="Search test" />
        <EnvsActions openDialog={() => openDialog(true)} />
      </S.PageHeader>
      <EnvList
        query={state.query}
        openDialog={openDialog}
        setEnvironment={(environment: IEnvironment) => setState(st => ({...st, environment}))}
      />
      <EnvsModal setState={setState} state={state} />
    </S.Wrapper>
  );
};

export default EnvsContent;
