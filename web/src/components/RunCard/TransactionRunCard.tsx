import {Tooltip} from 'antd';
import {Link} from 'react-router-dom';

import TestState from 'components/TestState';
import TransactionRunActionsMenu from 'components/TransactionRunActionsMenu';
import {TestState as TestStateEnum} from 'constants/TestRun.constants';
import {TTestRun} from 'types/TestRun.types';
import {TTransactionRun} from 'types/TransactionRun.types';
import Date from 'utils/Date';
import * as S from './RunCard.styled';

interface IProps {
  linkTo: string;
  run: TTransactionRun;
  transactionId: string;
}

function getIcon(state: TTestRun['state']) {
  if (state !== TestStateEnum.FAILED && state !== TestStateEnum.FINISHED) {
    return null;
  }
  if (state === TestStateEnum.FAILED) {
    return <S.IconFail />;
  }
  return <S.IconSuccess />;
}

const TransactionRunCard = ({run: {id: runId, createdAt, state, metadata, version}, transactionId, linkTo}: IProps) => {
  const metadataName = metadata?.name;
  const metadataBuildNumber = metadata?.buildNumber;
  const metadataBranch = metadata?.branch;
  const metadataUrl = metadata?.url;

  return (
    <Link to={linkTo}>
      <S.Container $isWhite>
        <S.IconContainer>{getIcon(state)}</S.IconContainer>

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

        <div>
          <TransactionRunActionsMenu runId={runId} transactionId={transactionId} transactionVersion={version} />
        </div>
      </S.Container>
    </Link>
  );
};

export default TransactionRunCard;
