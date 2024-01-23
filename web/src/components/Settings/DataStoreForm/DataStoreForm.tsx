import {Form} from 'antd';
import {useCallback, useEffect, useMemo} from 'react';
import DataStoreService from 'services/DataStore.service';
import {TDraftDataStore, TDataStoreForm, SupportedDataStores} from 'types/DataStore.types';
import DataStoreConfig from 'models/DataStoreConfig.model';
import DataStoreConfiguration from 'components/DataStoreConfiguration/DataStoreConfiguration';
import {DataStoreSelection} from 'components/Inputs';
import * as S from './DataStoreForm.styled';

export const FORM_ID = 'data-store';

interface IProps {
  form: TDataStoreForm;
  dataStoreConfig: DataStoreConfig;
  onSubmit(values: TDraftDataStore): Promise<void>;
  onIsFormValid(isValid: boolean): void;
  onTestConnection(): void;
  isLoading: boolean;
  isTestConnectionSuccess: boolean;
}

const DataStoreForm = ({
  form,
  onSubmit,
  dataStoreConfig,
  onIsFormValid,
  onTestConnection,
  isLoading,
  isTestConnectionSuccess,
}: IProps) => {
  const configuredDataStoreType = dataStoreConfig.defaultDataStore.type as SupportedDataStores;
  const initialValues = useMemo(
    () => DataStoreService.getInitialValues(dataStoreConfig.defaultDataStore, configuredDataStoreType),
    [configuredDataStoreType, dataStoreConfig]
  );
  const dataStoreType = Form.useWatch('dataStoreType', form);

  useEffect(() => {
    form.setFieldsValue({
      dataStore: {name: '', type: SupportedDataStores.JAEGER, ...initialValues.dataStore},
    });
  }, [dataStoreType, form, initialValues.dataStore]);

  const onValidation = useCallback(
    async (_: any, draft: TDraftDataStore) => {
      try {
        const isValid = await DataStoreService.validateDraft(draft);
        onIsFormValid(isValid);
      } catch (e) {
        onIsFormValid(false);
      }
    },
    [onIsFormValid]
  );

  return (
    <Form<TDraftDataStore>
      autoComplete="off"
      form={form}
      layout="vertical"
      name={FORM_ID}
      onFinish={onSubmit}
      onValuesChange={onValidation}
      initialValues={initialValues}
    >
      <S.FormContainer>
        <Form.Item name="dataStoreType">
          <DataStoreSelection />
        </Form.Item>
        <S.FactoryContainer>
          <DataStoreConfiguration
            isTestConnectionSuccess={isTestConnectionSuccess}
            onSubmit={() => form.submit()}
            onTestConnection={onTestConnection}
            isSubmitLoading={isLoading}
            dataStoreType={dataStoreType ?? SupportedDataStores.JAEGER}
          />
        </S.FactoryContainer>
      </S.FormContainer>
    </Form>
  );
};

export default DataStoreForm;
