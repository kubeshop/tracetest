import {Form, Switch} from 'antd';
import {LinterPlugin} from 'models/Linter.model';
import * as S from '../common/Settings.styled';

interface IProps {
  formId: string;
  index: number;
  plugin: LinterPlugin;
}

const Plugin = ({formId, index, plugin}: IProps) => {
  return (
    <>
      <Form.Item hidden name={['plugins', index, 'name']} />
      <S.SwitchContainer>
        <label htmlFor={`${formId}_plugins_${index}_enabled`}>Enable {plugin.name}</label>
        <Form.Item name={['plugins', index, 'enabled']} valuePropName="checked">
          <Switch />
        </Form.Item>
      </S.SwitchContainer>
    </>
  );
};

export default Plugin;
