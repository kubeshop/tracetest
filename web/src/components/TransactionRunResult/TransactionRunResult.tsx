import KeyValueRow from 'components/KeyValueRow';
import {TTransactionRun} from 'types/TransactionRun.types';
import ExecutionStep from './ExecutionStep';
import * as S from './TransactionRunResult.styled';

interface IProps {
  transactionRun: TTransactionRun;
}

const TransactionRunResult = ({transactionRun: {steps, stepRuns, environment}}: IProps) => {
  return (
    <S.ResultContainer>
      <div>
        <S.Title>Execution Steps</S.Title>
        {steps.map((step, index) => {
          return <ExecutionStep index={index} key={step.id} test={step} testRun={stepRuns[index]} />;
        })}
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
