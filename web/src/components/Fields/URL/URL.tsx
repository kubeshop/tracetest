import {Col, Form, Row, Select} from 'antd';
import {HTTP_METHOD} from 'constants/Common.constants';
import {SupportedEditors} from 'constants/Editor.constants';
import {Editor, DockerTip} from 'components/Inputs';
import {useVariableSet} from 'providers/VariableSet';
import useUrlValidation from 'hooks/useValidateUrl';

interface IProps {
  showMethodSelector?: boolean;
}

const URL = ({showMethodSelector = true}: IProps) => {
  const {selectedVariableSet} = useVariableSet();
  const {onValidate} = useUrlValidation({
    variableSetId: selectedVariableSet?.id,
  });

  return (
    <>
      <Row>
        {showMethodSelector && (
          <Col span={3}>
            <Form.Item name="method" initialValue={HTTP_METHOD.GET} valuePropName="value">
              <Select
                data-cy="method-select"
                filterOption={(input, option) => option?.key?.toLowerCase().includes(input.toLowerCase())}
                showSearch
              >
                {Object.keys(HTTP_METHOD).map(method => {
                  return (
                    <Select.Option data-cy={`method-select-option-${method}`} key={method} value={method}>
                      {method}
                    </Select.Option>
                  );
                })}
              </Select>
            </Form.Item>
          </Col>
        )}

        <Col span={showMethodSelector ? 21 : 24}>
          <Form.Item
            data-cy="url"
            name="url"
            rules={[
              {required: true, message: 'Please enter a URL'},
              {
                message: 'Please enter a valid URL',
                validator: async (_, value) => {
                  const isValid = await onValidate(value);
                  if (!isValid) return Promise.reject(new Error('Please enter a valid URL'));

                  return Promise.resolve();
                },
                validateTrigger: 'onSubmit',
              },
            ]}
          >
            <Editor type={SupportedEditors.Interpolation} placeholder="Enter URL" />
          </Form.Item>
        </Col>
      </Row>

      <Form.Item name="url" noStyle>
        <DockerTip />
      </Form.Item>
    </>
  );
};

export default URL;
