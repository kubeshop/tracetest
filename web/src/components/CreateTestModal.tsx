import {Modal as AntModal, Form, Input, Button} from 'antd';
import styled from 'styled-components';
import {useCreateTestMutation} from 'services/TestService';

const Modal = styled(AntModal)`
  .ant-modal {
    width: 100%;
    max-width: unset;
    margin: unset;
  }
  .ant-modal-centered::before {
    content: unset;
  }
`;

interface IProps {
  visible: boolean;
  onClose: () => void;
}

const CreateTestModal = ({visible, onClose}: IProps): JSX.Element => {
  const [createTest, result] = useCreateTestMutation();
  const onFinish = (values: {name: string; url: string}) => {
    createTest({
      name: values.name,
      serviceUnderTest: {
        url: values.url,
      },
    });
    onClose();
  };

  const onFinishFailed = () => {};

  return (
    <Modal title="" visible={visible} footer={null} closable onCancel={onClose}>
      <h1 style={{height: '100%', width: '100%'}}>Create New Test</h1>
      <Form
        name="newTest"
        labelCol={{span: 8}}
        wrapperCol={{span: 16}}
        initialValues={{remember: true}}
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        autoComplete="off"
      >
        <Form.Item label="Name" name="name" rules={[{required: true, message: 'Please input test name!'}]}>
          <Input />
        </Form.Item>

        <Form.Item label="Endpoint" name="url" rules={[{required: true, message: 'Please input Endpoint!'}]}>
          <Input />
        </Form.Item>

        <Form.Item wrapperCol={{offset: 8, span: 16}}>
          <Button type="primary" htmlType="submit">
            Create Test
          </Button>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CreateTestModal;
