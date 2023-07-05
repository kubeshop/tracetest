import {Tooltip} from 'antd';
import {Link, useNavigate} from 'react-router-dom';
import AnalyzerScore from 'components/AnalyzerScore';
import RunActionsMenu from 'components/RunActionsMenu';
import TestState from 'components/TestState';
import TestRun, {isRunStateFailed, isRunStateFinished, isRunStateStopped} from 'models/TestRun.model';
import Date from 'utils/Date';
import * as S from './RunCard.styled';

const TEST_RUN_TRACE_TAB = 'trace';
const TEST_RUN_TEST_TAB = 'test';

interface IProps {
  run: TestRun;
  testId: string;
  linkTo: string;
}

function getIcon(state: TestRun['state'], failedAssertions: number, isFailedAnalyzer: boolean) {
  if (!isRunStateFinished(state)) {
    return null;
  }
  if (isRunStateStopped(state)) {
    return <S.IconInfo />;
  }
  if (isRunStateFailed(state) || failedAssertions > 0 || isFailedAnalyzer) {
    return <S.IconFail />;
  }
  return <S.IconSuccess />;
}

const TestRunCard = ({
  run: {
    id: runId,
    executionTime,
    passedAssertionCount,
    failedAssertionCount,
    state,
    createdAt,
    testVersion,
    metadata,
    transactionId,
    transactionRunId,
    linter,
  },
  testId,
  linkTo,
}: IProps) => {
  const navigate = useNavigate();
  const metadataName = metadata?.name;
  const metadataBuildNumber = metadata?.buildNumber;
  const metadataBranch = metadata?.branch;
  const metadataUrl = metadata?.url;

  const handleResultClick = (
    event: React.MouseEvent<HTMLDivElement, MouseEvent>,
    type: typeof TEST_RUN_TRACE_TAB | typeof TEST_RUN_TEST_TAB
  ) => {
    event.preventDefault();
    navigate(`${linkTo}/${type}`);
  };

  return (
    <Link to={linkTo}>
      <S.Container $isWhite data-cy={`run-card-${runId}`}>
        <S.IconContainer>{getIcon(state, failedAssertionCount, linter.isFailed)}</S.IconContainer>

        <S.Info>
          <div>
            <S.Title>v{testVersion}</S.Title>
          </div>
          <S.Row>
            <Tooltip title={Date.format(createdAt)}>
              <S.Text>{Date.getTimeAgo(createdAt)}</S.Text>
            </Tooltip>

            {isRunStateFinished(state) && <S.Text>&nbsp;• {executionTime}s</S.Text>}

            {metadataName && (
              <a href={metadataUrl} target="_blank" onClick={event => event.stopPropagation()}>
                <S.Text $hasLink={Boolean(metadataUrl)}>&nbsp;• {`${metadataName} ${metadataBuildNumber}`}</S.Text>
              </a>
            )}
            {metadataBranch && <S.Text>&nbsp;• Branch: {metadataBranch}</S.Text>}
          </S.Row>
        </S.Info>

        {!!transactionId && !!transactionRunId && <S.Text>Part of transaction</S.Text>}

        {!isRunStateFinished(state) && (
          <div data-cy={`test-run-result-status-${runId}`}>
            <TestState testState={state} />
          </div>
        )}

        {isRunStateFinished(state) && !!linter.plugins.length && (
          <Tooltip title="Trace Analyzer score">
            <div onClick={event => handleResultClick(event, TEST_RUN_TRACE_TAB)}>
              <AnalyzerScore fontSize={10} width="28px" height="28px" score={linter.score} />
            </div>
          </Tooltip>
        )}

        {isRunStateFinished(state) && (
          <S.Row $minWidth={70} onClick={event => handleResultClick(event, TEST_RUN_TEST_TAB)}>
            <Tooltip title="Passed assertions">
              <S.HeaderDetail>
                <S.HeaderDot $passed />
                {passedAssertionCount}
              </S.HeaderDetail>
            </Tooltip>
            <Tooltip title="Failed assertions">
              <S.HeaderDetail>
                <S.HeaderDot $passed={false} />
                {failedAssertionCount}
              </S.HeaderDetail>
            </Tooltip>
          </S.Row>
        )}

        <div>
          <RunActionsMenu
            resultId={runId}
            testId={testId}
            transactionRunId={transactionRunId}
            transactionId={transactionId}
          />
        </div>
      </S.Container>
    </Link>
  );
};

export default TestRunCard;
