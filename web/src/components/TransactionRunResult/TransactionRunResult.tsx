import {TTransactionRun} from 'types/TransactionRun.types';
import ExecutionStep from './ExecutionStep';
import * as S from './TransactionRunResult.styled';
import Variable from './Variable';

interface IProps {
  transactionRun: TTransactionRun;
}

const TransactionRunResult = ({transactionRun: {steps, stepRuns, environment}}: IProps) => {
  return (
    <>
      <S.Title>Execution Steps</S.Title>
      {stepRuns.map((stepRun, index) => {
        return <ExecutionStep index={index} key={stepRun.id} test={steps[index]} testRun={stepRun} />;
      })}
      <S.Title>Variables</S.Title>
      {environment?.values?.map(value => (
        <Variable key={value.key} value={value} />
      ))}
    </>
  );
};

export default TransactionRunResult;
