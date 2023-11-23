import CodeMirror from '@uiw/react-codemirror';
import {StreamLanguage} from '@codemirror/language';
import {shell} from '@codemirror/legacy-modes/mode/shell';
import {IEditorProps} from '../Editor';

const CurlCommand = ({value, onChange}: IEditorProps) => {
  return (
    <CodeMirror
      value={value}
      onChange={onChange}
      data-cy="curl-command-editor"
      basicSetup={{lineNumbers: true, indentOnInput: true}}
      extensions={[StreamLanguage.define(shell)]}
      spellCheck={false}
      placeholder="curl -X POST http://site.com"
      autoFocus
    />
  );
};

export default CurlCommand;
