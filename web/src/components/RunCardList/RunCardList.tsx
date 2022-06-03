import {Badge} from 'antd';
import ResultCard from 'components/RunCard';
import {useNavigate} from 'react-router-dom';
import {useDeleteRunByIdMutation} from 'redux/apis/TraceTest.api';
import TestAnalyticsService from '../../services/Analytics/TestAnalytics.service';
import {TTestRun} from '../../types/TestRun.types';
import {TooltipQuestion} from '../TooltipQuestion/TooltipQuestion';
import * as S from './RunCardList.styled';

interface IProps {
  resultList: TTestRun[];
  testId: string;
}

const ResultCardList = ({resultList, testId}: IProps) => {
  const navigate = useNavigate();
  const [deleteRunById] = useDeleteRunByIdMutation();

  const handleOnResultClick = (runId: string) => {
    TestAnalyticsService.onTestRunClick();
    navigate(`/test/${testId}/run/${runId}`);
  };

  const handleOnDelete = (runId: string) => {
    TestAnalyticsService.onDeleteTestRun();
    deleteRunById({testId, runId});
  };

  return (
    <S.ResultCardList data-cy="result-card-list">
      <S.Header>
        <S.TextContainer>
          <S.Text>Time</S.Text>
        </S.TextContainer>
        <S.TextContainer>
          <S.Text>Execution time</S.Text>
        </S.TextContainer>
        <S.TextContainer>
          <S.Text>Version</S.Text>
        </S.TextContainer>
        <S.TextContainer>
          <S.Text>Status</S.Text>
        </S.TextContainer>
        <S.TextContainer>
          <S.Text>Total</S.Text>
        </S.TextContainer>
        <S.TextContainer>
          <Badge count="P" style={{backgroundColor: '#49AA19'}} />
        </S.TextContainer>
        <S.FailedContainer>
          <Badge count="F" />
          <TooltipQuestion margin={0} title="The number of Total/Pass/Fail assertions" />
        </S.FailedContainer>
        <S.TextContainer style={{textAlign: 'right'}}>
          <S.Text>Actions</S.Text>
        </S.TextContainer>
      </S.Header>
      <S.List>
        {resultList.map(run => (
          <ResultCard key={run.id} run={run} onClick={handleOnResultClick} onDelete={handleOnDelete} />
        ))}
      </S.List>
    </S.ResultCardList>
  );
};

export default ResultCardList;
