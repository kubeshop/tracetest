import {noop} from 'lodash';
import {Tooltip} from 'antd';
import {EditorView} from '@codemirror/view';
import {Extension} from '@codemirror/state';
import {autocompletion} from '@codemirror/autocomplete';
import CodeMirror, {ReactCodeMirrorRef} from '@uiw/react-codemirror';
import {useCallback, useMemo, useRef} from 'react';
import {useEnvironment} from 'providers/Environment/Environment.provider';
import EditorService from 'services/Editor.service';
import {SupportedEditors} from 'constants/Editor.constants';

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
  context = {},
}: IEditorProps) => {
  const {testId = '', runId = ''} = context;
  const {selectedEnvironment} = useEnvironment();
  const editorTheme = useEditorTheme();
  const completionFn = useAutoComplete({testId, runId, onSelect: onSelectAutocompleteOption});
  const {onHover, expression} = useTooltip({environmentId: selectedEnvironment?.id, ...context});

  const ref = useRef<ReactCodeMirrorRef>(null);

  const extensionList: Extension[] = useMemo(
    () => [autocompletion({override: [completionFn]}), expressionQL(), EditorView.lineWrapping, ...extensions],
    [completionFn, extensions]
  );

  const handleHover = useCallback(() => {
    if (EditorService.getIsQueryValid(SupportedEditors.Expression, value)) onHover(value);
  }, [onHover, value]);

  return (
    <S.ExpressionEditorContainer $isEditable={editable}>
      <Tooltip placement="topLeft" title={expression}>
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
          onMouseOver={handleHover}
        />
      </Tooltip>
    </S.ExpressionEditorContainer>
  );
};

export default Expression;
