import {noop} from 'lodash';
import {EditorView, hoverTooltip} from '@codemirror/view';
import {Extension} from '@codemirror/state';
import {autocompletion} from '@codemirror/autocomplete';
import CodeMirror, {ReactCodeMirrorRef} from '@uiw/react-codemirror';
import {useMemo, useRef} from 'react';

import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {expressionQL} from './grammar';
import useEditorTheme from '../hooks/useEditorTheme';
import {IEditorProps} from '../Editor';
import * as S from '../Editor.styled';
import useTooltip from '../hooks/useTooltip';
import useAutoComplete from './hooks/useAutoComplete';

const Expression = ({
  basicSetup: {lineNumbers = false, ...basicSetup} = {},
  onChange,
  placeholder,
  value = '',
  editable = true,
  extensions = [],
  autoFocus = false,
  onFocus = noop,
  indentWithTab = false,
  onSelectAutocompleteOption = noop,
}: IEditorProps) => {
  const {
    test: {id: testId},
  } = useTest();
  const {
    run: {id: runId},
  } = useTestRun();
  const editorTheme = useEditorTheme();
  const completionFn = useAutoComplete({testId, runId, onSelect: onSelectAutocompleteOption});
  const tooltipFn = useTooltip({testId, runId});
  const ref = useRef<ReactCodeMirrorRef>(null);

  const extensionList: Extension[] = useMemo(
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
        ref={ref}
        onFocus={() => onFocus(ref.current?.view!)}
        id="expression-editor"
        basicSetup={{...basicSetup, lineNumbers}}
        data-cy="expression-editor"
        value={value}
        maxHeight="120px"
        extensions={extensionList}
        onChange={onChange}
        spellCheck={false}
        theme={editorTheme}
        editable={editable}
        autoFocus={autoFocus}
        placeholder={placeholder}
        indentWithTab={indentWithTab}
      />
    </S.ExpressionEditorContainer>
  );
};

export default Expression;
