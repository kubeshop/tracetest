import * as S from './AnalyzerScore.styled';

interface IProps {
  score: number;
  fontSize?: number;
  height?: string;
  width?: string;
}

const AnalyzerScore = ({score, fontSize, height, width}: IProps) => (
  <S.ScoreWrapper>
    <S.ScoreTexContainer>
      <S.Score level={1} $fontSize={fontSize}>
        {score}
      </S.Score>
    </S.ScoreTexContainer>
    <S.ScoreProgress $height={height} $width={width} format={() => ''} percent={score} $score={score} type="circle" />
  </S.ScoreWrapper>
);

export default AnalyzerScore;
