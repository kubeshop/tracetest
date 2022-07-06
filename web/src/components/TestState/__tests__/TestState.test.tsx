import {render} from 'test-utils';
import TestState from '../TestState';
import {TestState as TestStateEnum, TestStateMap} from '../../../constants/TestRun.constants';

describe('TestState', () => {
  it('should render badge component', () => {
    const {container, getByText} = render(<TestState testState={TestStateEnum.CREATED} />);
    const badge = container.getElementsByClassName('ant-badge');

    expect(badge.length).toBe(1);
    expect(getByText(TestStateMap.CREATED.label)).toBeInTheDocument();
  });

  it('should render progress component', () => {
    const {container, getByText} = render(<TestState testState={TestStateEnum.EXECUTING} />);
    const progress = container.getElementsByClassName('ant-progress');

    expect(progress.length).toBe(1);
    expect(getByText(TestStateMap.EXECUTING.label)).toBeInTheDocument();
  });
});
