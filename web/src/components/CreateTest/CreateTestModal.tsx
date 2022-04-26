import {useRef} from 'react';
import {useTour} from '@reactour/tour';
import {Modal, Form, Input, Button, Select, Checkbox} from 'antd';
import {DeleteOutlined} from '@ant-design/icons';
import {useCreateTestMutation, useRunTestMutation} from 'redux/services/TestService';
import {HTTP_METHOD} from 'types';
import './CreateTest.css';
import {useNavigate} from 'react-router-dom';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTourService';
import {Steps} from '../GuidedTour/homeStepList';
import CreateTestAnalyticsService from '../../services/analytics/CreateTestAnalyticsService';

interface IProps {
  visible: boolean;
  onClose: () => void;
}

const defaultHeaders = [{key: 'Content-Type', value: 'application/json', checked: true}];
const {onCreateTestFormSubmit} = CreateTestAnalyticsService;

const CreateTestModal = ({visible, onClose}: IProps): JSX.Element => {
  const navigate = useNavigate();
  const {setIsOpen} = useTour();
  const [createTest, {isLoading: isLoadingCreateTest}] = useCreateTestMutation();
  const [runTest, {isLoading: isLoadingRunTest}] = useRunTestMutation();

  const [form] = Form.useForm();
  const touchedHttpHeadersRef = useRef<{[key: string]: Boolean}>({});

  const onFinish = async (values: any) => {
    const headers = values.headersList
      .filter((i: {checked: boolean}) => i.checked)
      .map(({key, value}: {key: string; value: string}) => ({key, value}));
    const newTest = await createTest({
      name: values.name,
      serviceUnderTest: {
        request: {url: values.url, method: values.method, body: values.body, headers},
      },
    }).unwrap();
    const newTestRunResult = await runTest(newTest.testId).unwrap();
    onClose();
    setIsOpen(false);
    navigate(`/test/${newTest.testId}?resultId=${newTestRunResult.resultId}`);
  };

  const onFinishFailed = () => {};

  const renderActionButtons = () => {
    return (
      <>
        <Button type="ghost" htmlType="button" onClick={onClose}>
          Cancel
        </Button>

        <Button
          type="primary"
          form="newTest"
          onClick={() => {
            onCreateTestFormSubmit();
            form.submit();
          }}
          loading={isLoadingCreateTest || isLoadingRunTest}
          data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Run)}
        >
          Run Test
        </Button>
      </>
    );
  };

  return (
    <Modal
      closable
      title="Create New Test"
      visible={visible}
      footer={renderActionButtons()}
      onCancel={onClose}
      wrapClassName="test-modal"
    >
      <div style={{display: 'flex'}}>
        <Form
          name="newTest"
          form={form}
          initialValues={{remember: true}}
          onFinish={onFinish}
          onFinishFailed={onFinishFailed}
          autoComplete="off"
          layout="vertical"
        >
          <div style={{display: 'flex', marginBottom: 24}}>
            <Form.Item name="method" initialValue="GET" valuePropName="value" noStyle>
              <Select
                style={{minWidth: 120}}
                className="method-select"
                dropdownClassName="method-select-item"
                data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Method)}
              >
                {Object.keys(HTTP_METHOD).map(el => {
                  return (
                    <Select.Option key={el} value={el}>
                      {el}
                    </Select.Option>
                  );
                })}
              </Select>
            </Form.Item>

            <Form.Item name="url" rules={[{required: true, message: 'Please input Endpoint!'}]} noStyle>
              <Input
                placeholder="Enter request url"
                data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Url)}
              />
            </Form.Item>
          </div>

          <Form.Item
            name="name"
            label="Name"
            colon={false}
            data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Name)}
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
          >
            <Input.TextArea style={{maxHeight: 150, height: 120}} />
          </Form.Item>
        </Form>
      </div>
    </Modal>
  );
};

export default CreateTestModal;
