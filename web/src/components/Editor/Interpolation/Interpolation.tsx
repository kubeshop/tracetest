import CodeMirror from '@uiw/react-codemirror';
import {uniq} from 'lodash';
import {Tooltip} from 'antd';
import {autocompletion} from '@codemirror/autocomplete';
import {useCallback, useMemo} from 'react';
import {useVariableSet} from 'providers/VariableSet';
import EditorService from 'services/Editor.service';
import {SupportedEditors} from 'constants/Editor.constants';
import {interpolationQL} from './grammar';
import useEditorTheme from '../hooks/useEditorTheme';
import {IEditorProps} from '../Editor';
import useAutoComplete from './hooks/useAutoComplete';
import * as S from '../Editor.styled';
import useTooltip from '../hooks/useTooltip';

const Interpolation = ({
  basicSetup: {lineNumbers, ...basicSetup} = {},
  onChange,
  placeholder,
  value = '',
  extensions = [],
  indentWithTab = false,
}: IEditorProps) => {
  const {selectedVariableSet} = useVariableSet();
  const editorTheme = useEditorTheme();
  const completionFn = useAutoComplete();
  const {onHover, resolvedValues} = useTooltip({
    variableSetId: selectedVariableSet?.id,
  });

  const extensionList = useMemo(
    () => [autocompletion({override: [completionFn]}), interpolationQL(), ...extensions],
    [completionFn, extensions]
  );

  const handleHover = useCallback(() => {
    if (EditorService.getIsQueryValid(SupportedEditors.Interpolation, value)) onHover(value);
  }, [onHover, value]);

  return (
    <S.InterpolationEditorContainer $showLineNumbers={lineNumbers}>
      <Tooltip placement="topLeft" title={uniq(resolvedValues).join(',')}>
        <CodeMirror
          id="interpolation-editor"
          basicSetup={{...basicSetup, lineNumbers}}
          data-cy="interpolation-editor"
          value={value}
          maxHeight="120px"
          extensions={extensionList}
          onChange={onChange}
          spellCheck={false}
          theme={editorTheme}
          placeholder={placeholder}
          indentWithTab={indentWithTab}
          onMouseOver={handleHover}
        />
      </Tooltip>
    </S.InterpolationEditorContainer>
  );
};

export default Interpolation;
