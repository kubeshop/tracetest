import {startCompletion} from '@codemirror/autocomplete';
import {EditorView} from '@codemirror/view';
import {Form, FormInstance} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {SupportedEditors} from 'constants/Editor.constants';
import {delay} from 'lodash';
import {useCallback} from 'react';
import {TResolveExpressionContext} from 'types/Expression.types';
import {Editor} from 'components/Inputs';
import {IValues} from './TestSpecForm';

interface IProps {
  form: FormInstance<IValues>;
  name: number;
  field: Pick<FormListFieldData, never>;
  valueList: string[];
  editorContext: TResolveExpressionContext;
}

const AssertionCheckValue = ({name, valueList, field, editorContext}: IProps) => {
  const onFocus = useCallback((view: EditorView) => {
    if (!view?.state.doc.length) delay(() => startCompletion(view!), 0);
  }, []);

  return (
    <Form.Item
      {...field}
      name={[name, 'right']}
      rules={[{required: true, message: 'Expected value is required'}]}
      data-cy="assertion-check-value"
      style={{margin: 0}}
    >
      <Editor
        type={SupportedEditors.Expression}
        placeholder="Expected Value"
        context={editorContext}
        autocompleteCustomValues={valueList}
        onFocus={onFocus}
      />
    </Form.Item>
  );
};

export default AssertionCheckValue;
