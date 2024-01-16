import {useMemo} from 'react';
import {TDraftDataStore} from 'types/DataStore.types';
import DataStoreService from 'services/DataStore.service';
import {useDataStore} from 'providers/DataStore/DataStore.provider';
import {Form, Typography} from 'antd';
import * as S from './TestConnectionStatus.styled';

export type TConnectionStatus = 'loading' | 'success' | 'error' | 'idle';

const Icon = ({status}: {status: TConnectionStatus}) => {
  switch (status) {
    case 'loading':
      return <S.LoadingIcon />;
    case 'success':
      return <S.IconSuccess />;
    case 'error':
      return <S.IconFail />;
    case 'idle':
    default:
      return <S.IconInfo />;
  }
};

const getText = (status: TConnectionStatus, shouldTestConnection: boolean) => {
  switch (status) {
    case 'loading':
      return shouldTestConnection ? 'Connecting to Data Store' : 'Waiting for incoming traces';
    case 'success':
      return shouldTestConnection ? 'Successful Connection' : 'We received your traces';
    case 'error':
      return shouldTestConnection ? 'Connection Failed' : 'We could not receive your traces';
    case 'idle':
    default:
      return 'Waiting for Test Connection';
  }
};

const TestConnectionStatus = () => {
  const {isTestConnectionLoading, testConnectionResponse} = useDataStore();
  const form = Form.useFormInstance<TDraftDataStore>();

  const shouldTest = useMemo(() => DataStoreService.shouldTestConnection(form.getFieldsValue()), [form]);

  const status = useMemo<TConnectionStatus>(() => {
    if (isTestConnectionLoading) return 'loading';
    if (!testConnectionResponse) return 'idle';

    return testConnectionResponse.allPassed ? 'success' : 'error';
  }, [isTestConnectionLoading, testConnectionResponse]);

  const text = getText(status, shouldTest);

  return (
    <S.StatusContainer>
      <Icon status={status} />
      <Typography.Text type="secondary">
        <b>{text}</b>
      </Typography.Text>
    </S.StatusContainer>
  );
};

export default TestConnectionStatus;
