import * as S from './LintScore.styled';

interface IProps {
  score: number;
  passed: boolean;
  height?: string;
  width?: string;
}

const LintScore = ({score, passed, height, width}: IProps) => {
  return (
    <S.ScoreWrapper>
      <S.ScoreTexContainer>
        <S.Score level={1}>{score}</S.Score>
      </S.ScoreTexContainer>
      <S.ScoreProgress
        $height={height}
        $width={width}
        format={() => ''}
        percent={score}
        status={passed ? 'success' : 'exception'}
        type="circle"
      />
    </S.ScoreWrapper>
  );
};

export default LintScore;
