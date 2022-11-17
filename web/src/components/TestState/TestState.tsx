import {TestStateMap} from 'constants/TestRun.constants';
import {TTestRunState} from 'types/TestRun.types';
import TestStateBadge from './TestStateBadge';
import TestStateProgress from './TestStateProgress';

interface IProps {
  testState: TTestRunState;
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
