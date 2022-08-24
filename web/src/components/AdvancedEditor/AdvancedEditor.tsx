import {autocompletion} from '@codemirror/autocomplete';
import {linter} from '@codemirror/lint';
import {EditorView} from '@codemirror/view';
import CodeMirror from '@uiw/react-codemirror';
import {noop} from 'lodash';
import {useMemo} from 'react';

import {tracetest} from 'utils/grammar';
import * as S from './AdvancedEditor.styled';
import useAutoComplete from './hooks/useAutoComplete';
import useEditorTheme from './hooks/useEditorTheme';
import useLint from './hooks/useLint';

interface IProps {
  testId: string;
  runId: string;
  value?: string;
  onChange?(value: string): void;
}

const AdvancedEditor = ({testId, runId, onChange = noop, value = ''}: IProps) => {
  const completionFn = useAutoComplete({testId, runId});
  const lintFn = useLint({testId, runId});
  const editorTheme = useEditorTheme();

  const extensionList = useMemo(
    () => [autocompletion({override: [completionFn]}), linter(lintFn), tracetest(), EditorView.lineWrapping],
    [completionFn, lintFn]
  );

  return (
    <S.AdvancedEditor>
      <CodeMirror
        id="advanced-editor"
        data-cy="advanced-selector"
        value={value}
        maxHeight="120px"
        extensions={extensionList}
        onChange={onChange}
        spellCheck={false}
        autoFocus
        theme={editorTheme}
        placeholder="Leaving it empty will select All Spans"
      />
    </S.AdvancedEditor>
  );
};

export default AdvancedEditor;
