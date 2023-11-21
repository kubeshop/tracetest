import {Form} from 'antd';
import {useEffect, useMemo} from 'react';
import {TDeepLinkConfig} from 'services/DeepLink.service';
import VariableSet from 'models/VariableSet.model';
import Test from 'models/Test.model';
import {SwitchControl} from 'components/Inputs';
import * as S from './DeepLink.styled';
import Variables from './Variables';

interface IProps {
  onChange(deepLinkConfig: TDeepLinkConfig): void;
  variableSet?: VariableSet;
  test: Test;
  variableSetId?: string;
}

const Controls = ({onChange, variableSet: {values} = VariableSet({}), test, variableSetId}: IProps) => {
  const [form] = Form.useForm<TDeepLinkConfig>();
  const variables = Form.useWatch('variables', form);
  const useVariableSetId = Form.useWatch('useVariableSetId', form);

  const defaultValues = useMemo(
    () => ({
      variables: values,
      useVariableSetId: true,
    }),
    [values]
  );

  useEffect(() => {
    onChange({
      variables: variables ?? [],
      useVariableSetId: useVariableSetId ?? true,
      variableSetId,
      test,
    });
  }, [variableSetId, test, onChange, variables, useVariableSetId]);

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
          <Form.Item name="useVariableSetId" noStyle>
            <SwitchControl id="useVariableSetId" text="Use Current Variable Set" />
          </Form.Item>
        </S.OptionsContainer>
        <Variables />
      </S.ControlsContainer>
    </Form>
  );
};

export default Controls;
