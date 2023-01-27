import {Tooltip} from 'antd';
import {Link} from 'react-router-dom';

import RunActionsMenu from 'components/RunActionsMenu';
import TestState from 'components/TestState';
import {TestState as TestStateEnum} from 'constants/TestRun.constants';
import {TTestRun} from 'types/TestRun.types';
import Date from 'utils/Date';
import * as S from './RunCard.styled';

interface IProps {
  run: TTestRun;
  testId: string;
  linkTo: string;
}

function getIcon(state: TTestRun['state'], failedAssertions: number) {
  if (state !== TestStateEnum.FAILED && state !== TestStateEnum.FINISHED) {
    return null;
  }
  if (state === TestStateEnum.FAILED || failedAssertions > 0) {
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
  },
  testId,
  linkTo,
}: IProps) => {
  const metadataName = metadata?.name;
  const metadataBuildNumber = metadata?.buildNumber;
  const metadataBranch = metadata?.branch;
  const metadataUrl = metadata?.url;

  return (
    <Link to={linkTo}>
      <S.Container $isWhite data-cy={`run-card-${runId}`}>
        <S.IconContainer>{getIcon(state, failedAssertionCount)}</S.IconContainer>

        <S.Info>
          <div>
            <S.Title>v{testVersion}</S.Title>
          </div>
          <S.Row>
            <Tooltip title={Date.format(createdAt)}>
              <S.Text>{Date.getTimeAgo(createdAt)}</S.Text>
            </Tooltip>

            {(state === TestStateEnum.FAILED || state === TestStateEnum.FINISHED) && (
              <S.Text>&nbsp;• {executionTime}s</S.Text>
            )}

            {metadataName && (
              <a href={metadataUrl} target="_blank" onClick={event => event.stopPropagation()}>
                <S.Text $hasLink={Boolean(metadataUrl)}>&nbsp;• {`${metadataName} ${metadataBuildNumber}`}</S.Text>
              </a>
            )}
            {metadataBranch && <S.Text>&nbsp;• Branch: {metadataBranch}</S.Text>}
          </S.Row>
        </S.Info>

        {!!transactionId && !!transactionRunId && <S.Text>Part of transaction</S.Text>}

        {state !== TestStateEnum.FAILED && state !== TestStateEnum.FINISHED && (
          <div data-cy={`test-run-result-status-${runId}`}>
            <TestState testState={state} />
          </div>
        )}

        {(state === TestStateEnum.FAILED || state === TestStateEnum.FINISHED) && (
          <S.Row>
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
            testVersion={testVersion}
            transactionRunId={transactionRunId}
            transactionId={transactionId}
          />
        </div>
      </S.Container>
    </Link>
  );
};

export default TestRunCard;
