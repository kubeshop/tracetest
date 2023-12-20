import {Button, Form, Space, Typography} from 'antd';
import {useCallback, useEffect, useMemo} from 'react';
import AllowButton, {Operation} from 'components/AllowButton';
import DataStoreIcon from 'components/DataStoreIcon/DataStoreIcon';
import DataStoreService from 'services/DataStore.service';
import {TDraftDataStore, TDataStoreForm, SupportedDataStores} from 'types/DataStore.types';
import DataStoreConfig from 'models/DataStoreConfig.model';
import {DataStoreSelection} from 'components/Inputs';
import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import DataStoreComponentFactory from '../DataStorePlugin/DataStoreComponentFactory';
import * as S from './DataStoreForm.styled';

export const FORM_ID = 'data-store';

interface IProps {
  form: TDataStoreForm;
  dataStoreConfig: DataStoreConfig;
  onSubmit(values: TDraftDataStore): Promise<void>;
  onIsFormValid(isValid: boolean): void;
  onTestConnection(): void;
  isTestConnectionLoading: boolean;
  isLoading: boolean;
  isFormValid: boolean;
}

const DataStoreForm = ({
  form,
  onSubmit,
  dataStoreConfig,
  onIsFormValid,
  onTestConnection,
  isTestConnectionLoading,
  isLoading,
  isFormValid,
}: IProps) => {
  const configuredDataStoreType = dataStoreConfig.defaultDataStore.type as SupportedDataStores;
  const initialValues = useMemo(
    () => DataStoreService.getInitialValues(dataStoreConfig, configuredDataStoreType),
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
          <S.TopContainer>
            <Space>
              <DataStoreIcon dataStoreType={dataStoreType ?? SupportedDataStores.JAEGER} width="22" height="22" />

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
            <Button loading={isTestConnectionLoading} type="primary" ghost onClick={onTestConnection}>
              Test Connection
            </Button>
            <AllowButton
              operation={Operation.Configure}
              disabled={!isFormValid}
              loading={isLoading}
              type="primary"
              onClick={() => form.submit()}
            >
              Save and Set as DataStore
            </AllowButton>
          </S.ButtonsContainer>
        </S.FactoryContainer>
      </S.FormContainer>
    </Form>
  );
};

export default DataStoreForm;
