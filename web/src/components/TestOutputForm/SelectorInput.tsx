import {Form, FormInstance} from 'antd';

import Editor from 'components/Editor';
import useQuerySelector from 'components/TestSpecForm/hooks/useQuerySelector';
import {SupportedEditors} from 'constants/Editor.constants';
import {noop} from 'lodash';
import {TTestOutput} from '../../types/TestOutput.types';

interface IProps {
  form: FormInstance<TTestOutput>;
  runId: string;
  spanIdList: string[];
  testId: string;
}

const SelectorInput = ({form, runId, spanIdList, testId}: IProps) => {
  const {isLoading} = useQuerySelector({
    form,
    runId,
    testId,
    onValidSelector: noop,
  });

  return (
    <Form.Item
      name="selector"
      rules={[
        {required: true, message: 'Please enter a valid selector'},
        {
          message: 'Please select a single span',
          validator: () => {
            if (spanIdList.length !== 1 && !isLoading) {
              return Promise.reject(new Error('Please select a single span'));
            }
            return Promise.resolve();
          },
          validateTrigger: 'onSubmit',
        },
      ]}
      style={{marginBottom: 0}}
    >
      <Editor basicSetup={{lineNumbers: true}} type={SupportedEditors.Selector} placeholder="" />
    </Form.Item>
  );
};

export default SelectorInput;
