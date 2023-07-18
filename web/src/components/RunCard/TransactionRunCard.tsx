import {Tooltip} from 'antd';
import {Link} from 'react-router-dom';
import TestState from 'components/TestState';
import TransactionRunActionsMenu from 'components/TransactionRunActionsMenu';
import {TestState as TestStateEnum} from 'constants/TestRun.constants';
import TransactionRun from 'models/TransactionRun.model';
import Date from 'utils/Date';
import * as S from './RunCard.styled';
import {TransactionRunStatusIcon} from '../RunStatusIcon';

interface IProps {
  linkTo: string;
  run: TransactionRun;
  transactionId: string;
}

const TransactionRunCard = ({
  run: {id: runId, createdAt, state, metadata, version, pass, fail, allStepsRequiredGatesPassed},
  transactionId,
  linkTo,
}: IProps) => {
  const metadataName = metadata?.name;
  const metadataBuildNumber = metadata?.buildNumber;
  const metadataBranch = metadata?.branch;
  const metadataUrl = metadata?.url;

  return (
    <Link to={linkTo}>
      <S.Container $isWhite>
        <TransactionRunStatusIcon state={state} hasFailedTests={!allStepsRequiredGatesPassed} />
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
          <TransactionRunActionsMenu runId={runId} transactionId={transactionId} />
        </div>
      </S.Container>
    </Link>
  );
};

export default TransactionRunCard;
