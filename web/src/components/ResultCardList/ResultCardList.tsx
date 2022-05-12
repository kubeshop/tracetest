import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import {Badge, Tooltip} from 'antd';
import {QuestionCircleOutlined} from '@ant-design/icons';
import * as S from './ResultCardList.styled';
import ResultCard from '../ResultCard';
import {ITestRunResult} from '../../types/TestRunResult.types';

interface IResultCardListProps {
  resultList: ITestRunResult[];
}

const ResultCardList: React.FC<IResultCardListProps> = ({resultList}) => {
  const navigate = useNavigate();

  const onResultClick = useCallback(
    (testId: string, resultId: string) => {
      navigate(`/test/${testId}/result/${resultId}`);
    },
    [navigate]
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
          <Tooltip title="The number of Total/Pass/Fail assertions">
            <QuestionCircleOutlined style={{color: '#8C8C8C', cursor: 'pointer'}} />
          </Tooltip>
        </S.FailedContainer>
        <S.TextContainer style={{textAlign: 'right'}}>
          <S.Text>Actions</S.Text>
        </S.TextContainer>
      </S.Header>
      <S.List>
        {resultList.map(result => (
          <ResultCard
            key={result.resultId}
            result={result}
            onClick={onResultClick}
            onDelete={() => console.log('onDelete')}
          />
        ))}
      </S.List>
    </S.ResultCardList>
  );
};

export default ResultCardList;
