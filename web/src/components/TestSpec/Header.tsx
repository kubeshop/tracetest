import {singularOrPlural} from 'utils/Common';
import {SupportedEditors} from 'constants/Editor.constants';
import Editor from 'components/Editor';
import * as S from './TestSpec.styled';

interface IProps {
  affectedSpans: number;
  assertionsFailed: number;
  assertionsPassed: number;
  selector: string;
  title: string;
}

const Header = ({affectedSpans, assertionsFailed, assertionsPassed, selector, title}: IProps) => {
  return (
    <S.Column>
      {title ? (
        <S.Title>{title}</S.Title>
      ) : (
        <Editor
          type={SupportedEditors.Selector}
          editable={false}
          basicSetup={{lineNumbers: false}}
          placeholder="Selecting All Spans"
          value={selector}
        />
      )}

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
