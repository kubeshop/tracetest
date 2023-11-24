import {Form} from 'antd';
import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';
import DefinitionService from 'services/Importers/Definition.service';

const Definition = () => {
  return (
    <Form.Item
      label="Paste Tracetest Definition"
      name="definition"
      rules={[
        {required: true, message: 'Please enter a valid Tracetest Definition'},
        {
          validator: (_, definition) => {
            const errors = DefinitionService.validate(definition);
            if (errors.length) return Promise.reject(new Error(errors.join(', ')));

            return Promise.resolve();
          },
        },
      ]}
    >
      <Editor type={SupportedEditors.Definition} />
    </Form.Item>
  );
};

export default Definition;
