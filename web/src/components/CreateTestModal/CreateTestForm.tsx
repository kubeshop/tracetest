import {useCallback, useRef, useState} from 'react';
import {camelCase} from 'lodash';
import {Form, Input, Button, Select, Checkbox, Dropdown, Menu, Typography, FormInstance} from 'antd';
import {DeleteOutlined, DownOutlined} from '@ant-design/icons';
import './CreateTestModal.styled.ts';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {Steps} from '../GuidedTour/homeStepList';
import {HTTP_METHOD} from '../../constants/Common.constants';
import {DemoTestExampleList} from '../../constants/Test.constants';
import * as S from './CreateTestModal.styled';

const defaultHeaders = [{key: 'Content-Type', value: 'application/json', checked: true}];

export interface ICreateTestValues {
  name: string;
  url: string;
  method: HTTP_METHOD;
  body: string;
  headersList: {
    key: string;
    value: string;
    checked: boolean;
  }[];
}

interface ICreateTestFormProps {
  onSubmit(values: ICreateTestValues): Promise<void>;
  form: FormInstance<ICreateTestValues>;
}

const CreateTestForm: React.FC<ICreateTestFormProps> = ({onSubmit, form}) => {
  const touchedHttpHeadersRef = useRef<{[key: string]: Boolean}>({});
  const [selectedDemo, setSelectedDemo] = useState();

  const onDemoClick = useCallback(
    ({key}) => {
      setSelectedDemo(key);
      const {method, url, name, description, body} = DemoTestExampleList.find(example => example.name === key) || {};

      form.setFieldsValue({
        method,
        url,
        body,
        name: `${name} - ${description}`,
      });
    },
    [form]
  );

  const menuLayout = (
    <Menu onClick={onDemoClick}>
      {DemoTestExampleList.map(({name}) => (
        <Menu.Item key={name}>
          <span data-cy={`demo-example-${camelCase(name)}`}>{name}</span>
        </Menu.Item>
      ))}
    </Menu>
  );

  return (
    <Form
      name="newTest"
      form={form}
      initialValues={{remember: true}}
      onFinish={onSubmit}
      autoComplete="off"
      layout="vertical"
    >
      <S.DemoTextContainer>
        <Typography.Text>Try these examples in our demo env: </Typography.Text>
        <Dropdown overlay={menuLayout} placement="bottomCenter" trigger={['click']}>
          <S.DropdownText data-cy="example-button" className="ant-dropdown-link" id="create-test-example-selector">
            {selectedDemo || 'Choose Example'} <DownOutlined />
          </S.DropdownText>
        </Dropdown>
      </S.DemoTextContainer>
      <div style={{display: 'flex', marginBottom: 24}}>
        <Form.Item name="method" initialValue="GET" valuePropName="value" noStyle>
          <Select
            style={{minWidth: 120}}
            data-cy="method-select"
            className="method-select"
            dropdownClassName="method-select-item"
            data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Method)}
          >
            {Object.keys(HTTP_METHOD).map(el => {
              return (
                <Select.Option data-cy={`method-select-option-${el}`} key={el} value={el}>
                  {el}
                </Select.Option>
              );
            })}
          </Select>
        </Form.Item>

        <Form.Item name="url" rules={[{required: true, message: 'Please input Endpoint!'}]} noStyle>
          <Input
            placeholder="Enter request url"
            data-cy="url"
            data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Url)}
          />
        </Form.Item>
      </div>

      <Form.Item
        name="name"
        label="Name"
        colon={false}
        data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Name)}
        data-cy="name"
        wrapperCol={{span: 24, offset: 0}}
        rules={[{required: true, message: 'Please input test name!'}]}
      >
        <Input />
      </Form.Item>

      <Form.Item
        label="Headers List"
        wrapperCol={{span: 24, offset: 0}}
        data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Headers)}
      >
        <div style={{minHeight: 80}}>
          <Form.List name="headersList" initialValue={[...defaultHeaders, {}, {}]}>
            {(fields, {add, remove}) => (
              <>
                {fields.map((field, index) => (
                  <div key={field.name} style={{display: 'flex', alignItems: 'center'}}>
                    <Form.Item name={[field.name, 'checked']} valuePropName="checked" noStyle>
                      <Checkbox style={{marginRight: 8}} />
                    </Form.Item>

                    <Form.Item name={[field.name, 'key']} noStyle>
                      <Input
                        placeholder={`Header${index + 1}`}
                        onChangeCapture={() => {
                          if (!touchedHttpHeadersRef.current[field.name]) {
                            touchedHttpHeadersRef.current[field.name] = true;
                            const headers = form.getFieldsValue().headersList;
                            headers[field.name].checked = true;
                            form.setFieldsValue({headersList: headers});
                          }

                          if (fields.length - 1 === index) {
                            add({checked: false});
                          }
                        }}
                      />
                    </Form.Item>
                    <Form.Item noStyle name={[field.name, 'value']}>
                      <Input
                        placeholder={`Value${index + 1}`}
                        onChangeCapture={() => {
                          if (!touchedHttpHeadersRef.current[field.name]) {
                            touchedHttpHeadersRef.current[field.name] = true;
                            const headers = form.getFieldsValue().headersList;
                            headers[field.name].checked = true;
                            form.setFieldsValue({headersList: headers});
                          }

                          if (fields.length - 1 === index) {
                            add({checked: false});
                          }
                        }}
                      />
                    </Form.Item>

                    <Form.Item noStyle>
                      <Button
                        style={{marginLeft: 8}}
                        type="text"
                        icon={<DeleteOutlined style={{fontSize: 24, color: '#D9D9D9'}} />}
                        onClick={() => {
                          touchedHttpHeadersRef.current[field.name] = false;
                          remove(index);
                          if (fields.length === 1 || fields.length - 1 === index) {
                            add();
                          }
                        }}
                      />
                    </Form.Item>
                  </div>
                ))}
              </>
            )}
          </Form.List>
        </div>
      </Form.Item>
      <Form.Item
        label="Request body"
        name="body"
        colon={false}
        wrapperCol={{span: 24, offset: 0}}
        data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Body)}
        data-cy="body"
      >
        <Input.TextArea style={{maxHeight: 150, height: 120}} />
      </Form.Item>
    </Form>
  );
};

export default CreateTestForm;
