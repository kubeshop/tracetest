import {Button, Form} from 'antd';
import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import {useDataStore} from 'providers/DataStore/DataStore.provider';
import {useDataStoreConfig} from 'providers/DataStoreConfig/DataStoreConfig.provider';
import {useCallback} from 'react';
import {TDraftDataStore, ConfigMode} from 'types/Config.types';
import DataStoreForm from '../DataStoreForm';
import * as S from './DataStore.styled';

const DataStore = () => {
  const {dataStoreConfig} = useDataStoreConfig();
  const {
    isLoading,
    isFormValid,
    onIsFormValid,
    onSaveConfig,
    isTestConnectionLoading,
    onTestConnection,
    onDeleteConfig,
  } = useDataStore();
  const isConfigReady = dataStoreConfig.mode === ConfigMode.READY;
  const [form] = Form.useForm<TDraftDataStore>();

  const handleOnSubmit = useCallback(
    async (values: TDraftDataStore) => {
      onSaveConfig(values, dataStoreConfig.defaultDataStore);
    },
    [onSaveConfig, dataStoreConfig.defaultDataStore]
  );

  const handleTestConnection = useCallback(async () => {
    const draft = form.getFieldsValue();
    onTestConnection(draft, dataStoreConfig.defaultDataStore);
  }, [form, onTestConnection, dataStoreConfig.defaultDataStore]);

  return (
    <S.Wrapper>
      <S.FormContainer>
        <div>
          <S.Description>
            Tracetest needs configuration information to be able to retrieve your trace from your distributed tracing
            solution. Select your tracing data store and enter the configuration info.
          </S.Description>
          <S.Title>Choose your OpenTelemetry data store</S.Title>
          <DataStoreForm
            form={form}
            dataStoreConfig={dataStoreConfig}
            onSubmit={handleOnSubmit}
            onIsFormValid={onIsFormValid}
          />
        </div>
        <S.ButtonsContainer>
          {isConfigReady ? (
            <Button
              disabled={isLoading}
              type="primary"
              ghost
              onClick={() => onDeleteConfig(dataStoreConfig.defaultDataStore)}
              danger
            >
              {`Delete ${SupportedDataStoresToName[dataStoreConfig.defaultDataStore.type]} Data Store`}
            </Button>
          ) : (
            <div />
          )}
          <S.SaveContainer>
            <Button
              loading={isTestConnectionLoading}
              disabled={!isFormValid}
              type="primary"
              ghost
              onClick={handleTestConnection}
            >
              Test Connection
            </Button>
            <Button disabled={!isFormValid} loading={isLoading} type="primary" onClick={() => form.submit()}>
              Save
            </Button>
          </S.SaveContainer>
        </S.ButtonsContainer>
      </S.FormContainer>
    </S.Wrapper>
  );
};

export default DataStore;
