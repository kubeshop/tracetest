import {TTransactionRun} from 'types/TransactionRun.types';
import ExecutionStep from './ExecutionStep';
import * as S from './TransactionRunResult.styled';
import Variable from './Variable';

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
          <Variable key={value.key} value={value} />
        ))}
      </div>
    </S.ResultContainer>
  );
};

export default TransactionRunResult;
