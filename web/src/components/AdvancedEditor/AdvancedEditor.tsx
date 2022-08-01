import {useMemo} from 'react';
import {noop} from 'lodash';
import CodeMirror from '@uiw/react-codemirror';
import {autocompletion} from '@codemirror/autocomplete';
import {tracetest} from 'utils/grammar';
import {linter} from '@codemirror/lint';
import useAutoComplete from './hooks/useAutoComplete';
import useLint from './hooks/useLint';
import useEditorTheme from './hooks/useEditorTheme';
import * as S from './AdvancedEditor.styled';

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
    () => [autocompletion({override: [completionFn]}), linter(lintFn), tracetest()],
    [completionFn, lintFn]
  );

  return (
    <S.AdvancedEditor>
      <CodeMirror
        data-cy="advanced-selector"
        value={value}
        maxHeight="120px"
        extensions={extensionList}
        onChange={onChange}
        spellCheck={false}
        autoFocus
        theme={editorTheme}
        placeholder="Selecting All Spans"
      />
    </S.AdvancedEditor>
  );
};

export default AdvancedEditor;
