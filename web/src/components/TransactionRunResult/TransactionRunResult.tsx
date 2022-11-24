import KeyValueRow from 'components/KeyValueRow';
import {TTransaction} from 'types/Transaction.types';
import {TTransactionRun} from 'types/TransactionRun.types';
import {TestState} from 'constants/TestRun.constants';
import ExecutionStep from './ExecutionStep';
import * as S from './TransactionRunResult.styled';

interface IProps {
  transaction: TTransaction;
  transactionRun: TTransactionRun;
}

const TransactionRunResult = ({transactionRun: {steps, environment, state}, transaction}: IProps) => {
  const hasRunFailed = state === TestState.FAILED;

  return (
    <S.ResultContainer>
      <div>
        <S.Title>Execution Steps</S.Title>
        {steps.map((step, index) => (
          <ExecutionStep
            index={index}
            key={step.id}
            test={transaction.steps[index]}
            testRun={step}
            hasRunFailed={hasRunFailed}
          />
        ))}
      </div>
      <div>
        <S.Title>Variables</S.Title>
        {environment?.values?.map(value => (
          <KeyValueRow key={value.key} keyName={value.key} value={value.value} />
        ))}
      </div>
    </S.ResultContainer>
  );
};

export default TransactionRunResult;
