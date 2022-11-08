import {Tag} from 'antd';
import {TTransactionTestResult} from 'types/Transaction.types';
import * as S from './TransactionRunResult.styled';

const iconBasedOnResult = (result: 'success' | 'fail' | 'running', index: number) => {
  switch (result) {
    case 'success':
      return <S.IconSuccess />;
    case 'fail':
      return <S.IconSuccess />;
    default:
      return index + 1;
  }
};

interface IProps {
  executionStepResult: TTransactionTestResult;
  index: number;
}

const ExecutionStep = ({executionStepResult: {name, id, result, version, trigger}, index}: IProps) => {
  return (
    <S.Container data-cy={`run-card-${name}`} key={id}>
      <div>{iconBasedOnResult(result, index)}</div>
      <S.Info>
        <S.Title>{`${name} v${version}`}</S.Title>
        <S.TagContainer>
          {[trigger.method, trigger.type].map(d => (
            <Tag key={d}>{d}</Tag>
          ))}
        </S.TagContainer>
      </S.Info>
    </S.Container>
  );
};

export default ExecutionStep;
