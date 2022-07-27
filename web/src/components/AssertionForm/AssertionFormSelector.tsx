import {Form, FormInstance} from 'antd';
import {IValues} from './AssertionForm';
import * as S from './AssertionForm.styled';
import AdvancedEditor from '../AdvancedEditor';
import useQuerySelector from './hooks/useQuerySelector';

interface IProps {
  form: FormInstance<IValues>;
  testId: string;
  runId: string;
  onValidSelector(isValid: boolean): void;
}

const AssertionFormSelector = ({form, testId, runId, onValidSelector}: IProps) => {
  const {isValid} = useQuerySelector({
    form,
    runId,
    testId,
    onValidSelector,
  });

  return (
    <S.AdvancedSelectorInputContainer>
      <Form.Item
        name="selector"
        validateTrigger={[]}
        hasFeedback
        help={!isValid ? 'Invalid selector' : ''}
        validateStatus={!isValid ? 'error' : ''}
      >
        <AdvancedEditor runId={runId} testId={testId} />
      </Form.Item>
    </S.AdvancedSelectorInputContainer>
  );
};

export default AssertionFormSelector;
