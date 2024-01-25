import {TestStateMap} from 'constants/TestRun.constants';
import {TTestRunState} from 'types/TestRun.types';
import TestStateBadge from './TestStateBadge';
import TestStateProgress from './TestStateProgress';

interface IProps {
  testState: TTestRunState;
  info?: string;
}

const TestState = ({testState, info}: IProps) => {
  const {label, percent, status, showInfo} = TestStateMap[testState];

  return percent ? (
    <TestStateProgress label={label} percent={percent} showInfo={showInfo} info={info} />
  ) : (
    <TestStateBadge label={label} status={status} />
  );
};

export default TestState;
