import {Col, Form, Row, Select} from 'antd';
import {HTTP_METHOD} from 'constants/Common.constants';
import {SupportedEditors} from 'constants/Editor.constants';
import {Editor, DockerTip} from 'components/Inputs';

interface IProps {
  showMethodSelector?: boolean;
}

const URL = ({showMethodSelector = true}: IProps) => (
  <>
    <Row>
      {showMethodSelector && (
        <Col span={3}>
          <Form.Item name="method" initialValue={HTTP_METHOD.GET} valuePropName="value" style={{marginBottom: 0}}>
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
          rules={[{required: true, message: 'Please enter a valid URL'}]}
          style={{marginBottom: 0}}
        >
          <Editor type={SupportedEditors.Interpolation} placeholder="Enter URL" />
        </Form.Item>
      </Col>
    </Row>

    <Form.Item name="url" style={{marginBottom: 0}} noStyle>
      <DockerTip />
    </Form.Item>
  </>
);

export default URL;
