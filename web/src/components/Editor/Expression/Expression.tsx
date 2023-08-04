import {noop, uniq} from 'lodash';
import {Tooltip} from 'antd';
import {EditorView} from '@codemirror/view';
import {Extension} from '@codemirror/state';
import {autocompletion} from '@codemirror/autocomplete';
import CodeMirror, {ReactCodeMirrorRef} from '@uiw/react-codemirror';
import {useCallback, useMemo, useRef, useState} from 'react';
import {useVariableSet} from 'providers/VariableSet';
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
  autocompleteCustomValues = [],
  autoFocus = false,
  onFocus = noop,
  indentWithTab = false,
  onSelectAutocompleteOption = noop,
  context = {},
}: IEditorProps) => {
  const {testId = '', runId = ''} = context;
  const [isHovering, setIsHovering] = useState(false);
  const {selectedVariableSet} = useVariableSet();
  const editorTheme = useEditorTheme();
  const completionFn = useAutoComplete({testId, runId, onSelect: onSelectAutocompleteOption, autocompleteCustomValues});
  const {onHover, resolvedValues, isLoading} = useTooltip({variableSetId: selectedVariableSet?.id, ...context});
  const ref = useRef<ReactCodeMirrorRef>(null);
  const isValidQuery = useMemo(() => EditorService.getIsQueryValid(SupportedEditors.Expression, value), [value]);

  const extensionList: Extension[] = useMemo(
    () => [autocompletion({override: [completionFn]}), expressionQL(), EditorView.lineWrapping, ...extensions],
    [completionFn, extensions]
  );

  const handleHover = useCallback(() => {
    if (isValidQuery) onHover(value);
  }, [isValidQuery, onHover, value]);

  const title = (
    <>
      {uniq(resolvedValues).map((resolvedValue, index) => (
        // eslint-disable-next-line react/no-array-index-key
        <p key={`${resolvedValue}-${index}`}>{resolvedValue}</p>
      ))}
    </>
  );

  return (
    <S.ExpressionEditorContainer
      $isEditable={editable}
      onMouseEnter={() => Boolean(value) && !isLoading && isValidQuery && setIsHovering(true)}
      onMouseLeave={() => setIsHovering(false)}
    >
      <Tooltip placement="topLeft" title={title} visible={isHovering}>
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
