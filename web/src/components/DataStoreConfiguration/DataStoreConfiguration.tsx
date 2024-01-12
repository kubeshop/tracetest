import {Button, Space, Typography} from 'antd';
import DataStoreIcon from 'components/DataStoreIcon/DataStoreIcon';
import {SupportedDataStores} from 'types/DataStore.types';
import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import DataStoreComponentFactory from 'components/Settings/DataStorePlugin/DataStoreComponentFactory';
import AllowButton, {Operation} from 'components/AllowButton';
import * as S from './DataStoreConfiguration.styled';
import TestConnectionStatus from '../TestConnectionStatus';

interface IProps {
  isTestConnectionLoading: boolean;
  isSubmitLoading: boolean;
  isValid: boolean;
  dataStoreType: SupportedDataStores;
  withColor?: boolean;

  onSubmit(): void;
  onTestConnection(): void;
}

const DataStoreConfiguration = ({
  onSubmit,
  onTestConnection,
  isSubmitLoading,
  isTestConnectionLoading,
  isValid,
  dataStoreType,
  withColor = false,
}: IProps) => {
  return (
    <>
      <S.TopContainer>
        <Space>
          <DataStoreIcon
            withColor={withColor}
            dataStoreType={dataStoreType ?? SupportedDataStores.JAEGER}
            width="22"
            height="22"
          />

          <Typography.Title level={2}>
            {SupportedDataStoresToName[dataStoreType ?? SupportedDataStores.JAEGER]}
          </Typography.Title>
        </Space>

        <S.Description>
          Tracetest needs configuration information to be able to retrieve your trace from your distributed tracing
          solution. Select your Tracing Backend and enter the configuration info.
        </S.Description>
        {dataStoreType && <DataStoreComponentFactory dataStoreType={dataStoreType} />}
      </S.TopContainer>
      <S.ButtonsContainer>
        <TestConnectionStatus />
        <Button loading={isTestConnectionLoading} type="primary" ghost onClick={onTestConnection}>
          Test Connection
        </Button>
        <AllowButton
          operation={Operation.Configure}
          disabled={!isValid}
          loading={isSubmitLoading}
          type="primary"
          onClick={onSubmit}
        >
          Save
        </AllowButton>
      </S.ButtonsContainer>
    </>
  );
};

export default DataStoreConfiguration;
