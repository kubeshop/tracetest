import Percentage from 'components/AnalyzerScore/Percentage';
import {TooltipQuestion} from 'components/TooltipQuestion/TooltipQuestion';
import {ANALYZER_DOCUMENTATION_URL} from 'constants/Common.constants';
import * as S from './AnalyzerResult.styled';

interface IProps {
  score: number;
  minimumScore: number;
  allRulesPassed: boolean;
}

const GlobalResult = ({score, minimumScore, allRulesPassed}: IProps) => {
  const scorePassed = score >= minimumScore;
  const allPassed = scorePassed && allRulesPassed;

  return (
    <S.GlobalResultWrapper>
      <S.GlobalScoreWrapper>
        <S.Subtitle level={3}>
          Overall Trace Analyzer Score
          <TooltipQuestion
            title={
              <>
                Tracetest core system supports analyzer evaluation as part of the testing capabilities.{' '}
                <a href={ANALYZER_DOCUMENTATION_URL} target="_blank">
                  Learn more
                </a>{' '}
                about the Analyzer.
              </>
            }
          />
        </S.Subtitle>
        <S.GlobalScoreContainer>
          <Percentage score={score} />
        </S.GlobalScoreContainer>
      </S.GlobalScoreWrapper>
      <S.ScoreResultWrapper>
        <S.Subtitle level={3}>
          Minimum Acceptable Score: <S.ResultText $passed={scorePassed}>{minimumScore}%</S.ResultText>
        </S.Subtitle>
        <S.Subtitle level={3}>
          Errors Found: <S.ResultText $passed={allRulesPassed}>{!allRulesPassed ? 'Yes' : 'No'}</S.ResultText>
        </S.Subtitle>
        <S.Subtitle level={3}>
          Result: <S.ResultText $passed={allPassed}>{allPassed ? 'Passed' : 'Failed'}</S.ResultText>
        </S.Subtitle>
      </S.ScoreResultWrapper>
    </S.GlobalResultWrapper>
  );
};

export default GlobalResult;
