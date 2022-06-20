import {Form, Input, Select} from 'antd';
import React from 'react';
import {HTTP_METHOD} from '../../constants/Common.constants';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {Steps} from '../GuidedTour/homeStepList';
import * as S from './CreateTestModal.styled';

export const CreateTestFormUrlInput: React.FC = () => {
  return (
    <S.Row>
      <div>
        <Form.Item name="method" initialValue={HTTP_METHOD.GET} valuePropName="value" noStyle>
          <Select
            className="select-method"
            data-cy="method-select"
            data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Method)}
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

      <Form.Item
        name="url"
        rules={[
          {required: true, message: 'Please enter a request url'},
          {type: 'url', message: 'Request url is not valid'},
        ]}
        style={{flex: 1}}
      >
        <Input
          data-cy="url"
          data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Url)}
          placeholder="Enter request url"
        />
      </Form.Item>
    </S.Row>
  );
};
