import {Button} from 'antd';
import {Link} from 'react-router-dom';
import {useTransaction} from 'providers/Transaction/Transaction.provider';
import {TestState as TestStateEnum} from 'constants/TestRun.constants';
import TestState from '../TestState';
import * as S from './TransactionHeader.styled';
import TransactionRunActionsMenu from '../TransactionRunActionsMenu';
import {TTransaction} from '../../types/Transaction.types';
import {TTransactionRun} from '../../types/TransactionRun.types';

interface IProps {
  transaction: TTransaction;
  transactionRun: TTransactionRun;
}

const TransactionHeader = ({
  transaction: {id: transactionId, name, version, description},
  transactionRun: {state, id: runId},
}: IProps) => {
  const {onRun} = useTransaction();

  return (
    <S.Container>
      <S.Section>
        <Link to={`/transaction/${transactionId}`} data-cy="transaction-header-back-button">
          <S.BackIcon />
        </Link>
        <div>
          <S.Title data-cy="transaction-details-name">
            {name} (v{version})
          </S.Title>
          <S.Text>{description}</S.Text>
        </div>
      </S.Section>
      <S.Section>
        {state && state !== TestStateEnum.FINISHED && (
          <S.StateContainer data-cy="transaction-run-result-status">
            <S.StateText>Status:</S.StateText>
            <TestState testState={state} />
          </S.StateContainer>
        )}
        {state && state === TestStateEnum.FINISHED && (
          <Button ghost onClick={onRun} type="primary" data-cy="transaction-run-button">
            Run Transaction
          </Button>
        )}
        <TransactionRunActionsMenu transactionId={transactionId} runId={runId} isRunView transactionVersion={version} />
      </S.Section>
    </S.Container>
  );
};

export default TransactionHeader;
