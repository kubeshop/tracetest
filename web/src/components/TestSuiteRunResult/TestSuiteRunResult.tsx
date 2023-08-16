import KeyValueRow from 'components/KeyValueRow';
import {TestState} from 'constants/TestRun.constants';
import TestSuite from 'models/TestSuite.model';
import TestSuiteRun from 'models/TestSuiteRun.model';
import ExecutionStep from './ExecutionStep';
import * as S from './TestSuiteRunResult.styled';

interface IProps {
  testSuite: TestSuite;
  testSuiteRun: TestSuiteRun;
}

const TestSuiteRunResult = ({testSuiteRun: {steps, variableSet, state}, testSuite}: IProps) => {
  const hasRunFailed = state === TestState.FAILED;

  return (
    <S.ResultContainer>
      <div>
        <S.Title>Execution Steps</S.Title>
        {steps.map((step, index) => (
          <ExecutionStep
            key={`${step.id}-${testSuite.fullSteps[index]?.id}`}
            test={testSuite.fullSteps[index]}
            testRun={step}
            hasRunFailed={hasRunFailed}
          />
        ))}
      </div>
      <div>
        <S.Title>Variables</S.Title>
        {variableSet?.values?.map(value => (
          <KeyValueRow key={value.key} keyName={value.key} value={value.value} />
        ))}
      </div>
    </S.ResultContainer>
  );
};

export default TestSuiteRunResult;
