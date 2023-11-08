import {Col, Form, Row, Select} from 'antd';
import {HTTP_METHOD} from 'constants/Common.constants';
import {SupportedEditors} from 'constants/Editor.constants';
import Editor from 'components/Editor';
import UrlDockerTipInput from './UrlDockerTipInput';

interface IProps {
  showMethodSelector?: boolean;
}

const RequestDetailsUrlInput = ({showMethodSelector = true}: IProps) => (
  <div>
    <Row>
      {showMethodSelector && (
        <Col span={2}>
          <Form.Item name="method" initialValue={HTTP_METHOD.GET} valuePropName="value" style={{marginBottom: 0}}>
            <Select
              showSearch
              className="select-method"
              data-cy="method-select"
              dropdownClassName="select-dropdown-method"
              filterOption={(input, option) => option?.key?.toLowerCase().includes(input.toLowerCase())}
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

      <Col span={showMethodSelector ? 22 : 24}>
        <Form.Item
          data-cy="url"
          name="url"
          rules={[{required: true, message: 'Please enter a request url'}]}
          style={{marginBottom: 0}}
        >
          <Editor type={SupportedEditors.Interpolation} placeholder="Enter request url" />
        </Form.Item>
      </Col>
    </Row>

    <Form.Item name="url" style={{marginBottom: 0}} noStyle>
      <UrlDockerTipInput />
    </Form.Item>
  </div>
);

export default RequestDetailsUrlInput;
