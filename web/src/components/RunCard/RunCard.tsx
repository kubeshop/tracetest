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

const ResultCard = ({
  run: {id: runId, executionTime, passedAssertionCount, failedAssertionCount, state, createdAt, testVersion, metadata},
  testId,
}: IProps) => {
  const metadataName = metadata?.name;
  const metadataBuildNumber = metadata?.buildNumber;
  const metadataBranch = metadata?.branch;
  const metadataUrl = metadata?.url;

  return (
    <Link to={`/test/${testId}/run/${runId}`}>
      <S.Container data-cy={`run-card-${runId}`}>
        <div>{getIcon(state, failedAssertionCount)}</div>

        <S.Info>
          <div>
            <S.Title>v{testVersion}</S.Title>
          </div>
          <div>
            <Tooltip title={Date.format(createdAt)}>
              <S.Text>{Date.getTimeAgo(createdAt)}</S.Text>
            </Tooltip>
            <S.Text> - {executionTime}s</S.Text>

            {metadataName && (
              <a href={metadataUrl} target="_blank" onClick={event => event.stopPropagation()}>
                <S.Text $hasLink={Boolean(metadataUrl)}> - {`${metadataName} ${metadataBuildNumber}`}</S.Text>
              </a>
            )}
            {metadataBranch && <S.Text> - Branch: {metadataBranch}</S.Text>}
          </div>
        </S.Info>

        <div>
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
        </div>

        <S.TestStateContainer>
          <TestState testState={state} />
        </S.TestStateContainer>

        <div>
          <RunActionsMenu resultId={runId} testId={testId} testVersion={testVersion} />
        </div>
      </S.Container>
    </Link>
  );
};

export default ResultCard;
