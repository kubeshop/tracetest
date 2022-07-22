import {Form, Input, Select} from 'antd';
import {HTTP_METHOD} from 'constants/Common.constants';
import Validator from "utils/Validator";
import * as S from './RequestDetails.styled';

interface IProps {
  showMethodSelector?: boolean;
  shouldValidateUrl?: boolean;
}

const RequestDetailsUrlInput = ({showMethodSelector = true, shouldValidateUrl = true}: IProps) => {
  return (
    <div>
      <S.URLInputContainer>
        {showMethodSelector && (
          <div>
            <Form.Item name="method" initialValue={HTTP_METHOD.GET} valuePropName="value" label="URL">
              <Select
                className="select-method"
                data-cy="method-select"
                dropdownClassName="select-dropdown-method"
                style={{minWidth: 120}}
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

        <Form.Item
          name="url"
          rules={[
            {
              validator: async (_, value: string) => {
                if (!shouldValidateUrl) {
                  return Promise.resolve(true);
                }
                if (value === '') {
                  return Promise.reject(new Error('Please enter a request url'));
                }
                const isValid = Validator.url(value);
                if (isValid) {
                  return Promise.resolve(isValid);
                }
                return Promise.reject(new Error('Request url is not valid'));
              },
            },
          ]}
          style={{flex: 1}}
          label={!showMethodSelector ? 'URL' : ''}
        >
          <Input data-cy="url" placeholder="Enter request url" />
        </Form.Item>
      </S.URLInputContainer>
    </div>
  );
};

export default RequestDetailsUrlInput;
