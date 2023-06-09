import * as S from './AnalyzerResult.styled';
import PercentageScore from '../LintScore/PercentageScore';
import {TooltipQuestion} from '../TooltipQuestion/TooltipQuestion';

interface IProps {
  score: number;
  minimumScore: number;
}

const GlobalResult = ({score, minimumScore}: IProps) => {
  const passedScore = score >= minimumScore;

  return (
    <S.GlobalResultWrapper>
      <S.GlobalScoreWrapper>
        <S.Subtitle level={3}>
          Overall Trace Analyzer Score
          <TooltipQuestion title="Tracetest core system supports analyzer evaluation as part of the testing capabilities." />
        </S.Subtitle>
        <S.GlobalScoreContainer>
          <PercentageScore score={score} />
        </S.GlobalScoreContainer>
      </S.GlobalScoreWrapper>
      {!!minimumScore && (
        <S.ScoreResultWrapper>
          <S.Subtitle level={3}>Minimum Acceptable Score: {minimumScore}</S.Subtitle>
          <S.Subtitle level={3}>
            Result: <S.ResultText $passed={passedScore}>{passedScore ? 'Passed' : 'Failed'}</S.ResultText>
          </S.Subtitle>
        </S.ScoreResultWrapper>
      )}
    </S.GlobalResultWrapper>
  );
};

export default GlobalResult;
