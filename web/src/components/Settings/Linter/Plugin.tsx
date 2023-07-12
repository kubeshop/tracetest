import {Form} from 'antd';
import Rule from './Rule';

interface IProps {
  fieldKey: number;
  baseName: string[];
}

const Plugin = ({fieldKey, baseName}: IProps) => {
  const isEnabled = Form.useWatch<boolean>([...baseName, 'enabled']) ?? true;

  return (
    <Form.List name={[fieldKey, 'rules']}>
      {fields =>
        fields.map(field => (
          <Rule baseName={[...baseName, 'rules', `${field.name}`]} key={field.key} fieldKey={field.name} isDisabled={!isEnabled} />
        ))
      }
    </Form.List>
  );
};

export default Plugin;
