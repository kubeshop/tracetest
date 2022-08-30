import {EditorView} from '@codemirror/view';

import useEditorTheme from 'components/AdvancedEditor/hooks/useEditorTheme';
import {singularOrPlural} from 'utils/Common';
import {tracetest} from 'utils/grammar';
import * as S from './TestSpec.styled';

interface IProps {
  affectedSpans: number;
  assertionsFailed: number;
  assertionsPassed: number;
  title: string;
}

const Header = ({affectedSpans, assertionsFailed, assertionsPassed, title}: IProps) => {
  const editorTheme = useEditorTheme({isEditable: false});

  return (
    <S.Column>
      <S.HeaderTitle
        data-cy="advanced-selector"
        editable={false}
        extensions={[tracetest(), editorTheme, EditorView.lineWrapping]}
        maxHeight="120px"
        placeholder="Selecting All Spans"
        spellCheck={false}
        value={title || 'All Spans'}
      />

      <div>
        {Boolean(assertionsPassed) && (
          <S.HeaderDetail>
            <S.HeaderDot $passed />
            {assertionsPassed}
          </S.HeaderDetail>
        )}
        {Boolean(assertionsFailed) && (
          <S.HeaderDetail>
            <S.HeaderDot $passed={false} />
            {assertionsFailed}
          </S.HeaderDetail>
        )}
        <S.HeaderDetail>
          <S.HeaderSpansIcon />
          {`${affectedSpans} ${singularOrPlural('span', affectedSpans)}`}
        </S.HeaderDetail>
      </div>
    </S.Column>
  );
};

export default Header;
