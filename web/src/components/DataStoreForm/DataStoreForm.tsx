import {Form} from 'antd';
import {useMemo} from 'react';
import SetupConfigService from 'services/SetupConfig.service';
import {TConfig, TDraftConfig, TDraftConfigForm} from 'types/Config.types';
import DataStoreComponentFactory from '../DataStorePlugin/DataStoreComponentFactory';
import * as S from './DataStoreForm.styled';

export const FORM_ID = 'edit-test';

interface IProps {
  form: TDraftConfigForm;
  config: TConfig;
  onSubmit(values: TDraftConfig): Promise<void>;
  onValidation(allValues: any, values: TDraftConfig): void;
}

const DataStoreForm = ({form, onSubmit, config, onValidation}: IProps) => {
  const initialValues = useMemo(() => SetupConfigService.getInitialValues(config), [config]);
  const dataStoreType = Form.useWatch('dataStoreType', form);

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
        <DataStoreComponentFactory dataStoreType={dataStoreType} />
      </S.FormContainer>
    </Form>
  );
};

export default DataStoreForm;
