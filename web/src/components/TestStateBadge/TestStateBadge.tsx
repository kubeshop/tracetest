import {Badge, BadgeProps} from 'antd';
import {TTestRunResult, TTestState} from '../../types/TestRunResult.types';

const BadgeStatusMap: Record<
  TTestState,
  {status: 'processing' | 'success' | 'error' | 'default' | 'warning' | undefined; label: string}
> = {
  CREATED: {
    status: 'default',
    label: 'Created',
  },
  EXECUTING: {
    status: 'processing',
    label: 'Running',
  },
  AWAITING_TRACE: {
    status: 'warning',
    label: 'Awaiting trace',
  },
  AWAITING_TEST_RESULTS: {
    status: 'success',
    label: 'Awaiting test results',
  },
  FINISHED: {
    status: 'success',
    label: 'Finished',
  },
  FAILED: {
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
