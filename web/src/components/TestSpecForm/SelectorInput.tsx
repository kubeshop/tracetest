import {Form, FormInstance} from 'antd';
import Editor from 'components/Editor';
import {IValues} from './TestSpecForm';
import useQuerySelector from './hooks/useQuerySelector';
import {SupportedEditors} from '../../constants/Editor.constants';

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
      <Editor
        type={SupportedEditors.Selector}
        basicSetup={{lineNumbers: true}}
        placeholder="Leaving it empty will select All Spans"
      />
    </Form.Item>
  );
};

export default SelectorInput;
