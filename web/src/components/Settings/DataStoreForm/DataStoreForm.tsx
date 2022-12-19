import {Form} from 'antd';
import {useCallback, useMemo} from 'react';
import SetupConfigService from 'services/DataStore.service';
import {TDraftDataStore, TDataStoreForm, TDataStoreConfig} from 'types/Config.types';
import DataStoreDocsBanner from '../DataStoreDocsBanner/DataStoreDocsBanner';
import DataStoreComponentFactory from '../DataStorePlugin/DataStoreComponentFactory';
import * as S from './DataStoreForm.styled';
import DataStoreSelectionInput from './DataStoreSelectionInput';

export const FORM_ID = 'data-store';

interface IProps {
  form: TDataStoreForm;
  dataStoreConfig: TDataStoreConfig;
  onSubmit(values: TDraftDataStore): Promise<void>;
  onIsFormValid(isValid: boolean): void;
}

const DataStoreForm = ({form, onSubmit, dataStoreConfig, onIsFormValid}: IProps) => {
  const initialValues = useMemo(() => SetupConfigService.getInitialValues(dataStoreConfig), [dataStoreConfig]);
  const dataStoreType = Form.useWatch('dataStoreType', form);

  const onValidation = useCallback(
    async (_: any, draft: TDraftDataStore) => {
      const isValid = await SetupConfigService.validateDraft(draft);
      onIsFormValid(isValid);
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
          <DataStoreSelectionInput />
        </Form.Item>
        <DataStoreDocsBanner dataStoreType={dataStoreType!} />
        <S.Title>Provide connection info</S.Title>
        <DataStoreComponentFactory dataStoreType={dataStoreType} />
      </S.FormContainer>
    </Form>
  );
};

export default DataStoreForm;
