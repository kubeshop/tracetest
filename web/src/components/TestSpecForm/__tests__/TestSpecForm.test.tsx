import {render} from 'test-utils';
import {ReactFlowProvider} from 'react-flow-renderer';
import TestSpecForm from '../TestSpecForm';

const defaultProps = {
  onSubmit: jest.fn(),
  onCancel: jest.fn(),
  testId: 'testId',
  runId: 'runId',
};

describe('TestSpecForm', () => {
  it('should render correctly', () => {
    fetchMock.mockResponse(JSON.stringify(['spanId']));
    const {container} = render(
      <ReactFlowProvider>
        <TestSpecForm {...defaultProps} />
      </ReactFlowProvider>
    );

    expect(container.querySelector('form')).toBeInTheDocument();
  });
});
