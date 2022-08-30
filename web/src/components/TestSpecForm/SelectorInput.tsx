import {Form, FormInstance} from 'antd';
import {IValues} from './TestSpecForm';
import AdvancedEditor from '../AdvancedEditor';
import useQuerySelector from './hooks/useQuerySelector';

interface IProps {
  form: FormInstance<IValues>;
  testId: string;
  runId: string;
  onValidSelector(isValid: boolean): void;
}

const SelectorInput = ({form, testId, runId, onValidSelector}: IProps) => {
  useQuerySelector({
    form,
    runId,
    testId,
    onValidSelector,
  });

  return (
    <Form.Item name="selector" validateTrigger={[]}>
      <AdvancedEditor lineNumbers runId={runId} testId={testId} />
    </Form.Item>
  );
};

export default SelectorInput;
