import {Form, FormInstance} from 'antd';

import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';
import useQuerySelector from './hooks/useQuerySelector';
import {IValues} from './TestSpecForm';

interface IProps {
  form: FormInstance<IValues>;
  testId: string;
  runId: number;
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
    <Form.Item name="selector" validateTrigger={[]} style={{marginBottom: '8px'}}>
      <Editor
        type={SupportedEditors.Selector}
        basicSetup={{lineNumbers: true}}
        placeholder="Leaving it empty will select All Spans"
      />
    </Form.Item>
  );
};

export default SelectorInput;
