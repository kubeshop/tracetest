import {autocompletion} from '@codemirror/autocomplete';
import {linter} from '@codemirror/lint';
import {EditorState} from '@codemirror/state';
import {EditorView} from '@codemirror/view';
import CodeMirror from '@uiw/react-codemirror';
import {useMemo} from 'react';
import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';

import {selectorQL} from './grammar';
import useAutoComplete from './hooks/useAutoComplete';
import useEditorTheme from '../hooks/useEditorTheme';
import useLint from './hooks/useLint';
import {IEditorProps} from '../Editor';
import * as S from '../Editor.styled';

const Selector = ({
  basicSetup,
  onChange,
  placeholder = 'Leaving it empty will select All Spans',
  value,
  editable = true,
}: IEditorProps) => {
  const {
    test: {id: testId},
  } = useTest();
  const {
    run: {id: runId},
  } = useTestRun();
  const completionFn = useAutoComplete({testId, runId});
  const lintFn = useLint({testId, runId});
  const editorTheme = useEditorTheme();

  const extensionList = useMemo(
    () => [
      autocompletion({override: [completionFn]}),
      linter(lintFn),
      selectorQL(),
      EditorView.lineWrapping,
      EditorState.transactionFilter.of(tr => (tr.newDoc.lines > 1 ? [] : tr)),
    ],
    [completionFn, lintFn]
  );

  return (
    <S.SelectorEditorContainer $isEditable={editable}>
      <CodeMirror
        id="selector-editor"
        basicSetup={{...basicSetup, lineNumbers: false}}
        data-cy="selector-editor"
        value={value}
        maxHeight="120px"
        extensions={extensionList}
        onChange={onChange}
        spellCheck={false}
        theme={editorTheme}
        placeholder={placeholder}
        editable={editable}
      />
    </S.SelectorEditorContainer>
  );
};

export default Selector;
