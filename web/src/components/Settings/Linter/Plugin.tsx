import {Checkbox, Form, Switch} from 'antd';
import {LinterPlugin} from 'models/Linter.model';
import * as S from '../common/Settings.styled';

interface IProps {
  formId: string;
  index: number;
  isEnabled: boolean;
  plugin: LinterPlugin;
}

const Plugin = ({formId, index, isEnabled, plugin}: IProps) => {
  return (
    <>
      <Form.Item hidden name={['plugins', index, 'name']} />
      <S.SwitchContainer>
        <label htmlFor={`${formId}_plugins_${index}_enabled`}>Enable {plugin.name}</label>
        <Form.Item name={['plugins', index, 'enabled']} valuePropName="checked">
          <Switch />
        </Form.Item>
      </S.SwitchContainer>

      {isEnabled && (
        <Form.Item name={['plugins', index, 'required']} valuePropName="checked" wrapperCol={{span: 8}}>
          <Checkbox>Required</Checkbox>
        </Form.Item>
      )}
    </>
  );
};

export default Plugin;
