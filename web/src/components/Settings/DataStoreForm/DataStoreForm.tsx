import {Form} from 'antd';
import {useCallback, useMemo} from 'react';
import SetupConfigService from 'services/SetupConfig.service';
import {TConfig, TDraftConfig, TDraftConfigForm} from 'types/Config.types';
import DataStoreComponentFactory from '../DataStorePlugin/DataStoreComponentFactory';
import * as S from './DataStoreForm.styled';
import DataStoreSelectionInput from './DataStoreSelectionInput';

export const FORM_ID = 'data-store';

interface IProps {
  form: TDraftConfigForm;
  config: TConfig;
  onSubmit(values: TDraftConfig): Promise<void>;
  onIsFormValid(isValid: boolean): void;
}

const DataStoreForm = ({form, onSubmit, config, onIsFormValid}: IProps) => {
  const initialValues = useMemo(() => SetupConfigService.getInitialValues(config), [config]);
  const dataStoreType = Form.useWatch('dataStoreType', form);

  const onValidation = useCallback(
    async (_: any, draft: TDraftConfig) => {
      try {
        const isValid = await SetupConfigService.validateDraft(draft);
        form.validateFields();
        onIsFormValid(isValid);
      } catch (error) {
        onIsFormValid(false);
      }
    },
    [form, onIsFormValid]
  );

  return (
    <Form<TDraftConfig>
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
        <S.Title>Provide connection info</S.Title>
        <DataStoreComponentFactory dataStoreType={dataStoreType} />
      </S.FormContainer>
    </Form>
  );
};

export default DataStoreForm;
