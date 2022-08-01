import {render} from 'test-utils';
import {ReactFlowProvider} from 'react-flow-renderer';
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
      </ReactFlowProvider>
    );

    expect(container.querySelector('form')).toBeInTheDocument();
  });
});
