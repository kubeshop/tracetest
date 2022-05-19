import {format, parseISO} from 'date-fns';
import {RESULT_DATE_FORMAT} from '../../constants/Date.constants';
import {ITestRunResult} from '../../types/TestRunResult.types';
import TestState from '../TestState';
import * as S from './ResultCard.styled';
import ResultCardActions from './ResultCardActions';

interface IResultCardProps {
  result: ITestRunResult;
  onDelete(resultId: string): void;
  onClick(testId: string, resultId: string): void;
}

const ResultCard: React.FC<IResultCardProps> = ({
  result: {
    resultId,
    testId,
    executionTime,
    totalAssertionCount,
    passedAssertionCount,
    failedAssertionCount,
    state,
    createdAt,
  },
  onClick,
  onDelete,
}) => {
  const startDate = format(parseISO(createdAt), RESULT_DATE_FORMAT);

  return (
    <S.ResultCard onClick={() => onClick(testId, resultId)} data-cy={`result-card-${resultId}`}>
      <S.TextContainer>
        <S.Text>{startDate}</S.Text>
      </S.TextContainer>
      <S.TextContainer>
        <S.Text>{executionTime}s</S.Text>
      </S.TextContainer>
      <S.TextContainer data-cy={`test-run-result-status-${resultId}`}>
        <TestState testState={state} />
      </S.TextContainer>
      <S.TextContainer>
        <S.Text>{totalAssertionCount}</S.Text>
      </S.TextContainer>
      <S.TextContainer>
        <S.Text>{passedAssertionCount}</S.Text>
      </S.TextContainer>
      <S.TextContainer>
        <S.Text>{failedAssertionCount}</S.Text>
      </S.TextContainer>
      <ResultCardActions resultId={resultId} onDelete={onDelete} />
    </S.ResultCard>
  );
};

export default ResultCard;
