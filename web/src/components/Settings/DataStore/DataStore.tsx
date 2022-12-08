import {Button, Form} from 'antd';
import {useSetupConfig} from 'providers/SetupConfig/SetupConfig.provider';
import {useCallback} from 'react';
import {TConfig, TDraftConfig} from 'types/Config.types';
import DataStoreForm from '../../DataStoreForm';
import * as S from './DataStore.styled';

interface IProps {
  config: TConfig;
}

const DataStore = ({config}: IProps) => {
  const {isLoading, isFormValid, onIsFormValid, onSaveConfig} = useSetupConfig();

  const [form] = Form.useForm<TDraftConfig>();

  const handleOnSubmit = useCallback(
    async (values: TDraftConfig) => {
      onSaveConfig(values);
    },
    [onSaveConfig]
  );

  return (
    <S.Wrapper data-cy="config-datastore-form">
      <S.FormContainer>
        <S.Description>
          Tracetest needs configuration information to be able to retrieve your trace from your distributed tracing
          solution. Select your tracing data store and enter the configuration info.
        </S.Description>
        <S.Title>Choose OpenTelemetry data store</S.Title>
        <DataStoreForm form={form} config={config} onSubmit={handleOnSubmit} onValidation={onIsFormValid} />
        <S.ButtonsContainer>
          <Button
            data-cy="config-datastore-submit"
            loading={isLoading}
            type="primary"
            ghost
            onClick={() => form.submit()}
          >
            Save
          </Button>
          <Button
            data-cy="config-datastore-submit"
            loading={isLoading}
            disabled={!isFormValid}
            type="primary"
            onClick={() => console.log('@@test connection')}
          >
            Test Connection
          </Button>
        </S.ButtonsContainer>
      </S.FormContainer>
    </S.Wrapper>
  );
};

export default DataStore;
