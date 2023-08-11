import {Tooltip} from 'antd';
import Link from 'components/Link';
import TestState from 'components/TestState';
import TestSuiteRunActionsMenu from 'components/TestSuiteRunActionsMenu';
import {TestState as TestStateEnum} from 'constants/TestRun.constants';
import TestSuiteRun from 'models/TestSuiteRun.model';
import Date from 'utils/Date';
import * as S from './RunCard.styled';
import {TestSuiteRunStatusIcon} from '../RunStatusIcon';

interface IProps {
  linkTo: string;
  run: TestSuiteRun;
  testSuiteId: string;
}

const TestSuiteRunCard = ({
  run: {id: runId, createdAt, state, metadata, version, pass, fail, allStepsRequiredGatesPassed},
  testSuiteId,
  linkTo,
}: IProps) => {
  const metadataName = metadata?.name;
  const metadataBuildNumber = metadata?.buildNumber;
  const metadataBranch = metadata?.branch;
  const metadataUrl = metadata?.url;

  return (
    <Link to={linkTo}>
      <S.Container $isWhite>
        <TestSuiteRunStatusIcon state={state} hasFailedTests={!allStepsRequiredGatesPassed} />
        <S.Info>
          <div>
            <S.Title>v{version}</S.Title>
          </div>
          <S.Row>
            <Tooltip title={Date.format(createdAt)}>
              <S.Text>{Date.getTimeAgo(createdAt)}</S.Text>
            </Tooltip>
            {/* Adding this latter when is available */}
            {/* <S.Text>&nbsp;• 0s (executionTime missing from API)</S.Text> */}

            {metadataName && (
              <a href={metadataUrl} target="_blank" onClick={event => event.stopPropagation()}>
                <S.Text $hasLink={Boolean(metadataUrl)}>&nbsp;• {`${metadataName} ${metadataBuildNumber}`}</S.Text>
              </a>
            )}
            {metadataBranch && <S.Text>&nbsp;• Branch: {metadataBranch}</S.Text>}
          </S.Row>
        </S.Info>

        {state !== TestStateEnum.FAILED && state !== TestStateEnum.FINISHED && (
          <div>
            <TestState testState={state} />
          </div>
        )}

        {(state === TestStateEnum.FAILED || state === TestStateEnum.FINISHED) && (
          <S.Row>
            <Tooltip title="Passed assertions">
              <S.HeaderDetail>
                <S.HeaderDot $passed />
                {pass}
              </S.HeaderDetail>
            </Tooltip>
            <Tooltip title="Failed assertions">
              <S.HeaderDetail>
                <S.HeaderDot $passed={false} />
                {fail}
              </S.HeaderDetail>
            </Tooltip>
          </S.Row>
        )}

        <div>
          <TestSuiteRunActionsMenu runId={runId} testSuiteId={testSuiteId} />
        </div>
      </S.Container>
    </Link>
  );
};

export default TestSuiteRunCard;
