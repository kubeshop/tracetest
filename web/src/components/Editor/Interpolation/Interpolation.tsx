import CodeMirror from '@uiw/react-codemirror';
// import {autocompletion} from '@codemirror/autocomplete';
// import {hoverTooltip} from '@codemirror/view';
// import {useMemo} from 'react';
// import {interpolationQL} from './grammar';
import useEditorTheme from '../hooks/useEditorTheme';
import {IEditorProps} from '../Editor';
// import useAutoComplete from './hooks/useAutoComplete';
import * as S from '../Editor.styled';
// import useTooltip from '../hooks/useTooltip';

const Interpolation = ({
  basicSetup: {lineNumbers, ...basicSetup} = {},
  onChange,
  placeholder,
  value,
  extensions = [],
}: IEditorProps) => {
  const editorTheme = useEditorTheme();
  // const completionFn = useAutoComplete();
  // const tooltipFn = useTooltip();

  // const extensionList = useMemo(
  //   () => [autocompletion({override: [completionFn]}), interpolationQL(), hoverTooltip(tooltipFn), ...extensions],
  //   [completionFn, extensions, tooltipFn]
  // );

  return (
    <S.InterpolationContainer $showLineNumbers={lineNumbers}>
      <CodeMirror
        id="interpolation-editor"
        basicSetup={{...basicSetup, lineNumbers}}
        data-cy="interpolation-editor"
        value={value}
        maxHeight="120px"
        extensions={extensions}
        onChange={onChange}
        spellCheck={false}
        theme={editorTheme}
        placeholder={placeholder}
        indentWithTab={false}
      />
    </S.InterpolationContainer>
  );
};

export default Interpolation;
