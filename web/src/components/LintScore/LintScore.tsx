import * as S from './LintScore.styled';

interface IProps {
  score: number;
  height?: string;
  width?: string;
}

const LintScore = ({score, height, width}: IProps) => {
  return (
    <S.ScoreWrapper>
      <S.ScoreTexContainer>
        <S.Score level={1}>{score}</S.Score>
      </S.ScoreTexContainer>
      <S.ScoreProgress
        $height={height}
        $width={width}
        format={() => ''}
        percent={score || 100}
        $score={score}
        type="circle"
      />
    </S.ScoreWrapper>
  );
};

export default LintScore;
