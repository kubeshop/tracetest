import {tracetest} from '../../utils/grammar';
import useEditorTheme from '../AdvancedEditor/hooks/useEditorTheme';
import * as S from './AssertionItem.styled';

interface IProps {
  affectedSpans: number;
  failedChecks: number;
  passedChecks: number;
  title: string;
}

const AssertionHeader = ({affectedSpans, failedChecks, passedChecks, title}: IProps) => {
  const editorTheme = useEditorTheme({isEditable: false});

  return (
    <S.HeaderContainer>
      <S.HeaderTitleText
        editable={false}
        data-cy="advanced-selector"
        value={(title || 'All Spans').replace(/\n/g, '')}
        maxHeight="120px"
        spellCheck={false}
        extensions={[tracetest(), editorTheme]}
        placeholder="Selecting All Spans"
      />

      <div>
        {Boolean(passedChecks) && (
          <S.HeaderDetail>
            <S.HeaderDot $passed />
            {passedChecks}
          </S.HeaderDetail>
        )}
        {Boolean(failedChecks) && (
          <S.HeaderDetail>
            <S.HeaderDot $passed={false} />
            {failedChecks}
          </S.HeaderDetail>
        )}
        <S.HeaderDetail>
          <S.HeaderSpansIcon />
          {`${affectedSpans} ${affectedSpans > 1 ? 'spans' : 'span'}`}
        </S.HeaderDetail>
      </div>
    </S.HeaderContainer>
  );
};

export default AssertionHeader;
