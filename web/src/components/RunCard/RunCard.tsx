import {format, parseISO} from 'date-fns';
import {RESULT_DATE_FORMAT} from '../../constants/Date.constants';
import {TTestRun} from '../../types/TestRun.types';
import RunActionsDropdown from '../RunActionsMenu';
import TestState from '../TestState';
import * as S from './RunCard.styled';

interface IResultCardProps {
  run: TTestRun;
  testId: string;
  onClick(resultId: string): void;
}

const ResultCard: React.FC<IResultCardProps> = ({
  run: {
    id: runId,
    executionTime,
    totalAssertionCount,
    passedAssertionCount,
    failedAssertionCount,
    state,
    createdAt,
    testVersion,
  },
  testId,
  onClick,
}) => {
  const startDate = format(parseISO(createdAt), RESULT_DATE_FORMAT);

  return (
    <S.ResultCard onClick={() => onClick(runId)} data-cy={`result-card-${runId}`}>
      <S.TextContainer>
        <S.Text>{startDate}</S.Text>
      </S.TextContainer>
      <S.TextContainer>
        <S.Text>{executionTime}s</S.Text>
      </S.TextContainer>
      <S.TextContainer>
        <S.Text>v{testVersion}</S.Text>
      </S.TextContainer>
      <S.TextContainer data-cy={`test-run-result-status-${runId}`}>
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
      <RunActionsDropdown resultId={runId} testId={testId} />
    </S.ResultCard>
  );
};

export default ResultCard;
