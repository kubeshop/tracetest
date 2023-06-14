import * as S from './AnalyzerScore.styled';

interface IProps {
  score: number;
  height?: string;
  width?: string;
}

const AnalyzerScore = ({score, height, width}: IProps) => (
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

export default AnalyzerScore;
