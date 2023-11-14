import {Form, Input} from 'antd';
import {StepsID} from 'components/GuidedTour/testRunSteps';
import * as S from './BasicDetailsForm.styled';

export const FORM_ID = 'create-test';

const BasicDetailsForm = () => {
  return (
    <S.InputContainer>
      <Form.Item
        className="input-name"
        data-cy="create-test-name-input"
        data-tour={StepsID.Trigger}
        label="Name"
        name="name"
        rules={[{required: true, message: 'Please enter a test name'}]}
        style={{marginBottom: 0}}
      >
        <Input placeholder="Enter test name" />
      </Form.Item>
      <Form.Item
        className="input-description"
        data-cy="create-test-description-input"
        label="Description"
        name="description"
        style={{marginBottom: 0}}
      >
        <Input.TextArea placeholder="Enter a brief description" />
      </Form.Item>
    </S.InputContainer>
  );
};

export default BasicDetailsForm;
