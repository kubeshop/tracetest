import {Button, Space} from 'antd';
import DataStoreIcon from 'components/DataStoreIcon/DataStoreIcon';
import {SupportedDataStores} from 'types/DataStore.types';
import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import DataStoreComponentFactory from 'components/Settings/DataStorePlugin/DataStoreComponentFactory';
import AllowButton, {Operation} from 'components/AllowButton';
import * as S from './DataStoreConfiguration.styled';
import TestConnectionStatus from '../TestConnectionStatus';

interface IProps {
  isTestConnectionLoading: boolean;
  isTestConnectionSuccess?: boolean;
  isSubmitLoading: boolean;
  isValid: boolean;
  dataStoreType: SupportedDataStores;
  onSubmit(): void;
  onTestConnection(): void;
  isWizard?: boolean;
}

const DataStoreConfiguration = ({
  onSubmit,
  onTestConnection,
  isSubmitLoading,
  isTestConnectionLoading,
  isTestConnectionSuccess,
  isValid,
  dataStoreType,
  isWizard = false,
}: IProps) => {
  return (
    <>
      <S.TopContainer>
        <Space align="start">
          <DataStoreIcon
            withColor={isWizard}
            dataStoreType={dataStoreType ?? SupportedDataStores.JAEGER}
            width="22"
            height="22"
          />
          <S.Title level={2}>{SupportedDataStoresToName[dataStoreType ?? SupportedDataStores.JAEGER]}</S.Title>
        </Space>
        {!isWizard && (
          <S.Description>
            Tracetest needs configuration information to be able to retrieve your trace from your distributed tracing
            solution. Select your Tracing Backend and enter the configuration info.
          </S.Description>
        )}
        {dataStoreType && <DataStoreComponentFactory dataStoreType={dataStoreType} />}
      </S.TopContainer>

      <S.ButtonsContainer>
        <TestConnectionStatus />
        <Button loading={isTestConnectionLoading} type="primary" ghost onClick={onTestConnection}>
          Test Connection
        </Button>
        <AllowButton
          operation={Operation.Configure}
          disabled={!isValid || (isWizard && !isTestConnectionSuccess)}
          loading={isSubmitLoading}
          type="primary"
          onClick={onSubmit}
        >
          {isWizard ? 'Continue' : 'Save'}
        </AllowButton>
      </S.ButtonsContainer>
    </>
  );
};

export default DataStoreConfiguration;
