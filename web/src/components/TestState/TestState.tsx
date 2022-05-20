import TestStateBadge from './TestStateBadge';
import TestStateProgress from './TestStateProgress';
import {TestStateMap} from '../../constants/TestRun.constants';
import {TTestRun} from '../../types/TestRun.types';

interface IProps {
  testState: TTestRun['state'];
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
