import {Form} from 'antd';
import {useDataStore} from 'providers/DataStore/DataStore.provider';
import {useDataStoreConfig} from 'providers/DataStoreConfig/DataStoreConfig.provider';
import {useCallback} from 'react';
import {TDraftDataStore, ConfigMode} from 'types/DataStore.types';
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
        <DataStoreForm
          form={form}
          dataStoreConfig={dataStoreConfig}
          onSubmit={handleOnSubmit}
          onTestConnection={handleTestConnection}
          isConfigReady={isConfigReady}
          isTestConnectionLoading={isTestConnectionLoading}
          onDeleteConfig={onDeleteConfig}
          isLoading={isLoading}
          isFormValid={isFormValid}
          onIsFormValid={onIsFormValid}
        />
      </S.FormContainer>
    </S.Wrapper>
  );
};

export default DataStore;
