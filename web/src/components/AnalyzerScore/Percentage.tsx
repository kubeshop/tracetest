import * as S from './AnalyzerScore.styled';

interface IProps {
  score: number;
  height?: string;
  width?: string;
}

const Percentage = ({score, height, width}: IProps) => (
  <S.PercentageScoreWrapper>
    <S.Score level={1} $fontSize={24}>
      {score}%
    </S.Score>
    <S.ScoreProgress $height={height} $width={width} $score={score} format={() => ''} percent={score} type="circle" />
  </S.PercentageScoreWrapper>
);

export default Percentage;
