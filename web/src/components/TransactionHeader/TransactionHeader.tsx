import {Button} from 'antd';
import {useLocation, useNavigate} from 'react-router-dom';
import {TransactionRunStatusIcon} from 'components/RunStatusIcon';
import TestState from 'components/TestState';
import TransactionRunActionsMenu from 'components/TransactionRunActionsMenu';
import {TestState as TestStateEnum} from 'constants/TestRun.constants';
import {useTransaction} from 'providers/Transaction/Transaction.provider';
import {useTransactionRun} from 'providers/TransactionRun/TransactionRun.provider';
import * as S from './TransactionHeader.styled';

const transactionLastPathRegex = /\/transaction\/[\w-]+\/run\/[\d-]+\/([\w-]+)/;

function getLastPath(pathname: string): string {
  const match = pathname.match(transactionLastPathRegex);
  if (match === null) {
    return '';
  }

  return match.length > 1 ? match[1] : '';
}

const LINKS = [
  {id: 'trigger', label: 'Trigger'},
  {id: 'automate', label: 'Automate'},
];

const TransactionHeader = () => {
  const {transaction, onRun} = useTransaction();
  const {transactionRun} = useTransactionRun();
  const navigate = useNavigate();
  const {pathname} = useLocation();
  const {id: transactionId, name, version, description} = transaction;
  const {state, id: runId, allStepsRequiredGatesPassed} = transactionRun;
  const lastPath = getLastPath(pathname);

  return (
    <S.Container>
      <S.Section>
        <a onClick={() => navigate('/')} data-cy="transaction-header-back-button">
          <S.BackIcon />
        </a>
        <div>
          <S.Title data-cy="transaction-details-name">
            {name} (v{version})
          </S.Title>
          <S.Text>{description}</S.Text>
        </div>
      </S.Section>

      <S.LinksContainer>
        {LINKS.map(({id, label}) => (
          <S.Link
            key={id}
            to={`/transaction/${transactionId}/run/${runId}/${id}`}
            $isActive={lastPath === id || (!lastPath && id === LINKS[0].id)}
          >
            {label}
          </S.Link>
        ))}
      </S.LinksContainer>

      <S.SectionRight>
        {state && state !== TestStateEnum.FINISHED && (
          <S.StateContainer data-cy="transaction-run-result-status">
            <S.StateText>Status:</S.StateText>
            <TestState testState={state} />
          </S.StateContainer>
        )}
        {state && state === TestStateEnum.FINISHED && (
          <>
            <TransactionRunStatusIcon state={state!} hasFailedTests={!allStepsRequiredGatesPassed} />
            <Button ghost onClick={() => onRun(runId)} type="primary" data-cy="transaction-run-button">
              Run Transaction
            </Button>
          </>
        )}
        <TransactionRunActionsMenu transactionId={transactionId} runId={runId} isRunView />
      </S.SectionRight>
    </S.Container>
  );
};

export default TransactionHeader;
