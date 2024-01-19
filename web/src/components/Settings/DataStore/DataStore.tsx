import {Form} from 'antd';
import {useDataStore} from 'providers/DataStore/DataStore.provider';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import {useCallback} from 'react';
import {TDraftDataStore} from 'types/DataStore.types';
import DataStoreForm from '../DataStoreForm';
import * as S from './DataStore.styled';

const DataStore = () => {
  const {dataStoreConfig} = useSettingsValues();
  const {isLoading, onIsFormValid, onSaveConfig, onTestConnection, isTestConnectionSuccessful} = useDataStore();
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
          isTestConnectionSuccess={isTestConnectionSuccessful}
          form={form}
          dataStoreConfig={dataStoreConfig}
          onSubmit={handleOnSubmit}
          onTestConnection={handleTestConnection}
          isLoading={isLoading}
          onIsFormValid={onIsFormValid}
        />
      </S.FormContainer>
    </S.Wrapper>
  );
};

export default DataStore;
