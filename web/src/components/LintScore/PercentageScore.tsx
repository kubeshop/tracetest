import * as S from './LintScore.styled';

interface IProps {
  score: number;
  height?: string;
  width?: string;
}

const Percentage = ({score, height, width}: IProps) => {
  return (
    <S.PercentageScoreWrapper>
      <S.Score level={1}>{score}%</S.Score>
      <S.ScoreProgress $height={height} $width={width} $score={score} format={() => ''} percent={score || 100} type="circle" />
    </S.PercentageScoreWrapper>
  );
};

export default Percentage;
