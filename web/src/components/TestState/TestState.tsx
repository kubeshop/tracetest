import TestStateBadge from './TestStateBadge';
import TestStateProgress from './TestStateProgress';
import {TestStateMap} from '../../constants/TestRunResult.constants';
import {ITestRunResult} from '../../types/TestRunResult.types';

interface IProps {
  testState: ITestRunResult['state'];
}

const TestState = ({testState}: IProps) => {
  const {label, percent, status} = TestStateMap[testState];

  return percent ? (
    <TestStateProgress label={label} percent={percent} />
  ) : (
    <TestStateBadge label={label} status={status} />
  );
};

export default TestState;
