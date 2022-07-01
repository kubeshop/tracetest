import {render} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import {ReduxWrapperProvider} from 'redux/ReduxWrapperProvider';
import AssertionForm from '../AssertionForm';

const defaultProps = {
  onSubmit: jest.fn(),
  onCancel: jest.fn(),
  testId: 'testId',
  runId: 'runId',
};

describe('AssertionForm', () => {
  it('should render correctly', () => {
    fetchMock.mockResponse(JSON.stringify(['spanId']));
    const {container} = render(
      <ReactFlowProvider>
        <AssertionForm {...defaultProps} />
      </ReactFlowProvider>,
      {wrapper: ReduxWrapperProvider}
    );

    expect(container.querySelector('form')).toBeInTheDocument();
  });
});
