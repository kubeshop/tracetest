import {Form} from 'antd';
import {useEffect, useMemo} from 'react';
import {TDeepLinkConfig} from 'services/DeepLink.service';
import Environment from 'models/Environment.model';
import Test from 'models/Test.model';
import * as S from './DeepLink.styled';
import SwitchControl from '../CLICommand/SwitchControl';
import Variables from './Variables';

interface IProps {
  onChange(deepLinkConfig: TDeepLinkConfig): void;
  environment?: Environment;
  test: Test;
  environmentId?: string;
}

const Controls = ({onChange, environment: {values} = Environment({}), test, environmentId}: IProps) => {
  const [form] = Form.useForm<TDeepLinkConfig>();
  const variables = Form.useWatch('variables', form);
  const useEnvironmentId = Form.useWatch('useEnvironmentId', form);

  const defaultValues = useMemo(
    () => ({
      variables: values,
      useEnvironmentId: false,
    }),
    [values]
  );

  useEffect(() => {
    onChange({
      variables: variables ?? [],
      useEnvironmentId: useEnvironmentId ?? false,
      environmentId,
      test,
    });
  }, [environmentId, test, onChange, variables, useEnvironmentId]);

  return (
    <Form<TDeepLinkConfig>
      form={form}
      autoComplete="off"
      initialValues={defaultValues}
      layout="horizontal"
      name="DEEP_LINK"
    >
      <S.ControlsContainer>
        <S.Title>Manage Execution</S.Title>
        <S.OptionsContainer>
          <Form.Item name="useEnvironmentId" noStyle>
            <SwitchControl id="useEnvironmentId" text="Use Current Environment" />
          </Form.Item>
        </S.OptionsContainer>
        <Variables />
      </S.ControlsContainer>
    </Form>
  );
};

export default Controls;
