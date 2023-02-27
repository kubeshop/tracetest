import {Form, Select} from 'antd';
import {HTTP_METHOD} from 'constants/Common.constants';
import {SupportedEditors} from 'constants/Editor.constants';
import Editor from 'components/Editor';
import * as S from './RequestDetails.styled';

interface IProps {
  showMethodSelector?: boolean;
}

const RequestDetailsUrlInput = ({showMethodSelector = true}: IProps) => {
  return (
    <div>
      <S.Label>URL</S.Label>
      <S.URLInputContainer $showMethodSelector={showMethodSelector}>
        {showMethodSelector && (
          <div>
            <Form.Item name="method" initialValue={HTTP_METHOD.GET} valuePropName="value">
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
          </div>
        )}

        <Form.Item data-cy="url" name="url" rules={[{required: true, message: 'Please enter a request url'}]}>
          <Editor type={SupportedEditors.Interpolation} placeholder="Enter request url" />
        </Form.Item>
      </S.URLInputContainer>
    </div>
  );
};

export default RequestDetailsUrlInput;
