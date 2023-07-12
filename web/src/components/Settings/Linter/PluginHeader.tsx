import {Form, Space, Switch, Typography} from 'antd';
import {LinterPlugin} from 'models/Linter.model';
import {TDraftLinter} from 'types/Settings.types';
import * as S from '../common/Settings.styled';

interface IProps {
  fieldKey: number;
}

const PluginHeader = ({fieldKey}: IProps) => {
  const form = Form.useFormInstance<TDraftLinter>();
  const plugin = Form.useWatch<LinterPlugin | undefined>(['plugins', `${fieldKey}`], form) ?? LinterPlugin({});

  return (
    <Space>
      <S.SwitchContainer>
        <Form.Item name={[fieldKey, 'enabled']} valuePropName="checked" noStyle>
          <Switch onClick={(_, event) => event.stopPropagation()} />
        </Form.Item>
      </S.SwitchContainer>
      <Typography.Text strong>{plugin.name}</Typography.Text>
    </Space>
  );
};

export default PluginHeader;
