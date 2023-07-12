import {singularOrPlural} from 'utils/Common';
import Span from 'models/Span.model';
import * as S from './TestResults.styled';
import AddTestSpecButton from './AddTestSpecButton';

interface IProps {
  selectedSpan: Span;
  totalFailedSpecs: number;
  totalPassedSpecs: number;
}

const Header = ({selectedSpan, totalFailedSpecs, totalPassedSpecs}: IProps) => (
  <S.HeaderContainer>
    <S.Row>
      <div>
        {Boolean(totalPassedSpecs) && (
          <S.HeaderDetail>
            <S.HeaderDot $passed />
            {`${totalPassedSpecs} ${singularOrPlural('spec', totalPassedSpecs)} passed`}
          </S.HeaderDetail>
        )}

        {Boolean(totalFailedSpecs) && (
          <S.HeaderDetail>
            <S.HeaderDot $passed={false} />
            {`${totalFailedSpecs} ${singularOrPlural('spec', totalFailedSpecs)} failed`}
          </S.HeaderDetail>
        )}
      </div>
    </S.Row>

    <AddTestSpecButton selectedSpan={selectedSpan} />
  </S.HeaderContainer>
);

export default Header;
