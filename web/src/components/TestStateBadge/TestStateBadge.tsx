import {Badge, BadgeProps} from 'antd';
import { TestState } from '../../entities/TestRunResult/TestRunResult.constants';
import { TTestRunResult } from '../../entities/TestRunResult/TestRunResult.types';

const BadgeStatusMap: Record<
  TestState,
  {status: 'processing' | 'success' | 'error' | 'default' | 'warning' | undefined; label: string}
> = {
  [TestState.CREATED]: {
    status: 'default',
    label: 'Created',
  },
  [TestState.EXECUTING]: {
    status: 'processing',
    label: 'Running',
  },
  [TestState.AWAITING_TRACE]: {
    status: 'warning',
    label: 'Awaiting trace',
  },
  [TestState.AWAITING_TEST_RESULTS]: {
    status: 'success',
    label: 'Awaiting test results',
  },
  [TestState.FINISHED]: {
    status: 'success',
    label: 'Finished',
  },
  [TestState.FAILED]: {
    status: 'error',
    label: 'Failed executing test run',
  },
};

interface IProps extends BadgeProps {
  testState: TTestRunResult['state'];
}

const TestStateBadge = ({testState, ...rest}: IProps) => {
  const {status, label} = BadgeStatusMap[testState] || BadgeStatusMap.CREATED;
  return <Badge {...rest} status={status} text={label} />;
};

export default TestStateBadge;
