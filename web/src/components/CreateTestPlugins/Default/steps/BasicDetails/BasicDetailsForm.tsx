import {Form, Input} from 'antd';
import {noop} from 'lodash';
import {StepsID} from 'components/GuidedTour/testRunSteps';
import {TDraftTest} from 'types/Test.types';
import BasicDetailsDemoHelper from './BasicDetailsDemoHelper';
import * as S from './BasicDetails.styled';

export const FORM_ID = 'create-test';

interface IProps {
  onSelectDemo?(demo: TDraftTest): void;
  demoList?: TDraftTest[];
  selectedDemo?: TDraftTest;
  isEditing?: boolean;
  isDemoEnabled?: boolean;
}

const BasicDetailsForm = ({
  onSelectDemo = noop,
  selectedDemo,
  isEditing = false,
  demoList = [],
  isDemoEnabled = false,
}: IProps) => {
  return (
    <>
      {!isEditing && Boolean(demoList.length) && isDemoEnabled && (
        <BasicDetailsDemoHelper selectedDemo={selectedDemo} onSelectDemo={onSelectDemo} demoList={demoList} />
      )}
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
    </>
  );
};

export default BasicDetailsForm;
