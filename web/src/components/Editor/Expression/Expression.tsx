import {EditorView, hoverTooltip} from '@codemirror/view';
import {autocompletion} from '@codemirror/autocomplete';
import CodeMirror from '@uiw/react-codemirror';
import {useMemo} from 'react';

import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {expressionQL} from './grammar';
import useEditorTheme from '../hooks/useEditorTheme';
import {IEditorProps} from '../Editor';
import * as S from '../Editor.styled';
import useTooltip from '../hooks/useTooltip';
import useAutoComplete from './hooks/useAutoComplete';

const Expression = ({
  basicSetup,
  onChange,
  placeholder = 'Leaving it empty will select All Spans',
  value = '',
  editable = true,
  extensions = [],
}: IEditorProps) => {
  const {
    test: {id: testId},
  } = useTest();
  const {
    run: {id: runId},
  } = useTestRun();
  const editorTheme = useEditorTheme();
  const completionFn = useAutoComplete({testId, runId});
  const tooltipFn = useTooltip();

  const extensionList = useMemo(
    () => [
      autocompletion({override: [completionFn]}),
      expressionQL(),
      hoverTooltip(tooltipFn),
      EditorView.lineWrapping,
      ...extensions,
    ],
    [completionFn, extensions, tooltipFn]
  );

  return (
    <S.ExpressionEditorContainer $isEditable={editable}>
      <CodeMirror
        id="expression-editor"
        basicSetup={basicSetup}
        data-cy="expression-selector"
        value={value}
        maxHeight="120px"
        extensions={extensionList}
        onChange={onChange}
        spellCheck={false}
        autoFocus
        theme={editorTheme}
        placeholder={placeholder}
      />
    </S.ExpressionEditorContainer>
  );
};

export default Expression;
