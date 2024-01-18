import {Form, Input} from 'antd';
import {LeftOutlined} from '@ant-design/icons';
import {useCallback, useEffect, useMemo} from 'react';
import {SupportedDataStores, TDraftDataStore} from 'types/DataStore.types';
import DataStoreConfiguration from 'components/DataStoreConfiguration/DataStoreConfiguration';
import {useDataStore} from 'providers/DataStore/DataStore.provider';
import DataStore from 'models/DataStore.model';
import DataStoreService from 'services/DataStore.service';
import * as S from './TracingBackend.styled';

interface IProps {
  dataStore: DataStore;
  onBack(): void;
}

const Configuration = ({dataStore, onBack}: IProps) => {
  const {isLoading, isFormValid, onIsFormValid, onSaveConfig, onTestConnection} =
    useDataStore();
  const [form] = Form.useForm<TDraftDataStore>();

  const handleTestConnection = useCallback(async () => {
    const draft = form.getFieldsValue();
    onTestConnection(draft, dataStore);
  }, [form, onTestConnection, dataStore]);

  const initialValues = useMemo(
    () => DataStoreService.getInitialValues(dataStore, dataStore.type as SupportedDataStores),
    [dataStore]
  );

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

  useEffect(() => {
    onValidation({}, initialValues);
  }, [initialValues, onIsFormValid, onValidation]);

  const handleOnSubmit = useCallback(
    async (values: TDraftDataStore) => {
      onSaveConfig(values, dataStore);
    },
    [onSaveConfig, dataStore]
  );

  return (
    <Form<TDraftDataStore>
      autoComplete="off"
      form={form}
      layout="vertical"
      name="wizard-datastore-config"
      onFinish={handleOnSubmit}
      onValuesChange={onValidation}
      initialValues={initialValues}
    >
      <S.Header>
        <S.NoPaddingButton ghost type="link" icon={<LeftOutlined />} onClick={onBack}>
          Back
        </S.NoPaddingButton>
      </S.Header>
      <Form.Item name="dataStoreType" hidden>
        <Input type="hidden" />
      </Form.Item>
      <DataStoreConfiguration
        onSubmit={() => form.submit()}
        onTestConnection={handleTestConnection}
        dataStoreType={dataStore.type as SupportedDataStores}
        isSubmitLoading={isLoading}
        isValid={isFormValid}
        withColor
      />
    </Form>
  );
};

export default Configuration;
