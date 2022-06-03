import {QuestionCircleOutlined} from '@ant-design/icons';
import {Badge, Tooltip} from 'antd';
import {useNavigate} from 'react-router-dom';

import ResultCard from 'components/RunCard';
import {useDeleteRunByIdMutation} from 'redux/apis/TraceTest.api';
import {TTestRun} from 'types/TestRun.types';
import * as S from './RunCardList.styled';

interface IProps {
  resultList: TTestRun[];
  testId: string;
}

const ResultCardList = ({resultList, testId}: IProps) => {
  const navigate = useNavigate();
  const [deleteRunById] = useDeleteRunByIdMutation();

  const handleOnResultClick = (runId: string) => {
    navigate(`/test/${testId}/run/${runId}`);
  };

  const handleOnDelete = (runId: string) => {
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
          <Tooltip color="#FBFBFF" title="The number of Total/Pass/Fail assertions">
            <QuestionCircleOutlined style={{color: '#8C8C8C', cursor: 'pointer'}} />
          </Tooltip>
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
