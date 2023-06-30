import * as S from '../Timeline.styled';

interface IProps {
  totalFailedChecks?: number;
  totalPassedChecks?: number;
}

const Header = ({totalFailedChecks, totalPassedChecks}: IProps) => {
  const failedChecksX = totalPassedChecks ? 20 : 0;

  return (
    <>
      {!!totalPassedChecks && (
        <>
          <S.CircleCheck cx={0} cy={0} r={4} $passed />
          <S.TextDescription dominantBaseline="middle" x={6} y={0}>
            {totalPassedChecks}
          </S.TextDescription>
        </>
      )}
      {!!totalFailedChecks && (
        <>
          <S.CircleCheck cx={failedChecksX} cy={0} r={4} $passed={false} />
          <S.TextDescription dominantBaseline="middle" x={failedChecksX + 6} y={0}>
            {totalFailedChecks}
          </S.TextDescription>
        </>
      )}
    </>
  );
};

export default Header;
