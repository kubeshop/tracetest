import {Tag} from 'antd';
import {TestState} from 'constants/TestRun.constants';
import {TTest} from 'types/Test.types';
import {TTestRun, TTestRunState} from 'types/TestRun.types';
import * as S from './TransactionRunResult.styled';

const iconBasedOnResult = (result: TTestRunState, index: number) => {
  switch (result) {
    case TestState.FINISHED:
      return <S.IconSuccess />;
    case TestState.FAILED:
      return <S.IconSuccess />;
    default:
      return index + 1;
  }
};

interface IProps {
  index: number;
  test: TTest;
  testRun: TTestRun;
}

const ExecutionStep = ({index, test: {name, trigger}, testRun: {id, state, testVersion}}: IProps) => {
  return (
    <S.Container data-cy={`run-card-${name}`} key={id}>
      <div>{iconBasedOnResult(state, index)}</div>
      <S.Info>
        <S.Title>{`${name} v${testVersion}`}</S.Title>
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
