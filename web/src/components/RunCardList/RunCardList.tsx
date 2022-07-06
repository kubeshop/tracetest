import {Badge} from 'antd';
import {useNavigate} from 'react-router-dom';
import {useTheme} from 'styled-components';

import ResultCard from 'components/RunCard';
import {TooltipQuestion} from 'components/TooltipQuestion/TooltipQuestion';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {TTestRun} from 'types/TestRun.types';
import * as S from './RunCardList.styled';

interface IProps {
  resultList: TTestRun[];
  testId: string;
}

const ResultCardList = ({resultList, testId}: IProps) => {
  const theme = useTheme();
  const navigate = useNavigate();

  const handleOnResultClick = (runId: string) => {
    TestAnalyticsService.onTestRunClick();
    navigate(`/test/${testId}/run/${runId}`);
  };

  return (
    <S.ResultCardList data-cy="result-card-list">
      <S.Header>
        <S.TextContainer>
          <S.Title>Time</S.Title>
        </S.TextContainer>
        <S.TextContainer>
          <S.Title>Execution time</S.Title>
        </S.TextContainer>
        <S.TextContainer>
          <S.Title>Version</S.Title>
        </S.TextContainer>
        <S.TextContainer>
          <S.Title>Status</S.Title>
        </S.TextContainer>
        <S.TextContainer>
          <S.Title>Total</S.Title>
        </S.TextContainer>
        <S.TextContainer>
          <Badge count="P" style={{backgroundColor: theme.color.success}} />
        </S.TextContainer>
        <S.FailedContainer>
          <Badge count="F" style={{backgroundColor: theme.color.error}} />
          <TooltipQuestion margin={0} title="The number of Total/Pass/Fail checks" />
        </S.FailedContainer>
        <S.TextContainer style={{textAlign: 'right'}}>
          <S.Title>Actions</S.Title>
        </S.TextContainer>
      </S.Header>
      <S.List>
        {resultList.map(run => (
          <ResultCard key={run.id} run={run} testId={testId} onClick={handleOnResultClick} />
        ))}
      </S.List>
    </S.ResultCardList>
  );
};

export default ResultCardList;
