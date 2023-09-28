import {Tooltip} from 'antd';
import AnalyzerScore from 'components/AnalyzerScore';
import Link from 'components/Link';
import RunActionsMenu from 'components/RunActionsMenu';
import TestState from 'components/TestState';
import TestRun, {isRunStateFinished} from 'models/TestRun.model';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import Date from 'utils/Date';
import * as S from './RunCard.styled';
import RunStatusIcon from '../RunStatusIcon';

const TEST_RUN_TRACE_TAB = 'trace';
const TEST_RUN_TEST_TAB = 'test';

interface IProps {
  run: TestRun;
  testId: string;
  linkTo: string;
  origin?: string;
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
    testSuiteId,
    testSuiteRunId,
    linter,
    requiredGatesResult,
  },
  testId,
  linkTo,
  origin,
}: IProps) => {
  const {navigate} = useDashboard();
  const metadataName = metadata?.name;
  const metadataBuildNumber = metadata?.buildNumber;
  const metadataBranch = metadata?.branch;
  const metadataUrl = metadata?.url;

  const handleResultClick = (
    event: React.MouseEvent<HTMLDivElement, MouseEvent>,
    type: typeof TEST_RUN_TRACE_TAB | typeof TEST_RUN_TEST_TAB
  ) => {
    event.preventDefault();
    navigate(`${linkTo}/${type}`, {state: {origin}});
  };

  return (
    <Link to={linkTo} state={{origin}}>
      <S.Container $isWhite data-cy={`run-card-${runId}`}>
        <RunStatusIcon state={state} requiredGatesResult={requiredGatesResult} />
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

        {!!testSuiteId && !!testSuiteRunId && <S.Text>Part of test suite</S.Text>}

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
            testSuiteRunId={testSuiteRunId}
            testSuiteId={testSuiteId}
            origin={origin}
          />
        </div>
      </S.Container>
    </Link>
  );
};

export default TestRunCard;
