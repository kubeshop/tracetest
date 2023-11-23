import CodeMirror from '@uiw/react-codemirror';
import {StreamLanguage} from '@codemirror/language';
import {yaml} from '@codemirror/legacy-modes/mode/yaml';
import {IEditorProps} from '../Editor';

const placeholder = `type: Test
spec:
  name: My Test
  trigger:
    type: http
    httpRequest:
      method: GET
      url: google.com
`;

const Definition = ({value, onChange}: IEditorProps) => {
  return (
    <CodeMirror
      value={value}
      onChange={onChange}
      data-cy="definition-editor"
      basicSetup={{indentOnInput: true}}
      extensions={[StreamLanguage.define(yaml)]}
      spellCheck={false}
      placeholder={placeholder}
      minHeight="200px"
      autoFocus
    />
  );
};

export default Definition;
