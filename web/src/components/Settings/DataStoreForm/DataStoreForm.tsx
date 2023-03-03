import {Button, Form} from 'antd';
import {useCallback, useEffect, useMemo} from 'react';
import SetupConfigService from 'services/DataStore.service';
import {TDraftDataStore, TDataStoreForm, SupportedDataStores} from 'types/DataStore.types';
import DataStore from 'models/DataStore.model';
import {SupportedDataStoresToExplanation, SupportedDataStoresToName} from 'constants/DataStore.constants';
import DataStoreConfig from 'models/DataStoreConfig.model';
import DataStoreDocsBanner from '../DataStoreDocsBanner/DataStoreDocsBanner';
import DataStoreComponentFactory from '../DataStorePlugin/DataStoreComponentFactory';
import * as S from './DataStoreForm.styled';
import DataStoreSelectionInput from './DataStoreSelectionInput';

export const FORM_ID = 'data-store';

interface IProps {
  form: TDataStoreForm;
  dataStoreConfig: DataStoreConfig;
  onSubmit(values: TDraftDataStore): Promise<void>;
  onIsFormValid(isValid: boolean): void;
  onTestConnection(): void;
  isConfigReady: boolean;
  isTestConnectionLoading: boolean;
  onDeleteConfig(dataStore: DataStore): void;
  isLoading: boolean;
  isFormValid: boolean;
}

const DataStoreForm = ({
  form,
  onSubmit,
  dataStoreConfig,
  onIsFormValid,
  onTestConnection,
  isConfigReady,
  isTestConnectionLoading,
  onDeleteConfig,
  isLoading,
  isFormValid,
}: IProps) => {
  const initialValues = useMemo(() => SetupConfigService.getInitialValues(dataStoreConfig), [dataStoreConfig]);
  const dataStoreType = Form.useWatch('dataStoreType', form);

  useEffect(() => {
    form.setFieldsValue({dataStore: {name: '', type: SupportedDataStores.JAEGER, ...initialValues.dataStore}});
  }, [dataStoreType, form, initialValues.dataStore]);

  const onValidation = useCallback(
    async (_: any, draft: TDraftDataStore) => {
      const isValid = await SetupConfigService.validateDraft(draft);
      onIsFormValid(isValid);
    },
    [onIsFormValid]
  );
  const explanation = SupportedDataStoresToExplanation[dataStoreType!];

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
          <DataStoreSelectionInput />
        </Form.Item>
        <S.FactoryContainer>
          <S.TopContainer>
            <S.Description>
              Tracetest needs configuration information to be able to retrieve your trace from your distributed tracing
              solution. Select your tracing data store and enter the configuration info.
            </S.Description>
            {explanation ? (
              <S.Explanation>{explanation}</S.Explanation>
            ) : (
              <S.Title>Provide the connection info for {SupportedDataStoresToName[dataStoreType!]}</S.Title>
            )}
            <DataStoreDocsBanner dataStoreType={dataStoreType!} />
            {dataStoreType && <DataStoreComponentFactory dataStoreType={dataStoreType} />}
          </S.TopContainer>
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
              <Button loading={isTestConnectionLoading} type="primary" ghost onClick={onTestConnection}>
                Test Connection
              </Button>
              <Button disabled={!isFormValid} loading={isLoading} type="primary" onClick={() => form.submit()}>
                Save
              </Button>
            </S.SaveContainer>
          </S.ButtonsContainer>
        </S.FactoryContainer>
      </S.FormContainer>
    </Form>
  );
};

export default DataStoreForm;
