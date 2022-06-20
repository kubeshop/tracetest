import {DeleteOutlined, PlusOutlined} from '@ant-design/icons';
import {Button, Form, Input} from 'antd';
import React from 'react';
import {DEFAULT_HEADERS} from '../../constants/Test.constants';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {Steps} from '../GuidedTour/homeStepList';
import * as S from './CreateTestModal.styled';

export const CreateTestFormHeadersInput: React.FC = () => (
  <Form.Item
    className="input-headers"
    data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Headers)}
    label="Headers list"
  >
    <Form.List name="headers" initialValue={DEFAULT_HEADERS}>
      {(fields, {add, remove}) => (
        <>
          {fields.map((field, index) => (
            <S.HeaderContainer key={field.name}>
              <Form.Item name={[field.name, 'key']} noStyle>
                <Input placeholder={`Header ${index + 1}`} />
              </Form.Item>

              <Form.Item name={[field.name, 'value']} noStyle>
                <Input placeholder={`Value ${index + 1}`} />
              </Form.Item>

              <Form.Item noStyle>
                <Button
                  icon={<DeleteOutlined style={{fontSize: 12, color: 'rgba(3, 24, 73, 0.5)'}} />}
                  onClick={() => remove(field.name)}
                  style={{marginLeft: 12}}
                  type="link"
                />
              </Form.Item>
            </S.HeaderContainer>
          ))}

          <Button
            data-cy="add-header"
            icon={<PlusOutlined />}
            onClick={() => add()}
            style={{fontWeight: 600, height: 'auto', padding: 0}}
            type="link"
          >
            Add Header
          </Button>
        </>
      )}
    </Form.List>
  </Form.Item>
);
