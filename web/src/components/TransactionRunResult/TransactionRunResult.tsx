import {TTransactionRun} from 'types/Transaction.types';
import ExecutionStep from './ExecutionStep';
import * as S from './TransactionRunResult.styled';
import Variable from './Variable';

interface IProps {
  transactionRun: TTransactionRun;
}

const TransactionRunResult = ({
  transactionRun: {
    environment: {values},
    results,
  },
}: IProps) => {
  return (
    <>
      <S.Title>Execution Steps</S.Title>
      {results.map((executionStepResult, index) => {
        return <ExecutionStep executionStepResult={executionStepResult} index={index} key={executionStepResult.id} />;
      })}
      <S.Title>Variables</S.Title>
      {values.map(value => (
        <Variable key={value.key} value={value} />
      ))}
    </>
  );
};

export default TransactionRunResult;
