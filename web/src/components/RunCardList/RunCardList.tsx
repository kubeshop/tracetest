import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import {Badge, Tooltip} from 'antd';
import {QuestionCircleOutlined} from '@ant-design/icons';
import * as S from './RunCardList.styled';
import ResultCard from '../RunCard';
import {TTestRun} from '../../types/TestRun.types';

interface IResultCardListProps {
  resultList: TTestRun[];
  testId: string;
}

const ResultCardList: React.FC<IResultCardListProps> = ({resultList, testId}) => {
  const navigate = useNavigate();

  const onResultClick = useCallback(
    (resultId: string) => {
      navigate(`/test/${testId}/run/${resultId}`);
    },
    [navigate, testId]
  );

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
          <ResultCard key={run.id} run={run} onClick={onResultClick} onDelete={() => console.log('onDelete')} />
        ))}
      </S.List>
    </S.ResultCardList>
  );
};

export default ResultCardList;
