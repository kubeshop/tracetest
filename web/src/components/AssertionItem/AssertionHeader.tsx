import * as S from './AssertionItem.styled';

interface IProps {
  affectedSpans: number;
  failedChecks: number;
  passedChecks: number;
  title: string;
}

const AssertionHeader = ({affectedSpans, failedChecks, passedChecks, title}: IProps) => (
  <S.Column>
    <div>
      <S.HeaderTitle level={3}>{title}</S.HeaderTitle>
    </div>
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
  </S.Column>
);

export default AssertionHeader;
