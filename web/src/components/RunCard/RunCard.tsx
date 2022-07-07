import {format, parseISO} from 'date-fns';
import {Link} from 'react-router-dom';
import {RESULT_DATE_FORMAT} from '../../constants/Date.constants';
import {TTestRun} from '../../types/TestRun.types';
import RunActionsMenu from '../RunActionsMenu';
import TestState from '../TestState';
import * as S from './RunCard.styled';

interface IResultCardProps {
  run: TTestRun;
  testId: string;
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
}) => {
  const startDate = format(parseISO(createdAt), RESULT_DATE_FORMAT);

  return (
    <Link to={`/test/${testId}/run/${runId}`}>
      <S.ResultCard data-cy={`result-card-${runId}`}>
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
        <RunActionsMenu resultId={runId} testId={testId} testVersion={testVersion} />
      </S.ResultCard>
    </Link>
  );
};

export default ResultCard;
