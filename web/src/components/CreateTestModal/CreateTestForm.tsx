import {DeleteOutlined, DownOutlined, PlusOutlined} from '@ant-design/icons';
import {Button, Dropdown, Form, FormInstance, Input, Menu, Select, Typography} from 'antd';

import {Steps} from 'components/GuidedTour/homeStepList';
import {HTTP_METHOD} from 'constants/Common.constants';
import {DEFAULT_HEADERS, DemoTestExampleList} from 'constants/Test.constants';
import {camelCase} from 'lodash';
import React from 'react';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import Validator from 'utils/Validator';
import {TooltipQuestion} from '../TooltipQuestion/TooltipQuestion';
import * as S from './CreateTestModal.styled';
import CreateTestAnalyticsService from '../../services/Analytics/CreateTestAnalytics.service';

export const FORM_ID = 'create-test';

export interface ICreateTestValues {
  body: string;
  headers: {
    key: string;
    value: string;
  }[];
  method: HTTP_METHOD;
  name: string;
  url: string;
}

interface IProps {
  form: FormInstance<ICreateTestValues>;
  onSelectDemo(value: string): void;
  onSubmit(values: ICreateTestValues): Promise<void>;
  onValidation(isValid: boolean): void;
  selectedDemo: string;
}

const CreateTestForm = ({form, onSelectDemo, onSubmit, onValidation, selectedDemo}: IProps) => {
  const handleOnDemoClick = ({key}: {key: string}) => {
    CreateTestAnalyticsService.onDemoTestClick();
    onSelectDemo(key);
    const {body, description, method, name, url} = DemoTestExampleList.find(demo => demo.name === key) || {};

    form.setFieldsValue({
      body,
      method,
      name: `${name} - ${description}`,
      url,
    });

    onValidation(true);
  };

  const handleOnValuesChange = (changedValues: any, allValues: ICreateTestValues) => {
    const isValid =
      Validator.required(allValues.name) && Validator.required(allValues.url) && Validator.url(allValues.url);
    onValidation(isValid);
  };

  const menuLayout = (
    <Menu
      items={DemoTestExampleList.map(({name}) => ({
        key: name,
        label: <span data-cy={`demo-example-${camelCase(name)}`}>{name}</span>,
      }))}
      onClick={handleOnDemoClick}
      style={{width: '190px'}}
    />
  );

  return (
    <Form
      autoComplete="off"
      data-cy="create-test-modal"
      form={form}
      layout="vertical"
      name={FORM_ID}
      onFinish={onSubmit}
      onValuesChange={handleOnValuesChange}
    >
      <S.GlobalStyle />

      <S.DemoContainer>
        <Typography.Text>Try these examples in our demo env: </Typography.Text>
        <Dropdown overlay={menuLayout} placement="bottomLeft" trigger={['click']}>
          <Typography.Link data-cy="example-button" onClick={e => e.preventDefault()} strong>
            {selectedDemo || 'Choose Example'} <DownOutlined />
          </Typography.Link>
        </Dropdown>
        <TooltipQuestion
          margin={8}
          title="We have a microservice based on the Pokemon API installed - pick an example to autofill this form"
        />
      </S.DemoContainer>

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

      <Form.Item
        className="input-name"
        data-cy="name"
        data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Name)}
        label="Name"
        name="name"
        rules={[{required: true, message: 'Please enter a test name'}]}
      >
        <Input placeholder="Enter test name" />
      </Form.Item>

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

      <Form.Item
        className="input-body"
        data-cy="body"
        data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Body)}
        label="Request body"
        name="body"
        style={{marginBottom: 0}}
      >
        <Input.TextArea placeholder="Enter request body text" />
      </Form.Item>
    </Form>
  );
};

export default CreateTestForm;
